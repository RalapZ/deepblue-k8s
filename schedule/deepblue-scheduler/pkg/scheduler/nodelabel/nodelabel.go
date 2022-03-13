package plugins

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"log"
)

// 插件名称
const Name = "DeepBlue"

type  DeepBluePluginArgs  struct {}

type DeepBluePlugin struct {
	handle framework.Handle
	args   DeepBluePluginArgs
}

func (DeepBluePlugin)Name()string {
	return Name
}



//func (DeepBluePlugin)Less(*framework.QueuedPodInfo, *framework.QueuedPodInfo) bool{
//
//	return false
//}



func (*DeepBluePlugin)Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeduInfo) *framework.Status{
	log.Printf("filter pod: %##v, node: %v", pod.Name, nodeInfo)
	log.Printf("filter context:%##v",ctx)
	log.Printf("filter state: %##v",state)
	log.Printf("filter pod: %##v",pod)
	log.Printf("filter nodeInfo: %##v",nodeInfo)
	log.Printf("filter pod: %##v, node: %v", pod.Name, nodeInfo)
	log.Println(state)
	// 排除没有cpu=true标签的节点
	if nodeInfo.Node().Labels["cpu"] != "true" {
		return framework.NewStatus(framework.Unschedulable, "Node: "+nodeInfo.Node().Name)
	}

	return framework.NewStatus(framework.Success, "Node: "+nodeInfo.Node().Name)
}

//func (DeepBluePlugin) PostFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, filteredNodeStatusMap framework.NodeToStatusMap) (*framework.PostFilterResult, *framework.Status){
//
//}


//func (DeepBluePlugin)PreScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodes []*v1.Node) *framework.Status{
//
//}




func New(configuration runtime.Object,f framework.Handle)(framework.Plugin,error){
	args := &DeepBluePluginArgs{}
	//err := runtime2.DecodeInto(configuration, args)
	//if err != nil {
	//	panic(err)
	//}
	return &DeepBluePlugin{
		f,
		*args,
	},nil
}

