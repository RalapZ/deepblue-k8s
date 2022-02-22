package deployment

import (
	"k8s.io/apimachinery/pkg/runtime"
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

//
//func (r *v1.de) Reconcile(req ctrl.Request) (ctrl.Result, error) {
//	sample := &samplev1.Sample{}
//	_ := r.Get(ctx, client.ObjectKey{Name: req.Name, Namespace: req.Namespace}, sample)
//	// do something
//	r.Recorder.Eventf(sample, corev1.EventTypeWarning, "Error", "some error")
//}