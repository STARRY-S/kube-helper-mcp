package main

import (
	"fmt"
	"os"
	"path"

	controllergen "github.com/rancher/wrangler/v3/pkg/controller-gen"
	"github.com/rancher/wrangler/v3/pkg/controller-gen/args"
	"github.com/rancher/wrangler/v3/pkg/crd"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	os.Unsetenv("GOPATH")

	controllergen.Run(args.Options{
		OutputPackage: "github.com/STARRY-S/kube-helper-mcp/pkg/generated",
		Boilerplate:   "pkg/codegen/boilerplate.go.txt",
		Groups: map[string]args.Group{
			corev1.GroupName: {
				Types: []any{
					corev1.Node{},
					corev1.Pod{},
					corev1.Service{},
					corev1.Namespace{},
					// corev1.Secret{},
					corev1.Endpoints{},
					// corev1.ConfigMap{},
					corev1.Event{},
				},
			},
			appsv1.GroupName: {
				Types: []any{
					appsv1.Deployment{},
					appsv1.DaemonSet{},
					appsv1.StatefulSet{},
					appsv1.ReplicaSet{},
				},
			},
			batchv1.GroupName: {
				Types: []any{
					batchv1.CronJob{},
					batchv1.Job{},
				},
			},
			discoveryv1.GroupName: {
				Types: []any{
					discoveryv1.EndpointSlice{},
				},
			},
			networkingv1.GroupName: {
				Types: []any{
					networkingv1.Ingress{},
				},
			},
		},
	})
}

func newCRD(obj any, customize func(crd.CRD) crd.CRD) crd.CRD {
	crd := crd.CRD{
		GVK: schema.GroupVersionKind{
			Group:   "flatnetwork.pandaria.io",
			Version: "v1",
		},
		Status:       true,
		SchemaObject: obj,
	}
	if customize != nil {
		crd = customize(crd)
	}
	return crd
}

func saveCRDYaml(name, data string) error {
	filepath := fmt.Sprintf("./charts/%s/templates/", name)
	if err := os.MkdirAll(filepath, 0755); err != nil {
		return fmt.Errorf("failed to mkdir %q: %w", filepath, err)
	}

	filename := path.Join(filepath, "crds.yaml")
	if err := os.WriteFile(filename, []byte(data), 0644); err != nil {
		return err
	}

	return nil
}
