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
	"dpb/controllers/deployment"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	podv1alpha1 "dpb/api/v1alpha1"
)

const(
	ERP="erp"
	CRM="crm"
	DLM="dlm"
)

const (
	IsExist=iota
	IsUpdated
	IsAdded
	Error
)

var (

	STATUSCODE=map[int]string{
	IsExist:"IsExist",
	IsUpdated:"IsUpdated",
	IsAdded:"IsAdded",
	Error:"error",
}
	BusService=map[string]string{
		"user":ERP,
		"login":ERP,
		"order": CRM,
		"member": DLM,
	}
)



// DeepBlueReconciler reconciles a DeepBlue object
type DeepBlueReconciler struct {
	client.Client
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
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=core,resources=deployment,verbs=get;list;watch;create;update;delete

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
				respcode := CheckAndUpdateLabels(v)
				klog.Info("update deployment labels")
				switch respcode{
				case IsAdded,IsUpdated:
					r.Update(ctx, P)
				case IsExist:
					continue
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

func CheckAndUpdateLabels(deploy v1.Deployment) (int) {
	if v,ok := deploy.ObjectMeta.Labels["node"];ok {
		return IsExist
	}else{
		deploy.ObjectMeta.Labels["node"]=BusService[v]
		return IsUpdated
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeepBlueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podv1alpha1.DeepBlue{}).
		Complete(r)
}


