package deployment

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func Add(manager.Manager) error {
	return nil
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pod corev1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		klog.Error(err, "unable to fetch Pod")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Deployment{}).
		Complete(r)
}

//

//
//func (r *v1.de) Reconcile(req ctrl.Request) (ctrl.Result, error) {
//	sample := &samplev1.Sample{}
//	_ := r.Get(ctx, client.ObjectKey{Name: req.Name, Namespace: req.Namespace}, sample)
//	// do something
//	r.Recorder.Eventf(sample, corev1.EventTypeWarning, "Error", "some error")
//}
