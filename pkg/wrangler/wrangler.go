package wrangler

import (
	"github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/apps"
	"github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/batch"
	"github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/core"
	corecontroller "github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/core/v1"
	"github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/discovery.k8s.io"
	"github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/networking.k8s.io"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	appsv1 "github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/apps/v1"
	batchv1 "github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/batch/v1"
	k8sgpt "github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/core.k8sgpt.ai"
	k8sgptv1alpha1 "github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/core.k8sgpt.ai/v1alpha1"
	discoveryv1 "github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/discovery.k8s.io/v1"
	networkingv1 "github.com/STARRY-S/kube-helper-mcp/pkg/generated/controllers/networking.k8s.io/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Context struct {
	RESTConfig        *rest.Config
	Kubernetes        kubernetes.Interface
	ControllerFactory controller.SharedControllerFactory

	Core       corecontroller.Interface
	Apps       appsv1.Interface
	Networking networkingv1.Interface
	Batch      batchv1.Interface
	Discovery  discoveryv1.Interface
	K8sGPT     k8sgptv1alpha1.Interface
}

func NewContextOrDie(
	restCfg *rest.Config,
) *Context {
	// panic on error
	core := core.NewFactoryFromConfigOrDie(restCfg)
	apps := apps.NewFactoryFromConfigOrDie(restCfg)
	networking := networking.NewFactoryFromConfigOrDie(restCfg)
	batch := batch.NewFactoryFromConfigOrDie(restCfg)
	discovery := discovery.NewFactoryFromConfigOrDie(restCfg)
	k8sgpt := k8sgpt.NewFactoryFromConfigOrDie(restCfg)

	k8s, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		logrus.Fatalf("kubernetes.NewForConfig: %v", err)
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: k8s.CoreV1().Events("")})

	c := &Context{
		RESTConfig: restCfg,
		Kubernetes: k8s,

		Core:       core.Core().V1(),
		Apps:       apps.Apps().V1(),
		Networking: networking.Networking().V1(),
		Batch:      batch.Batch().V1(),
		Discovery:  discovery.Discovery().V1(),
		K8sGPT:     k8sgpt.Core().V1alpha1(),
	}
	return c
}
