FROM golang:1.22.0@sha256:ef61a20960397f4d44b0e729298bf02327ca94f1519239ddc6d91689615b1367 as build
WORKDIR /src

# First cache the dependencies to avoid downloading again on code change
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o scheduler ./cmd/scheduler/main.go


FROM gcr.io/distroless/base:nonroot@sha256:1a8ece87bd75cde87d0484ef48eb60ea25811baf90967265956ae4fa2098dd9d

COPY --from=build /src/scheduler /usr/local/bin/scheduler

ENTRYPOINT ["/usr/local/bin/scheduler"]
