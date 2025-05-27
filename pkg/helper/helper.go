package helper

import (
	"context"

	"github.com/STARRY-S/kube-helper-mcp/pkg/helper/internal/common"
	"github.com/STARRY-S/kube-helper-mcp/pkg/helper/internal/k8s"
	"github.com/STARRY-S/kube-helper-mcp/pkg/helper/internal/k8sgpt"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"k8s.io/client-go/rest"
)

type Helper interface {
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Options struct {
	Cfg      *rest.Config
	Protocol utils.MCPProtocol
	Listen   string
	Port     int
}

func NewK8sHelper(o *Options) Helper {
	return k8s.NewKubeHelper(&k8s.Options{
		Options: &common.Options{
			Cfg:      o.Cfg,
			Protocol: o.Protocol,
			Listen:   o.Listen,
			Port:     o.Port,
		},
	})
}

func NewK8sGPTHelper(o *Options) Helper {
	return k8sgpt.NewK8sGPTHelper(&k8sgpt.Options{
		Options: &common.Options{
			Cfg:      o.Cfg,
			Protocol: o.Protocol,
			Listen:   o.Listen,
			Port:     o.Port,
		},
	})
}
