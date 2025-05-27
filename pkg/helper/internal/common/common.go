package common

import (
	"context"
	"fmt"
	"os"

	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/STARRY-S/kube-helper-mcp/pkg/wrangler"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

type Common struct {
	Wctx       *wrangler.Context
	Protocol   utils.MCPProtocol
	Listen     string
	Port       int
	ShutdownFn func(ctx context.Context) error
}

type Options struct {
	Cfg      *rest.Config
	Protocol utils.MCPProtocol
	Listen   string
	Port     int
}

func NewCommon(o *Options) *Common {
	wctx := wrangler.NewContextOrDie(o.Cfg)
	return &Common{
		Wctx:       wctx,
		Protocol:   o.Protocol,
		Listen:     o.Listen,
		Port:       o.Port,
		ShutdownFn: nil,
	}
}

func (c *Common) Start(ctx context.Context, s *server.MCPServer) error {
	var (
		addr string
		url  string
	)
	if c.Listen == "" {
		addr = fmt.Sprintf(":%v", c.Port)
		url = fmt.Sprintf("http://0.0.0.0:%v", c.Port)
	} else {
		addr = fmt.Sprintf("%v:%v", c.Listen, c.Port)
		url = fmt.Sprintf("http://%v:%v", c.Listen, c.Port)
	}

	switch c.Protocol {
	case utils.ProtocolHTTP:
		logrus.Printf("HTTP server listening on [%v/mcp]", url)
		httpServer := server.NewStreamableHTTPServer(s,
			server.WithLogger(logrus.StandardLogger()))
		c.ShutdownFn = httpServer.Shutdown
		return httpServer.Start(addr)
	case utils.ProtocolSSE:
		logrus.Printf("SSE server listening on [%v]", url)
		sseServer := server.NewSSEServer(s, server.WithBaseURL(addr))
		c.ShutdownFn = sseServer.Shutdown
		return sseServer.Start(addr)
	case utils.ProtocolStdio, "":
		stdioServer := server.NewStdioServer(s)
		return stdioServer.Listen(ctx, os.Stdin, os.Stdout)
	}
	return fmt.Errorf("unrecognized protocol %q", c.Protocol)
}

func (c *Common) Stop(ctx context.Context) error {
	if c.ShutdownFn != nil {
		return c.ShutdownFn(ctx)
	}
	return nil
}
