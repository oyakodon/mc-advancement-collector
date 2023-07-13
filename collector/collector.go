package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"com.oykdn.mc-advancement-collector/config"
	"com.oykdn.mc-advancement-collector/lang"
	_logger "com.oykdn.mc-advancement-collector/logger"
	"com.oykdn.mc-advancement-collector/model"
	"com.oykdn.mc-advancement-collector/model/responses"
)

const (
	PROFILE_API_PATH = "https://sessionserver.mojang.com/session/minecraft/profile"

	MinecraftAdvancementTimeLayout = "2006-01-02 15:04:05 -0700"

	SpriteSize = 32
)

var (
	ErrPlayerNotFound      = fmt.Errorf("player not found")
	ErrOpenAdvancementJSON = fmt.Errorf("failed to open advancement json")
	ErrParseAdvancement    = fmt.Errorf("failed to parse advancement")

	ErrAdvancementKeyNotFound = fmt.Errorf("advancement key not found")
	ErrAdvancementConvert     = fmt.Errorf("failed to convert advancement")
)

type Collector interface {
	Player() (*responses.PlayersResponse, error)
	Load(string) (*model.PlayerAdvancementSummary, error)
	Filter(model.AdvancementFilterCondition, *model.PlayerAdvancementSummary) *model.PlayerAdvancementSummary
	Response(*model.PlayerAdvancementSummary) *responses.PlayerAdvancementResponse
}

type collector struct {
	basePath string

	ref  map[string]config.AdvancementRecord
	lang map[string]string

	playercache config.PlayerCache
	cacheSecond int
	cache       map[string]struct {
		Response model.PlayerAdvancementSummary
		Updated  time.Time
	}
}

var logger *_logger.ZapLogger = _logger.NewZapLogger()

func (c collector) Player() (*responses.PlayersResponse, error) {
	// 進捗フォルダ以下のUUID.jsonをスキャン
	files, err := os.ReadDir(c.basePath)
	if err != nil {
		return nil, err
	}

	var uuids []string
	for _, f := range files {
		_, basename := filepath.Split(f.Name())
		uuids = append(uuids, strings.Split(basename, ".")[0])
	}

	// uuidからプレイヤー名を取得
	var players []model.PlayerProfile

	eg := new(errgroup.Group)
	var mu sync.Mutex
	for _, id := range uuids {
		id := id

		eg.Go(func() error {
			// キャッシュに存在した場合はキャッシュから返却
			cache, exists := c.playercache.Players[id]
			if exists {
				// レスポンスとキャッシュに書き込み
				mu.Lock()
				defer mu.Unlock()

				players = append(players, cache)
				return nil
			}

			// Mojang APIでUUIDからProfileを得る
			profile, err := c.fetchPlayerProfile(id)
			if err != nil {
				logger.Warn(err)
			}

			if profile == nil {
				return nil
			}

			p := model.PlayerProfile{
				Id:   id,
				Name: profile.Name,
			}

			// レスポンスとキャッシュに書き込み
			mu.Lock()
			defer mu.Unlock()

			players = append(players, p)

			c.playercache.Players[id] = p
			if err := c.playercache.Save(config.PLAYERCACHE_PATH); err != nil {
				logger.Warn(err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &responses.PlayersResponse{
		Players: players,
	}, nil
}

func (c collector) Load(userId string) (*model.PlayerAdvancementSummary, error) {
	// キャッシュが存在 かつ 時間内の場合はキャッシュから返す
	if cache, exists := c.cache[userId]; exists {
		if time.Now().Before(cache.Updated.Add(time.Duration(c.cacheSecond) * time.Second)) {
			return &cache.Response, nil
		}
	}

	// jsonから進捗をロード
	var (
		filepath = path.Join(c.basePath, userId+".json")
	)
	original, updated, err := c.load(filepath)
	if err != nil {
		return nil, err
	}

	// yamlの進捗設定ファイルとjsonの進捗情報を突き合わせ, 変換処理
	var (
		mu           sync.Mutex
		wg           sync.WaitGroup
		advancements = make(map[string]*model.PlayerAdvancement)
	)
	for k := range c.ref {
		wg.Add(1)

		go func(k string) {
			defer wg.Done()
			adv, exists := original[k]
			if !exists {
				adv = &model.MinecraftAdvancement{
					Criteria: map[string]string{},
					Done:     false,
				}
			}

			converted, err := c.convert(k, adv)
			if err != nil {
				logger.Warn(err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			advancements[k] = converted
		}(k)
	}
	wg.Wait()

	// 集計結果も含めて全件レスポンス
	now := time.Now().UTC()
	resp := model.PlayerAdvancementSummary{
		Advancements: advancements,
		Progress:     *c.summarize(advancements),
		Updated:      *updated,
		Cached:       now,
	}

	// - キャッシュ更新
	c.cache[userId] = struct {
		Response model.PlayerAdvancementSummary
		Updated  time.Time
	}{
		Response: resp,
		Updated:  now,
	}

	return &resp, nil
}

func (c collector) load(filepath string) (map[string]*model.MinecraftAdvancement, *time.Time, error) {
	// JSONファイル存在確認 -> オープン
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		return nil, nil, ErrPlayerNotFound
	}

	b, err := os.ReadFile(filepath)
	if err != nil {
		logger.Warn(err)
		return nil, nil, ErrOpenAdvancementJSON
	}

	// 一旦interface{}で読み込み
	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		logger.Warn(err)
		return nil, nil, ErrParseAdvancement
	}
	// - DataVersionを除去
	delete(v, "DataVersion")

	// 再度、Advancementのmapとしてパース
	b, err = json.Marshal(v)
	if err != nil {
		logger.Warn(err)
		return nil, nil, ErrParseAdvancement
	}

	advancements := make(map[string]*model.MinecraftAdvancement)
	if err := json.Unmarshal(b, &advancements); err != nil {
		logger.Warn(err)
		return nil, nil, ErrParseAdvancement
	}

	updated := fileinfo.ModTime().UTC()
	return advancements, &updated, nil
}

func (c collector) convert(key string, original *model.MinecraftAdvancement) (*model.PlayerAdvancement, error) {
	// 設定ファイルから各種実績の属性値を読み込み
	if _, exists := c.ref[key]; !exists {
		return nil, ErrAdvancementKeyNotFound
	}
	ref := c.ref[key]

	// criteriaの達成日時をパース
	criteria := make(map[string]*time.Time)
	for _, k := range ref.Criteria {
		timestamp, exists := original.Criteria[k]
		if !exists {
			criteria[k] = nil
			continue
		}

		t, err := time.Parse(MinecraftAdvancementTimeLayout, timestamp)
		if err != nil {
			return nil, ErrAdvancementConvert
		}

		criteria[k] = &t
	}

	// 進捗集計
	total := 0
	switch ref.Calculate {
	case model.CalculateOneOf:
		total = 1
	case model.CalculateAllOf:
		total = len(ref.Criteria)
	}

	var percentage float64 = 0
	count := len(original.Criteria)
	if original.Done {
		count = total
		percentage = 1
	} else if len(original.Criteria) > 0 {
		percentage = float64(count) / float64(total)
		percentage = math.Floor(percentage*1000) / 1000
	}

	// アイコン表示
	var posx, posy int
	if ref.Icon.InvSprite {
		p := ref.Icon.Pos - 1
		posx = (p % SpriteSize) * SpriteSize
		posy = (p / SpriteSize) * SpriteSize
	}

	return &model.PlayerAdvancement{
		Key:    key,
		Parent: ref.Parent,
		Display: model.PlayerAdvancementDisplay{
			Title:       c.lang[ref.LanguageKey+lang.LANG_SUFFIX_TITLE],
			Description: c.lang[ref.LanguageKey+lang.LANG_SUFFIX_DESCRIPTION],
			Icon: model.PlayerAdvancementDisplayIcon{
				Url:       ref.Icon.Url,
				InvSprite: ref.Icon.InvSprite,
				PosX:      posx,
				PosY:      posy,
			},
		},
		Type:     ref.Type,
		Hidden:   ref.Hidden,
		Done:     original.Done,
		Criteria: criteria,
		Progress: model.AdvancementProgress{
			Total:      total,
			Done:       count,
			Percentage: percentage,
		},
	}, nil
}

func (c collector) summarize(advancements map[string]*model.PlayerAdvancement) *model.AdvancementProgress {
	total := len(advancements)

	done := 0
	var progress float64 = 0
	for _, v := range advancements {
		if v.Done {
			done += 1
			progress += 1

			continue
		}

		if len(v.Criteria) > 0 {
			progress += float64(v.Progress.Done) / float64(v.Progress.Total)
		}
	}

	return &model.AdvancementProgress{
		Total:      len(advancements),
		Done:       done,
		Percentage: math.Floor((progress/float64(total))*1000) / 1000,
	}
}

func (c collector) Filter(condition model.AdvancementFilterCondition, summary *model.PlayerAdvancementSummary) *model.PlayerAdvancementSummary {
	advancements := make(map[string]*model.PlayerAdvancement)

	switch condition {
	case model.ConditionDone:
		for k, v := range summary.Advancements {
			if v.Done {
				advancements[k] = v
			}
		}

	case model.ConditionProgress:
		checkParent := func(v *model.PlayerAdvancement) (*model.PlayerAdvancement, bool) {
			if v == nil {
				return nil, false
			}

			parent, exists := summary.Advancements[v.Parent]
			return parent, exists && parent.Done
		}

		for k, v := range summary.Advancements {
			// 達成済みの場合は表示
			if v.Done {
				advancements[k] = v
				continue
			}

			// 隠し進捗で未達成の場合は、非表示
			if v.Hidden {
				continue
			}

			// 親進捗が存在し、達成済みの場合は表示
			parent, done := checkParent(v)
			if done {
				advancements[k] = v
				continue
			}

			// 親進捗の親が存在し、達成済みの場合は表示
			if _, done := checkParent(parent); done {
				advancements[k] = v
			}
		}

	// all 及び 該当しない場合はすべて返す
	case model.ConditionAll:
		fallthrough
	default:
		for k, v := range summary.Advancements {
			advancements[k] = v
		}
	}

	return &model.PlayerAdvancementSummary{
		Advancements: advancements,
		Progress:     summary.Progress,
		Updated:      summary.Updated,
		Cached:       summary.Cached,
	}
}

func (c collector) Response(summary *model.PlayerAdvancementSummary) *responses.PlayerAdvancementResponse {
	advancements := make([]*model.PlayerAdvancement, 0, len(summary.Advancements))
	for _, v := range summary.Advancements {
		advancements = append(advancements, v)
	}

	return &responses.PlayerAdvancementResponse{
		Advancements: advancements,
		Progress:     summary.Progress,
		Updated:      summary.Updated,
		Cached:       summary.Cached,
	}
}

func (c collector) fetchPlayerProfile(id string) (*model.PlayerProfile, error) {
	resp, err := http.DefaultClient.Get(fmt.Sprintf("%s/%s", PROFILE_API_PATH, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profile model.MojangPlayerProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	return &model.PlayerProfile{
		Id:   profile.Id,
		Name: profile.Name,
	}, nil
}

func NewCollector(config *config.AppConfig, list *config.AdvancementList, lang *lang.Lang, playercache *config.PlayerCache) Collector {
	return &collector{
		basePath:    config.AdvancementPath,
		ref:         list.Advancements,
		lang:        lang.Mapping,
		playercache: *playercache,
		cacheSecond: config.Cache,
		cache: make(map[string]struct {
			Response model.PlayerAdvancementSummary
			Updated  time.Time
		}),
	}
}
