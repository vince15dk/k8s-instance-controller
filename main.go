package main

import (
	"flag"
	nhnClient "github.com/vince15dk/k8s-operator-nhncloud/pkg/client/clientset/versioned"
	nInfFac "github.com/vince15dk/k8s-operator-nhncloud/pkg/client/informers/externalversions"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

func main(){
	var kubeconfig *string
	if home := homedir.HomeDir(); home != ""{
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	}else{
		kubeconfig = flag.String("kubeconfig", "", "absulte path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil{
		log.Printf("Building ocnfig from flags, %s", err.Error())
	}

	nhnClientSet, err := nhnClient.NewForConfig(config)
	if err != nil{
		log.Printf("getting nhnclient set %s\n", err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil{
		log.Printf("getting std client %s\n", err.Error())
	}

	// period of re-sync by calling UpdateFunc of the event handler
	infoFactory := nInfFac.NewSharedInformerFactory(nhnClientSet, 20 * time.Minute)
	ch := make(chan struct{})
	c := controller.NewController(client, nhnClientSet, infoFactory.Nhncloud().V1beta1().Instances())

	infoFactory.Start(ch)

	if err := c.Run(ch); err != nil{
		log.Printf("error running controller %s\n", err.Error())
	}

}