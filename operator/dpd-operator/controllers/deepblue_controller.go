/*
Copyright 2022.

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

package controllers

import (
	"context"
	podv1alpha1 "dpb/api/v1alpha1"
	"dpb/controllers/deployment"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	//corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	ERP = "erp"
	CRM = "crm"
	DLM = "dlm"
)

const (
	IsExist = iota
	IsUpdated
	IsAdded
	Error
)

var (
	STATUSCODE = map[int]string{
		IsExist:   "IsExist",
		IsUpdated: "IsUpdated",
		IsAdded:   "IsAdded",
		Error:     "error",
	}
	//BusService  map[string]string
	//	"user":   ERP,
	//	"login":  ERP,
	//	"order":  CRM,
	//	"member": DLM,
	//}
)

// DeepBlueReconciler reconciles a DeepBlue object
type DeepBlueReconciler struct {
	client.Client
	Log logr.Logger
	Scheme *runtime.Scheme
}

//
//type DeepBlueList struct {
//	metav1.TypeMeta `json:",inline"`
//	metav1.ListMeta `json:"metadata,omitempty"`
//	items []DeepBlue `json:"items"`
//}
var controllerAddFuncs []func(manager.Manager) error

func init() {
	controllerAddFuncs = append(controllerAddFuncs, deployment.Add)
}

//+kubebuilder:rbac:groups=pod.dp.io,resources=deepblues,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pod.dp.io,resources=deepblues/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pod.dp.io,resources=deepblues/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeepBlue object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
//func (r *v1.Deployment)

func (r *DeepBlueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// TODO(user): your logic here

	obj := &podv1alpha1.DeepBlue{}
	err := r.Get(ctx, req.NamespacedName, obj)
	if err != nil {
		klog.Errorln(err.Error())
		return ctrl.Result{}, err
	}
	D := &v1.DeploymentList{}
	err = r.List(ctx, D, client.InNamespace("default"))
	if err != nil {
		klog.Errorln(err)
	} else {
		for _, v := range D.Items {
			t := types.NamespacedName{
				Namespace: v.Namespace,
				Name:      v.Name,
			}
			P := &v1.Deployment{}
			err := r.Get(ctx, t, P)
			if err != nil {
				klog.Errorln(err)
			} else {
				respcode := CheckAndUpdateLabels(obj.Spec.LabelMap,v)
				klog.Info("update deployment labels")
				switch respcode {
				case IsAdded, IsUpdated:
					r.Update(ctx, P)
				case IsExist:
					continue
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

func CheckAndUpdateLabels(BusService map[string]string,deploy v1.Deployment) int {
	if v, ok := deploy.ObjectMeta.Labels["node"]; ok {
		return IsExist
	} else {
		deploy.ObjectMeta.Labels["node"] = BusService[v]
		return IsUpdated
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeepBlueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "Deployment",
		Group:   "apps",
		Version: "v1",
	})
	return ctrl.NewControllerManagedBy(mgr).
		Watches(&source.Kind{Type: u}, &handler.EnqueueRequestForObject{}).
		//Watches(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForObject{}).
		For(&podv1alpha1.DeepBlue{}).
		//Owns(&v1.Deployment{}).
		Complete(r)
}
