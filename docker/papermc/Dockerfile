ARG JDK_VER
FROM arm64v8/eclipse-temurin:${JDK_VER}-jdk-jammy

# distribution type: vanilla, paper, waterfall
ARG DIST
ENV DIST ${DIST}

# version: 1.12.2(vanilla, paper), 1.12(waterfall)
ARG VERSION
ENV VERSION ${VERSION}

# build: for forge, velocity
ARG BUILD
ENV BUILD ${BUILD}

# install packages
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    bash \
    curl \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN pip3 install mcrcon --ignore-installed six

# download jar
WORKDIR /dist
COPY download.sh .
RUN bash ./download.sh ${DIST} ${VERSION} ${BUILD}

VOLUME ["/dist/server"]
WORKDIR /dist/server
EXPOSE 25565
