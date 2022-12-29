# 基础镜像,基于golang的alpine版本进行构建
FROM golang:alpine3.17

# 作者:库陈胜
MAINTAINER kucs<kucs@isyscore.com>

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

WORKDIR $GOPATH/workspace

RUN go mod init github.com/kuchensheng/project \
    && go get github.com/kuchensheng/bintools \
    && go get github.com/kuchensheng/bintools/json@v1.0.1 \
    && go get github.com/kuchensheng/bintools/tracer/trace \
    && go get github.com/gin-gonic/gin \
    && go get github.com/rs/zerolog

RUN mkdir ./json
COPY . ./json
RUN go build -o server ./json/server.go
RUN rm -rf ./json
RUN mkdir template
COPY ./template ./template
CMD ["./server"]


