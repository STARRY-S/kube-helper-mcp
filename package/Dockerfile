
FROM registry.suse.com/bci/golang:1.24 AS build

WORKDIR /app
COPY . .

# Pre-download dependencies
RUN go mod download && go mod verify
# Build the binary
RUN TAG=$(shell git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0") \
    COMMIT=$(shell git rev-parse HEAD) \
    CGO_ENABLED=0 \
    go build \
    -buildmode=pie \
    -ldflags="-extldflags='-static' -s -w -X github.com/STARRY-S/kube-helper-mcp/pkg/utils.Version=${TAG#v} -X github.com/STARRY-S/kube-helper-mcp/pkg/utils.Commit=${COMMIT}" \
    -o ./helper \
    .

FROM registry.suse.com/bci/python:3.12 AS server

RUN pip install --no-cache-dir mcpo
COPY --from=build /app/helper /usr/local/bin/
COPY --from=build /app/package/mcpo/config.json /mcpo/config.json

CMD [ "mcpo", "--config", "/mcpo/config.json" ]
