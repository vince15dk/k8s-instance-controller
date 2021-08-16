package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="InstanceID",type=string,JSONPath=`.status.instanceID`
// +kubebuilder:printcolumn:name="Progress",type=string,JSONPath=`.status.progress`
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

type InstanceStatus struct {
	InstanceID string `json:"instanceID,omitempty"`
	Progress   string `json:"progress,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Instance `json:"items,omitempty"`
}

type InstanceSpec struct {
	Name                 string        `json:"name,omitempty"`
	KeyName              string        `json:"key_name,omitempty"`
	ImageRef             string        `json:"imageRef,omitempty"`
	FlavorRef            string        `json:"flavorRef,omitempty"`
	MinCount             int           `json:"min_count,omitempty"`
	Networks             []Subnet      `json:"networks,omitempty"`
	BlockDeviceMappingV2 []BlockDevice `json:"block_device_mapping_v2,omitempty"`
}

type BlockDevice struct {
	UUID                string `json:"uuid,omitempty"`
	BootIndex           *int   `json:"boot_index,omitempty"`
	VolumeSize          int    `json:"volume_size,omitempty"`
	DeviceName          string `json:"device_name,omitempty"`
	SourceType          string `json:"source_type,omitempty"`
	DestinationType     string `json:"destination_type,omitempty"`
	DeleteOnTermination int    `json:"delete_on_termination,omitempty"`
}

type Subnet struct {
	Subnet string `json:"subnet,omitempty"`
}
