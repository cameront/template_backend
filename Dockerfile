#
# Stage 1 - build binary
#

FROM golang:1.21-alpine as backend_builder

RUN apk add --no-cache --update gcc g++

WORKDIR /app

COPY go.* /app
COPY *.go /app
COPY stuff/ /app/stuff
#... COPY [folder] /app/folder

RUN go mod download

# CGO_CFLAGS flag because of https://github.com/mattn/go-sqlite3/issues/1164#issuecomment-1635253695
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CGO_CFLAGS="-D_LARGEFILE64_SOURCE" \
  go build \
  -ldflags '-w -extldflags "-static"' \
  -mod=readonly \
  -v \
  -o /app/server \
  ./cmd/addon_server/server_main.go

#
# Stage 2 - build frontend
#

FROM node:20-buster as ui_builder

WORKDIR /_ui

COPY _ui/*.json _ui/*.js _ui/*.cjs /_ui
COPY _ui/src /_ui/src
COPY _ui/dist /_ui/dist
# Dont' want artifacts from dev to leak through
RUN rm /_ui/public/build/*

RUN npm install
RUN npm run build

#
# Stage 3 - add litestream
#

FROM debian:stable-20240110-slim AS litestream_downloader

ARG litestream_version="v0.3.13"
ARG litestream_binary_tgz_filename="litestream-${litestream_version}-linux-amd64.tar.gz"

ADD https://github.com/benbjohnson/litestream/releases/download/${litestream_version}/${litestream_binary_tgz_filename} /tmp/litestream.tar.gz
RUN tar -C /usr/local/bin -xzf /tmp/litestream.tar.gz


#
# Stage 4 - add atlas
#

FROM arigaio/atlas:0.14.2 AS atlas_image

#
# Stage 5 - minimal app container
#

FROM alpine:3.18

RUN apk add --no-cache bash

COPY --from=backend_builder /app/server /app/server
COPY --from=ui_builder /_ui/dist /app/_ui/dist
COPY --from=litestream_downloader /usr/local/bin/litestream /app/litestream
COPY --from=atlas_image /atlas /atlas

# To give access to timezone info
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

COPY ent/migrate/migrations /app/migrations
COPY static/ /app/static
COPY docker_entrypoint /app/docker_entrypoint
COPY litestream.yml /etc/litestream.yml
COPY gcp_service_account_database_replicator.json /app/_creds/gcp_service_account_database_replicator.json

WORKDIR /app

ENTRYPOINT ["/app/docker_entrypoint"]