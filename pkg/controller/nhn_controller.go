package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/apis/nhncloud.com/v1beta1"
	nhnClientSet "github.com/vince15dk/k8s-operator-nhncloud/pkg/client/clientset/versioned"
	ninf "github.com/vince15dk/k8s-operator-nhncloud/pkg/client/informers/externalversions/nhncloud.com/v1beta1"
	nlister "github.com/vince15dk/k8s-operator-nhncloud/pkg/client/listers/nhncloud.com/v1beta1"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/model"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

var (
	secretName = "nhn-token"
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

	defer c.wq.Done(item)
	defer c.wq.Forget(item)

	if c.state == "update"{
		updateItem, ok := item.([2]interface{})
		if !ok {
			log.Printf("error %s", errors.New("item can not be converted"))
			return false
		}
		// get namespace and name of a new obj from que
		_ = updateItem[1].(*v1beta1.Instance).Namespace
		_ = updateItem[1].(*v1beta1.Instance).Name

	}

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
		for i :=0;i<len(instance.Spec.BlockDeviceMappingV2);i++{
			instance.Spec.BlockDeviceMappingV2[i].UUID = instance.Spec.ImageRef
		}
		rest.CreateInstance(c.client, instance, ns, secretName)
		//log.Printf("instance spec that we have it %+v\n", instance.Spec)
		err = c.updateStatus("", "", instance)
		if err != nil{
			log.Printf("error %s, updating status of the instance %s\n", err.Error(), instance.Name)
		}

	case "delete":
		fmt.Println("Delete...")
		rest.DeleteInstance(c.client, item.(*v1beta1.Instance),ns,secretName)
	case "update":
		fmt.Println("update state")
		fmt.Println(item.(*v1beta1.Instance))
		fmt.Println(ns, name)
		//rest.ValidateInstance(c.client, item.(*v1beta1.Instance),ns, secretName)
	}

	return true
}

func (c *Controller)updateStatus(id, progress string, instance *v1beta1.Instance)error{
	instance.Status.InstanceID = id
	instance.Status.Progress = progress
	_, err := c.nhnClient.NhncloudV1beta1().Instances(instance.Namespace).UpdateStatus(context.Background(), instance, metav1.UpdateOptions{})
	return err
}

func (c *Controller) handleAdd(obj interface{}) {
	log.Println("handleAdd was called")
	c.state = "create"
	c.wq.AddAfter(obj, time.Second * 2)
}

func (c *Controller) handleDelete(obj interface{}) {
	log.Println("handleDelete was called")
	c.state = "delete"
	c.wq.AddAfter(obj, time.Second * 2)
}

func (c *Controller) handleUpdate(old interface{}, new interface{}) {
	log.Println("handleUpdate was called")
	c.state = "update"
	s := [2]interface{}{
		old,
		new,
	}
	c.wq.AddRateLimited(s)
}
