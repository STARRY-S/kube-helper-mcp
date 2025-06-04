
FROM registry.suse.com/bci/golang:1.24 AS build

WORKDIR /app

# Pre-download dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Build the binary
ENV GOCACHE=/root/.cache/go-build
ARG TAG=v0.1.0-alpha
ARG COMMIT=unknown
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 \
    go build \
    -buildmode=pie \
    -ldflags="-extldflags='-static' -s -w -X github.com/STARRY-S/kube-helper-mcp/pkg/utils.Version=${TAG#v} -X github.com/STARRY-S/kube-helper-mcp/pkg/utils.Commit=${COMMIT}" \
    -o /usr/local/bin/kube-helper-mcp \
    .
RUN kube-helper-mcp version

FROM registry.suse.com/bci/python:3.12 AS server

RUN pip install --no-cache-dir mcpo
COPY --from=build /usr/local/bin/kube-helper-mcp /usr/local/bin/
COPY package/mcpo/ /mcpo/

# Run Streamable HTTP Server in default mode
ENTRYPOINT [ "kube-helper-mcp" ]
CMD [ "k8s", "--http" ]
