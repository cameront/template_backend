#
# Stage 1 - build binary (on debian)
#
ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-bookworm as backend_builder

WORKDIR /app
COPY go.* /app
COPY *.go /app
COPY stuff/ /app/stuff
#... COPY [folder] /app/folder

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -v -o /app/server ./cmd/server


#
# Stage 2 - build frontend
#
FROM node:23-bookworm as ui_builder

WORKDIR /_ui
COPY _ui/*.json _ui/*.js _ui/*.ts _ui/*.html /_ui
COPY _ui/src /_ui/src
COPY _ui/public /_ui/public

RUN npm install
RUN npm run build


#
# Stage 3 - add litestream
#
FROM debian:bookworm-slim AS litestream_downloader

ARG litestream_version="v0.3.13"
ADD https://github.com/benbjohnson/litestream/releases/download/${litestream_version}/litestream-${litestream_version}-linux-amd64.tar.gz /tmp/litestream.tar.gz
RUN tar -C /usr/local/bin -xzf /tmp/litestream.tar.gz


#
# Stage 4 - add atlas
#
ARG atlas_version="0.28.1"
FROM arigaio/atlas:${atlas_version} AS atlas_image


#
# Stage 5 - Small(ish, but debian-based) app container. I don't want to use
# alpine here because of the potential to be bitten by the glibc vs musl issues.
#
FROM debian:bookworm-slim

COPY --from=backend_builder /app/server /app/server
COPY --from=backend_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=ui_builder /_ui/dist /app/_ui/dist
COPY --from=litestream_downloader /usr/local/bin/litestream /usr/local/bin/litestream
COPY --from=atlas_image /atlas /atlas

COPY ent/migrate/migrations /app/migrations
COPY docker_entrypoint /app/docker_entrypoint
COPY static/ /app/static
#COPY gcp_service_account_database_replicator.json /app/_creds/gcp_service_account_database_replicator.json
#COPY ./litestream.yml /etc/litestream.yml

WORKDIR /app

ENTRYPOINT ["/app/docker_entrypoint"]