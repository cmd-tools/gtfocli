ARG GO_VERSION=1.21.3
FROM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /output/gtfocli . \
    && cd /output \
    && /output/gtfocli search

FROM alpine AS final
WORKDIR /app
RUN apk add git
COPY --from=build /output /app
ENTRYPOINT [ "/app/gtfocli" ]
