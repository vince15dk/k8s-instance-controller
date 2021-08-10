package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/model"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	getTokenURL = "https://api-identity.infrastructure.cloud.toast.com/v2.0/tokens"
	baseUrl = "https://kr1-api-instance.infrastructure.cloud.toast.com/v2/"
)

func GetToken(s *Secret)model.CreateAccessResponse{
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	b := &model.CreateAccessRequest{Auth: model.Tenant{
		TenantId: s.TenantId,
		PasswordCredentials: model.UserInfo{
			UserName: s.UserName,
			Password: s.Password,
		},
	}}

	response, _ := PostHandleFunc(getTokenURL, b, headers)
	bytes, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	var result model.CreateAccessResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal create repo successful response: %s", err.Error()))
	}
	return result
}

func PostHandleFunc(url string, body interface{}, headers http.Header)(*http.Response, error){
	jsonBytes, err := json.Marshal(body)
	if err != nil{
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers
	client := http.Client{}
	return client.Do(request)
}