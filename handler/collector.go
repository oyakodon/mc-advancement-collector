package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	_collector "com.oykdn.mc-advancement-collector/collector"
	"com.oykdn.mc-advancement-collector/config"
	"com.oykdn.mc-advancement-collector/lang"
	_logger "com.oykdn.mc-advancement-collector/logger"
	"com.oykdn.mc-advancement-collector/model"
	"com.oykdn.mc-advancement-collector/model/requests"
	"com.oykdn.mc-advancement-collector/model/responses"
)

const (
	LANG_PATH = "./lang"
)

var logger *_logger.ZapLogger = _logger.NewZapLogger()

func main() {
	if os.Getenv("GIN_DEBUG") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	lang, err := lang.LoadLang(fmt.Sprintf("%s/%s.json", LANG_PATH, conf.AppConfig.Language))
	if err != nil {
		panic(err)
	}

	collector := _collector.NewCollector(conf.AppConfig, conf.AdvancementList, lang, conf.PlayerCache)

	r := gin.New()

	r.Use(ginzap.Ginzap(logger.Zap(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.Zap(), true))
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://mc.oykdn.com",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			"GET",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		MaxAge: 24 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	v1.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})

	v1.GET("/players", func(c *gin.Context) {
		p, err := collector.Player()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.IndentedJSON(http.StatusOK, p)
	})

	advancement := v1.Group("/advancement")

	advancement.GET("/:id", func(c *gin.Context) {
		var p requests.PlayerAdvancementRequest

		if err := c.ShouldBindUri(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		// クエリで条件を指定できるようにする
		if err := c.ShouldBindQuery(&p); err != nil {
			logger.Debug(err)
		}
		condition := p.Condition
		switch condition {
		case model.ConditionAll:
			break
		case model.ConditionDone:
			break
		case model.ConditionProgress:
			break

		default:
			condition = model.ConditionProgress
		}

		// プレイヤーの進捗情報を取得
		advancements, err := collector.Load(p.PlayerId)
		if err != nil {
			code := http.StatusInternalServerError

			switch err {
			case _collector.ErrPlayerNotFound:
				// PlayerNotFoundの場合は404で返す
				code = http.StatusNotFound
			}

			c.JSON(code, gin.H{
				"message": err.Error(),
			})
			return
		}

		resp := collector.Response(collector.Filter(condition, advancements))

		// pretty print & no escape でJSONを返却
		c.Status(http.StatusOK)

		header := c.Writer.Header()
		if val := header["Content-Type"]; len(val) == 0 {
			header["Content-Type"] = []string{"application/json; charset=utf-8"}
		}

		enc := json.NewEncoder(c.Writer)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "    ")
		if err := enc.Encode(resp); err != nil {
			panic(err)
		}
	})

	advancement.GET("/assets", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, responses.ConvertToAdvancementAssetsResponse(conf.AppConfig.Assets.Background))
	})

	if err := r.Run(":18080"); err != nil {
		logger.Fatal(err)
	}
}
