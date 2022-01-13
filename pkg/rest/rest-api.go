package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/apis/nhncloud.com/v1beta1"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/model"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"strings"
)

type Secret struct {
	TenantId string
	UserName string
	Password string
}

func SettingAuthHeader(h *http.Header, token string) *http.Header {
	h.Set("Content-Type", "application/json")
	h.Set("X-Auth-Token", token)
	return h
}

func getSecret(client kubernetes.Interface, namespace, name string) (*Secret, error) {
	s, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	secret := &Secret{}
	secret.TenantId = string(s.Data["tenantId"])
	secret.UserName = string(s.Data["userName"])
	secret.Password = string(s.Data["password"])
	return secret, nil
}

func AddToBodyInstance(a *model.Instance, b *v1beta1.Instance){
	a.Server.Name = b.Spec.Name
	a.Server.ImageRef = model.Images[b.Spec.ImageRef]
	a.Server.FlavorRef = model.Flavors[b.Spec.FlavorRef]
	a.Server.Networks = b.Spec.Networks
	a.Server.KeyName = b.Spec.KeyName
	a.Server.MinCount = b.Spec.MinCount
	a.Server.BlockDeviceMappingV2 = b.Spec.BlockDeviceMappingV2
	a.Server.BlockDeviceMappingV2[0].UUID = model.Images[b.Spec.ImageRef]
	a.Server.BlockDeviceMappingV2[0].DeviceName = "vda"
	a.Server.BlockDeviceMappingV2[0].SourceType = "image"
	a.Server.BlockDeviceMappingV2[0].DestinationType = "volume"
	a.Server.BlockDeviceMappingV2[0].DeleteOnTermination = 1
}


func CreateInstance(client kubernetes.Interface, instance *v1beta1.Instance, namespace, name string) {
	s, err := getSecret(client, namespace, name)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
	// Get Token
	token := GetToken(s)

	// Set Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token.Access.Token.ID)

	// Create Instance
	url := baseUrl + s.TenantId + "/servers"
	inst := model.Instance{Server: instance.Spec}
	resp, err := PostHandleFunc(url, inst, *newHeader)
	if err != nil {
		fmt.Println(err)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func DeleteInstance(client kubernetes.Interface, instance *v1beta1.Instance, namespace, name string) {
	s, err := getSecret(client, namespace, name)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
	// Get Token
	token := GetToken(s)

	// Set Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token.Access.Token.ID)

	// Delete Instance
	url := baseUrl + s.TenantId + "/servers/detail"
	resp, err := ListHandleFunc(url, *newHeader)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	instances := &model.InstanceList{}
	var instanceIds []string
	var num = 1
	err = json.Unmarshal(bytes, instances)
	for _, v := range instances.Servers {
		for i := num; i <= instance.Spec.MinCount; i++ {
			if v.Name == fmt.Sprintf("%s-%d", instance.Spec.Name, i) {
				instanceIds = append(instanceIds, v.ID)
			}
		}
	}
	for _, v := range instanceIds{
		url := baseUrl + s.TenantId + "/servers/" + v
		_, err := DeleteHandleFunc(url, *newHeader)
		if err != nil{
		}
	}
}

func ListInstance(client kubernetes.Interface, instance *v1beta1.Instance, namespace, name string){
	secret, err := getSecret(client, namespace, name)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
	// Get Token
	token := GetToken(secret)

	// Set Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token.Access.Token.ID)

	// update list
	urlGetInstance := baseUrl + secret.TenantId + "/servers/detail"
	newResponse, err := ListHandleFunc(urlGetInstance, *newHeader)
	if err != nil {
		fmt.Println(err)
	}
	//servers := &model.ServerInfo{}
	newBytes, err := ioutil.ReadAll(newResponse.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer newResponse.Body.Close()
	servers := &model.InstanceList{}
	err = json.Unmarshal(newBytes, servers)
	var s []string
	for _, v := range servers.Servers{
		if strings.Split(v.Name, "-")[0] == instance.Spec.Name{
			s = append(s, strings.Split(v.Name, "-")[0])
		}
	}

	diff := len(s) -instance.Spec.MinCount
	url := baseUrl + secret.TenantId + "/servers"
	inst := &model.Instance{}
	if diff < 0 {
		AddToBodyInstance(inst, instance)
		diff *= -1
		if diff == 1{
			inst.Server.Name = instance.Spec.Name+ fmt.Sprintf("-%d", diff + 1)
		}else{
			inst.Server.Name = instance.Spec.Name
		}
		inst.Server.MinCount = diff
		resp, err := PostHandleFunc(url, inst, *newHeader)
		if err != nil {
			fmt.Println(err)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	}else if diff > 0 {
		var serverIds []string
		for _, v := range servers.Servers{
			if strings.Split(v.Name, "-")[0] == instance.Spec.Name{
				serverIds = append(serverIds, v.ID)
			}
		}
		serverIds = serverIds[:diff]
		for _, v := range serverIds{
			urlDeleteInstance := baseUrl + secret.TenantId + "/servers/" + v
			resp, err := DeleteHandleFunc(urlDeleteInstance, *newHeader)
			if err != nil {
				fmt.Println(err)
			}
			resp.Body.Close()
		}
	}

}

func ValidateInstance(client kubernetes.Interface, instance *v1beta1.Instance, namespace, name string) {
	s, err := getSecret(client, namespace, name)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
	// Get Token
	token := GetToken(s)

	// Set Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token.Access.Token.ID)
	url := baseUrl + s.TenantId + "/servers/detail"
	resp, err := ListHandleFunc(url, *newHeader)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	instances := &model.InstanceList{}
	err = json.Unmarshal(bytes, instances)
	var names []string
	for _, v := range instances.Servers {
		if strings.Split(v.Name, "-")[0] == instance.Spec.Name {
			names = append(names, strings.Split(v.Name, "-")[0])
		}
	}
	diff := len(names) - instance.Spec.MinCount
	url = baseUrl + s.TenantId + "/servers"
	inst := model.Instance{Server: instance.Spec}
	if diff < 0 {
		diff *= -1
		inst.Server.MinCount = diff
		_, err := PostHandleFunc(url, inst, *newHeader)
		if err != nil {
			fmt.Println(err)
			return
		}
	}else if diff > 0{
		var instanceIds []string
		for _, v := range instances.Servers{
			if strings.Split(v.Name, "-")[0] == instance.Spec.Name{
				instanceIds = append(instanceIds, v.ID)
			}
		}
		instanceIds = instanceIds[:diff]
		for _, v := range instanceIds{
			url = baseUrl + s.TenantId + "/servers/" + v
			_, err := DeleteHandleFunc(url, *newHeader)
			if err != nil{
				fmt.Println(err)
				return
			}
		}
	}
}