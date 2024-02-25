ARG GO_VERSION=1.21.3
FROM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN mkdir -p /output \
    chmod +w /output

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /output/gtfocli . \
    && cd /output \
    && /output/gtfocli update

FROM scratch AS final
WORKDIR /bin
COPY --from=build /output /bin
ENTRYPOINT [ "/bin/gtfocli" ]
