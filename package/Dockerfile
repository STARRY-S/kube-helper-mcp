
FROM registry.suse.com/bci/golang:1.24 AS build

WORKDIR /app

# Pre-download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Build the binary
COPY . .
RUN TAG=$(git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0") \
    COMMIT=$(git rev-parse HEAD) \
    CGO_ENABLED=0 \
    go build \
    -buildmode=pie \
    -ldflags="-extldflags='-static' -s -w -X github.com/STARRY-S/kube-helper-mcp/pkg/utils.Version=${TAG#v} -X github.com/STARRY-S/kube-helper-mcp/pkg/utils.Commit=${COMMIT}" \
    -o ./kube-helper-mcp \
    .

FROM registry.suse.com/bci/python:3.12 AS server

RUN pip install --no-cache-dir mcpo
COPY --from=build /app/kube-helper-mcp /usr/local/bin/
COPY --from=build /app/package/mcpo/ /mcpo/

# Run SSE server in default mode
CMD [ "kube-helper-mcp", "run", "--sse" ]
