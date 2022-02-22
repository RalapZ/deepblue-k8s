package main

import (
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

func main(){
	dir := homedir.HomeDir()
	//var tree string
	s := flag.String("tree", filepath.Join(dir, ".kube", "config"), "test")
	fmt.Println(s)
	config, err := clientcmd.BuildConfigFromFlags("", *s)
	if err != nil{
		panic(err.Error())
	}
	client, err := kubernetes.NewForConfig(config)
	stopCh:=make(chan struct{})
	defer close(stopCh)
	factory := informers.NewSharedInformerFactory(client, time.Minute)
	informer := factory.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}){
			myObj:=obj.(metav1.Object)
			log.Printf("New Pod Added to Store: %s", myObj.GetName())
		},
		UpdateFunc: func(oldObj interface{},newObj interface{}){
			oObj := oldObj.(metav1.Object)
			nObj := newObj.(metav1.Object)
			log.Printf("%s Pod Updated to %s", oObj.GetName(), nObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			myObj := obj.(metav1.Object)
			log.Printf("Pod Deleted from Store: %s", myObj.GetName())
		},
	})
	informer.Run(stopCh)
}