/*
Copyright 2025 SUSE Rancher

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"sync"
	"time"

	v1alpha1 "github.com/k8sgpt-ai/k8sgpt-operator/api/v1alpha1"
	"github.com/rancher/wrangler/v3/pkg/apply"
	"github.com/rancher/wrangler/v3/pkg/condition"
	"github.com/rancher/wrangler/v3/pkg/generic"
	"github.com/rancher/wrangler/v3/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// K8sGPTController interface for managing K8sGPT resources.
type K8sGPTController interface {
	generic.ControllerInterface[*v1alpha1.K8sGPT, *v1alpha1.K8sGPTList]
}

// K8sGPTClient interface for managing K8sGPT resources in Kubernetes.
type K8sGPTClient interface {
	generic.ClientInterface[*v1alpha1.K8sGPT, *v1alpha1.K8sGPTList]
}

// K8sGPTCache interface for retrieving K8sGPT resources in memory.
type K8sGPTCache interface {
	generic.CacheInterface[*v1alpha1.K8sGPT]
}

// K8sGPTStatusHandler is executed for every added or modified K8sGPT. Should return the new status to be updated
type K8sGPTStatusHandler func(obj *v1alpha1.K8sGPT, status v1alpha1.K8sGPTStatus) (v1alpha1.K8sGPTStatus, error)

// K8sGPTGeneratingHandler is the top-level handler that is executed for every K8sGPT event. It extends K8sGPTStatusHandler by a returning a slice of child objects to be passed to apply.Apply
type K8sGPTGeneratingHandler func(obj *v1alpha1.K8sGPT, status v1alpha1.K8sGPTStatus) ([]runtime.Object, v1alpha1.K8sGPTStatus, error)

// RegisterK8sGPTStatusHandler configures a K8sGPTController to execute a K8sGPTStatusHandler for every events observed.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterK8sGPTStatusHandler(ctx context.Context, controller K8sGPTController, condition condition.Cond, name string, handler K8sGPTStatusHandler) {
	statusHandler := &k8sGPTStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, generic.FromObjectHandlerToHandler(statusHandler.sync))
}

// RegisterK8sGPTGeneratingHandler configures a K8sGPTController to execute a K8sGPTGeneratingHandler for every events observed, passing the returned objects to the provided apply.Apply.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterK8sGPTGeneratingHandler(ctx context.Context, controller K8sGPTController, apply apply.Apply,
	condition condition.Cond, name string, handler K8sGPTGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &k8sGPTGeneratingHandler{
		K8sGPTGeneratingHandler: handler,
		apply:                   apply,
		name:                    name,
		gvk:                     controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterK8sGPTStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type k8sGPTStatusHandler struct {
	client    K8sGPTClient
	condition condition.Cond
	handler   K8sGPTStatusHandler
}

// sync is executed on every resource addition or modification. Executes the configured handlers and sends the updated status to the Kubernetes API
func (a *k8sGPTStatusHandler) sync(key string, obj *v1alpha1.K8sGPT) (*v1alpha1.K8sGPT, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type k8sGPTGeneratingHandler struct {
	K8sGPTGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
	seen  sync.Map
}

// Remove handles the observed deletion of a resource, cascade deleting every associated resource previously applied
func (a *k8sGPTGeneratingHandler) Remove(key string, obj *v1alpha1.K8sGPT) (*v1alpha1.K8sGPT, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.K8sGPT{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	if a.opts.UniqueApplyForResourceVersion {
		a.seen.Delete(key)
	}

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

// Handle executes the configured K8sGPTGeneratingHandler and pass the resulting objects to apply.Apply, finally returning the new status of the resource
func (a *k8sGPTGeneratingHandler) Handle(obj *v1alpha1.K8sGPT, status v1alpha1.K8sGPTStatus) (v1alpha1.K8sGPTStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.K8sGPTGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}
	if !a.isNewResourceVersion(obj) {
		return newStatus, nil
	}

	err = generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
	if err != nil {
		return newStatus, err
	}
	a.storeResourceVersion(obj)
	return newStatus, nil
}

// isNewResourceVersion detects if a specific resource version was already successfully processed.
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *k8sGPTGeneratingHandler) isNewResourceVersion(obj *v1alpha1.K8sGPT) bool {
	if !a.opts.UniqueApplyForResourceVersion {
		return true
	}

	// Apply once per resource version
	key := obj.Namespace + "/" + obj.Name
	previous, ok := a.seen.Load(key)
	return !ok || previous != obj.ResourceVersion
}

// storeResourceVersion keeps track of the latest resource version of an object for which Apply was executed
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *k8sGPTGeneratingHandler) storeResourceVersion(obj *v1alpha1.K8sGPT) {
	if !a.opts.UniqueApplyForResourceVersion {
		return
	}

	key := obj.Namespace + "/" + obj.Name
	a.seen.Store(key, obj.ResourceVersion)
}
