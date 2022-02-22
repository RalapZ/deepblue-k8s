package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)
func main(){
	//config
	dir := homedir.HomeDir()
	kubeconfig:=flag.String("kubeconfig",filepath.Join(dir,".kube","config"),"kube config")
	config, err := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err !=nil{
		panic(err.Error())
	}
	//client
	client, err := kubernetes.NewForConfig(config)
	if err !=nil{
		panic(err.Error())
	}
	//resource
	core:=client.CoreV1()
	podDefault:=core.Pods("default")
	ctx:=context.TODO()
	list, err := podDefault.List(context.TODO(), v1.ListOptions{})
	for _,v := range list.Items {
		get, _ := podDefault.Get(ctx, v.Name, v1.GetOptions{})
		fmt.Println(get.Status.PodIP)

	}


}