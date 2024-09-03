FROM golang:alpine AS builder

RUN go env -w GO111MODULE=auto \
  && go env -w CGO_ENABLED=0

ARG BUILD_TIME
ARG COMMIT_ID
ARG Tags

WORKDIR /build

COPY ./ .

RUN set -ex \
    && cd /build \
    && go build -ldflags "-X 'github.com/starskim/DDBOT/lsp.BuildTime=$BUILD_TIME' -X 'github.com/starskim/DDBOT/lsp.CommitId=$COMMIT_ID' -X 'github.com/starskim/DDBOT/lsp.Tags=$Tags'" -o DDBOT github.com/starskim/DDBOT/cmd

FROM alpine:latest

COPY docker-entrypoint.sh /docker-entrypoint.sh

RUN chmod +x /docker-entrypoint.sh && \
    apk add --no-cache --update \
      ffmpeg \
      coreutils \
      shadow \
      su-exec \
      tzdata && \
    rm -rf /var/cache/apk/* && \
    mkdir -p /app && \
    mkdir -p /data && \
    mkdir -p /config && \
    useradd -d /config -s /bin/sh abc && \
    chown -R abc /config && \
    chown -R abc /data

ENV TZ="Asia/Shanghai"
ENV UID=99
ENV GID=100
ENV UMASK=002

COPY --from=builder /build/DDBOT /app/

WORKDIR /data

VOLUME [ "/data" ]

ENTRYPOINT [ "/docker-entrypoint.sh" ]
CMD [ "/app/DDBOT" ]