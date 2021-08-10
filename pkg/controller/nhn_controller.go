package controller

import (
	"fmt"
	nhnClientSet "github.com/vince15dk/k8s-operator-nhncloud/pkg/client/clientset/versioned"
	ninf "github.com/vince15dk/k8s-operator-nhncloud/pkg/client/informers/externalversions/nhncloud.com/v1beta1"
	nlister "github.com/vince15dk/k8s-operator-nhncloud/pkg/client/listers/nhncloud.com/v1beta1"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/model"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

type Controller struct {
	client        kubernetes.Interface
	nhnClient     nhnClientSet.Interface
	clusterSynced cache.InformerSynced
	nLister       nlister.InstanceLister
	wq            workqueue.RateLimitingInterface
	state         string
}

func NewController(client kubernetes.Interface, nhnClient nhnClientSet.Interface, nInformer ninf.InstanceInformer) *Controller {
	c := &Controller{
		client:        client,
		nhnClient:     nhnClient,
		clusterSynced: nInformer.Informer().HasSynced,
		nLister:       nInformer.Lister(),
		wq:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "instance"),
	}

	nInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDelete,
			UpdateFunc: c.handleUpdate,
		})
	return c
}

func (c *Controller) Run(ch chan struct{}) error {
	if ok := cache.WaitForCacheSync(ch, c.clusterSynced); !ok {
		log.Println("cache was not synced")
	}

	go wait.Until(c.worker, time.Second, ch)

	<-ch
	return nil
}

func (c *Controller) worker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	item, shutDown := c.wq.Get()
	if shutDown {
		return false
	}

	defer c.wq.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("error %s called Namespace key func on cache for item", err.Error())
		return false
	}
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("error %s, Getting the namespace from lister", err.Error())
		return false
	}
	switch c.state{
	case "create":
		instance, err := c.nLister.Instances(ns).Get(name)
		if err != nil {
			log.Printf("error %s, Getting the instance resource from lister", err.Error())
		}
		instance.Spec.ImageRef = model.Images[instance.Spec.ImageRef]
		instance.Spec.FlavorRef = model.Flavors[instance.Spec.FlavorRef]
		for _, v := range instance.Spec.BlockDeviceMappingV2{
			v.UUID = instance.Spec.ImageRef
		}
		log.Printf("instance spec that we have it %+v\n", instance.Spec)
	case "delete":
		fmt.Println("delete state")
	case "update":
		fmt.Println("update state")
	}
	//InstanceID, err := do.Create(c.client, instance.Spec)
	return true
}

func (c *Controller) handleAdd(obj interface{}) {
	log.Println("handleAdd was called")
	c.state = "create"
	c.wq.Add(obj)
}

func (c *Controller) handleDelete(obj interface{}) {
	log.Println("handleDelete was called")
	c.state = "delete"
	c.wq.Add(obj)
}

func (c *Controller) handleUpdate(old interface{}, obj interface{}) {
	c.state = "update"
	log.Println("handleUpdate was called")
}
