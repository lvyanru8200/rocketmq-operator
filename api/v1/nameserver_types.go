/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NameserverSpec defines the desired state of Nameserver  NameserverSpec定义Nameserver的所需状态
type NameserverSpec struct {
	NameserverNumber   int                         `json:"nameserverNumber,omitempty"`
	Resource           corev1.ResourceRequirements `json:"resource,omitempty"`
	Image              ImageSetting                `json:"image,omitempty"`
	ServiceAccountName string                      `json:"serviceAccountName,omitempty" protobuf:"bytes,8,opt,name=serviceAccountName"`
	Env                []corev1.EnvVar             `json:"env,omitempty"`
	PodSpec            PodSpec                     `json:"podSpec,omitempty"`
	Export             ExportSetting               `json:"export,omitempty"`
}

// NameserverStatus defines the observed state of Nameserver  NameserverStatus定义观察到的Nameserver状态
type NameserverStatus struct {
	ConnectAddr string `json:"externalAddr,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Nameserver is the Schema for the nameservers API  Nameserver是nameserversapi的模式
type Nameserver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NameserverSpec   `json:"spec,omitempty"`
	Status NameserverStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NameserverList contains a list of Nameserver
type NameserverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nameserver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nameserver{}, &NameserverList{})
}
