FROM golang:1.21.6@sha256:7b575fe0d9c2e01553b04d9de8ffea6d35ca3ab3380d2a8db2acc8f0f1519a53 as build

# Note: Dockerfile uses paths relative to the top-level project directory,
# so it should be built from that directory, i.e:
# $ cd package-analysis
# $ docker build -f sandboxes/staticanalysis/Dockerfile .

# Cache dependencies to avoid downloading again on code change
WORKDIR /src
# Dependencies for package analysis
COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . ./

WORKDIR /src/sandboxes/staticanalysis
# If CGO is disabled then we don't need glibc
RUN CGO_ENABLED=0 go build -o staticanalyze staticanalyze.go

FROM alpine:3.19.0@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48
RUN apk add --no-cache file && \
	apk add --no-cache nodejs && \
	apk add --no-cache npm && \
	apk add --no-cache python3

COPY --from=build /src/sandboxes/staticanalysis/staticanalyze /usr/local/bin/staticanalyze
RUN chmod 755 /usr/local/bin/staticanalyze

RUN mkdir /npm_deps
COPY --from=build /src/internal/staticanalysis/parsing/package.json /src/internal/staticanalysis/parsing/package-lock.json /npm_deps/

# cache NPM installs in /npm_cache so that static analysis binary can use them
RUN mkdir -m 755 /npm_cache && \
	npm ci --prefix /npm_deps --cache /npm_cache

WORKDIR /app

ENTRYPOINT [ "sleep" ]
CMD [ "30m" ]
