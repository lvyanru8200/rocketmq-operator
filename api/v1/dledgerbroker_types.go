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

// DledgerBrokerSpec defines the desired state of DledgerBroker
type DledgerBrokerSpec struct {
	Dledger            `json:",inline"`             // dledger模式配置
	Resource           *corev1.ResourceRequirements `json:"resource,omitempty"` // dledger模式pod配置
	Storage            *DledgerStorage              `json:"storage,omitempty"`  // 存储设置
	Export             *ExportSetting               `json:"export,omitempty"`   // 监控设置
	ImageSetting       `json:",inline"`             // broker镜像设置
	ServiceAccountName string                       `json:"serviceAccountName,omitempty"` // serviceaccount名称
	PodSpec            *PodSpec                     `json:"podSpec,omitempty"`            // broker pod配置
	Env                []corev1.EnvVar              `json:"env,omitempty"`                // 环境变量设置
	Config             map[string]string            `json:"config,omitempty"`             // broker 配置文件
	Nameserver         string                       `json:"nameserver,omitempty"`         // 需要连接的nameserver实例名称
	Acl                *Acl                         `json:"acl,omitempty"`                // broker acl配置
}

// Dledger模式设置
type Dledger struct {
	BrokerGroupNumber    int   `json:"brokerGroupNumber,omitempty"`
	BrokerNumberPerGroup []int `json:"brokerNumberPerGroup,omitempty"` // broker每个group node数量
}

// 存储设置
type DledgerStorage struct {
	StorageClass string `json:"storageClass,omitempty"`
	Size         string `json:"size,omitempty"`
}

// export设置
type ExportSetting struct {
	Open         bool                         `json:"open"`
	Resource     *corev1.ResourceRequirements `json:"resource,omitempty"`
	ImageSetting `json:",inline"`
}

// image设置
type ImageSetting struct {
	Image           string                        `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy             `json:"imagePullPolicy,omitempty"`
	ImagePullSecret []corev1.LocalObjectReference `json:"imagePullSecret,omitempty"`
}

// pod设置
type PodSpec struct {
	HostAliases     []corev1.HostAlias         `json:"hostAliases,omitempty" patchStrategy:"merge" patchMergeKey:"ip" protobuf:"bytes,23,rep,name=hostAliases"`
	RestartPolicy   corev1.RestartPolicy       `json:"restartPolicy,omitempty" protobuf:"bytes,3,opt,name=restartPolicy,casttype=RestartPolicy"`
	NodeSelector    map[string]string          `json:"nodeSelector,omitempty" protobuf:"bytes,7,rep,name=nodeSelector"`
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty" protobuf:"bytes,14,opt,name=securityContext"`
	Affinity        *corev1.Affinity           `json:"affinity,omitempty" protobuf:"bytes,18,opt,name=affinity"`
	Tolerations     []corev1.Toleration        `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
}

// acl设置
type Acl struct {
	GlobalWhiteRemoteAddresses []string  `json:"globalWhiteRemoteAddresses"`
	Accounts                   []Account `json:"accounts"`
}

type Account struct {
	AccessKey          string   `json:"accessKey"`
	SecretKey          string   `json:"secretKey"`
	WhiteRemoteAddress string   `json:"whiteRemoteAddress,omitempty"`
	Admin              bool     `json:"admin,omitempty"`
	DefaultTopicPerm   string   `json:"defaultTopicPerm"`
	DefaultGroupPerm   string   `json:"defaultGroupPerm"`
	TopicPerms         []string `json:"topicPerms"`
	GroupPerms         []string `json:"groupPerms"`
}

// DledgerBrokerStatus defines the observed state of DledgerBroker
type DledgerBrokerStatus struct {
	BrokerConfigmap string              `json:"brokerConfigmap,omitempty"` // 当前实例挂载的broker配置
	NameserverAddr  []string            `json:"nameserverAddr,omitempty"`  // 当前实例上报的nameserver地址
	InternalAccess  string              `json:"InternalAccess,omitempty"`  // 内部访问地址
	ExternalAccess  string              `json:"ExternalAccess,omitempty"`  // 外部访问地址
	BrokerInfo      map[string][]string `json:"BrokerInfo,omitempty"`      // broker配置信息
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DledgerBroker is the Schema for the dledgerbrokers API
type DledgerBroker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DledgerBrokerSpec   `json:"spec,omitempty"`
	Status DledgerBrokerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DledgerBrokerList contains a list of DledgerBroker
type DledgerBrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DledgerBroker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DledgerBroker{}, &DledgerBrokerList{})
}
