package model

import (
	"github.com/vince15dk/k8s-operator-nhncloud/pkg/apis/nhncloud.com/v1beta1"
	"time"
)

type CreateAccessRequest struct {
	Auth Tenant `json:"auth"`
}

type Tenant struct {
	TenantId            string   `json:"tenantId"`
	PasswordCredentials UserInfo `json:"passwordCredentials"`
}

type UserInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type CreateAccessResponse struct {
	Access struct {
		Token struct {
			ID      string    `json:"id"`
			Expires time.Time `json:"expires"`
			Tenant  struct {
				ID            string `json:"id"`
				Name          string `json:"name"`
				GroupID       string `json:"groupId"`
				Description   string `json:"description"`
				Enabled       bool   `json:"enabled"`
				ProjectDomain string `json:"project_domain"`
			} `json:"tenant"`
			IssuedAt string `json:"issued_at"`
		} `json:"token"`
		ServiceCatalog []struct {
			Endpoints []struct {
				Region    string `json:"region"`
				PublicURL string `json:"publicURL"`
			} `json:"endpoints"`
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"serviceCatalog"`
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Name     string `json:"name"`
			Roles    []struct {
				Name string `json:"name"`
			} `json:"roles"`
			RolesLinks []interface{} `json:"roles_links"`
		} `json:"user"`
		Metadata struct {
			Roles   []string `json:"roles"`
			IsAdmin int      `json:"is_admin"`
		} `json:"metadata"`
	} `json:"access"`
}

type Instance struct {
	Server v1beta1.InstanceSpec `json:"server"`
}

type InstanceList struct {
	Servers []struct {
		Image struct {
			ID    string `json:"id"`
		} `json:"image"`
		Flavor             struct {
			ID    string `json:"id"`
		} `json:"flavor"`
		ID             string `json:"id"`
		Name                             string        `json:"name"`
		TenantID                         string        `json:"tenant_id"`
	} `json:"servers"`
}
