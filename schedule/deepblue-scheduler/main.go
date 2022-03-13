package main

import (
	plugins "github.com/RalapZ/deepblue-k8s-scheduler/pkg/scheduler/nodelabel"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"math/rand"
	"time"
)

func main(){
	rand.Seed(time.Now().UnixNano())
	command := app.NewSchedulerCommand(app.WithPlugin(plugins.Name, plugins.New))
	command.Execute()
}
