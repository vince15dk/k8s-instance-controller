package rest

import (
	"context"
	"fmt"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/apis/nhncloud.com/v1beta1"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/model"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
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

func getSecret(client kubernetes.Interface, namespace, name string)(*Secret, error){
	s, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil{
		return nil, err
	}
	secret := &Secret{}
	secret.TenantId = string(s.Data["tenantId"])
	secret.UserName = string(s.Data["userName"])
	secret.Password = string(s.Data["password"])
	return secret, nil
}

func CreateInstance(client kubernetes.Interface,instance *v1beta1.Instance, namespace, name string){
	s, err := getSecret(client, namespace, name)
	if err != nil{
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
	if err != nil{
		fmt.Println(err)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println(err)
	}
	defer resp.Body.Close()
}



