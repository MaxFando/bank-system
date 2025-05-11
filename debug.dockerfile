FROM golang:1.23-alpine AS build
RUN apk update && apk add --no-cache curl make git
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /src
COPY go.mod .
COPY go.sum .

RUN go env -w GOCACHE=/root/.cache/go-build
RUN --mount=type=secret,id=netrc,dst=/root/.netrc \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -gcflags="all=-N -l" -o app ./cmd/app

FROM alpine:latest
COPY --from=build /go/bin/dlv /usr/local/bin/dlv
WORKDIR /src
COPY --from=build /src/app .
COPY --from=build /src/.env ./.env

ENTRYPOINT ["dlv"]