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
	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ralapiov1alpha1 "sec/api/v1alpha1"
)

const (
	IsExist = iota
	IsUpdated
	IsAdded
	Error
)

// SecReconciler reconciles a Sec object
type SecReconciler struct {
	client.Client
	logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ralap.io.sec,resources=secs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ralap.io.sec,resources=secs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ralap.io.sec,resources=secs/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps/v1,resources=deployment,verbs=get;create;list;watch;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Sec object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *SecReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	c, err := r.ReconcileDeployment(ctx, req)
	if err == nil{
		return c ,err
	}

	return ctrl.Result{}, nil
}




func CheckAndUpdateLabels(obj *ralapiov1alpha1.Sec, deploy v1.Deployment) int {
	if v, ok := deploy.ObjectMeta.Labels["node"]; ok {
		obj.Status.LabelStatus[deploy.Name] = obj.Labels[deploy.Name]
		return IsExist
	} else {
		deploy.ObjectMeta.Labels["node"] = obj.ObjectMeta.Labels[v]
		obj.Status.LabelStatus[deploy.Name] = obj.Labels[deploy.Name]
		return IsUpdated
	}
}

func (r *SecReconciler)ReconcileDeployment(ctx context.Context, req ctrl.Request)(ctrl.Result,error){
	D:=&v1.Deployment{}
	if  err :=r.Get(ctx, req.NamespacedName, D); err!=nil {
		return  ctrl.Result{},err
	}else{
		sl:=&ralapiov1alpha1.SecList{}
		err := r.List(ctx,sl)
		if err !=nil{
			return ctrl.Result{},err
		}
		for _,v := range sl.Items{
			s:=&ralapiov1alpha1.Sec{}
			err := r.Get(ctx, types.NamespacedName{Namespace: v.Namespace, Name: v.Name}, s)
			if err !=nil{
				continue
			}else{
				seclabels := s.ObjectMeta.GetLabels()
				deploylabels:=D.ObjectMeta.GetLabels()
				if v,ok := deploylabels["node"];ok{
					if seclabels[D.Name]==v{
						return ctrl.Result{},nil
					}else{
						deploylabels["node"]=seclabels[D.Name]
						return ctrl.Result{},err
					}
				}
			}
		}
	}
	
	return ctrl.Result{},nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecReconciler) SetupWithManager(mgr ctrl.Manager) error {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "Deployment",
		Group:   "apps",
		Version: "v1",
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(&ralapiov1alpha1.Sec{}).
		Watches(&source.Kind{Type: u}, &handler.EnqueueRequestForObject{}).
		WithEventFilter(&NodeLabelPredicate{}).
		Complete(r)
}

type NodeLabelPredicate struct {
	predicate.Funcs
}

func (rl *NodeLabelPredicate) Update(e event.UpdateEvent) bool {
	//klog.Infof("new:%##v\n old:%##v\n", e.ObjectNew, e.ObjectOld)
	if !compareMaps(e.ObjectOld.GetLabels(), e.ObjectNew.GetLabels()) {
		klog.Infof("Newlabel:%##v \n,Oldlabels", e.ObjectNew.GetLabels(), e.ObjectOld.GetLabels())
		return true
	}
	return false
}

func (rl *NodeLabelPredicate) Create(e event.CreateEvent) bool {
	return false
}

func compareMaps(old map[string]string, new map[string]string) bool {
	if reflect.DeepEqual(old, new) {
		return true
	}
	return false
}
