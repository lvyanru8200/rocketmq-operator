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
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"rocketmq-operator-v2/pkg/configs"
	"rocketmq-operator-v2/pkg/logi"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strconv"
)

// log is for logging in this package.
var dledgerbrokerlog = logi.GetSugaredLogger().With(zap.String("Webhook", "Dledgerbroker"))

func (r *DledgerBroker) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-rocketmq-daocloud-io-v1-dledgerbroker,mutating=true,failurePolicy=fail,groups=rocketmq.daocloud.io,resources=dledgerbrokers,verbs=create;update,versions=v1,name=mdledgerbroker.kb.io

var _ webhook.Defaulter = &DledgerBroker{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *DledgerBroker) Default() {
	dledgerbrokerlog.Info("default", "name", r.Name)
	brokerGroups := r.Spec.Dledger.BrokerGroupNumber
	preGroup := r.Spec.BrokerNumberPerGroup
	cfg := configs.GetGlobalConfig()
	specEnv := r.Spec.Env
	defer func() {
		r.Spec.Env = specEnv
	}()
	if brokerGroups <= 0 {
		brokerGroups = 2
	}
	if len(preGroup) == 0 {
		preGroup = make([]int, 0, brokerGroups)
		for k := range preGroup {
			preGroup[k] = 3
		}
		r.Spec.BrokerNumberPerGroup = preGroup
	}
	if r.Spec.Image == "" {
		r.Spec.Image = cfg.IMAGE_ROCKETMQ
	}

	if r.Spec.Resource == nil || r.Spec.Resource.Size() == 0 {
		r.Spec.Resource = new(v1.ResourceRequirements)
		*r.Spec.Resource = defaultBrokerResource()
	}

	if r.Spec.Export.Open {
		if r.Spec.Export.Image == "" {
			r.Spec.Export.Image = cfg.IMAGE_EXPORTER
		}
		if r.Spec.Export.Resource == nil || r.Spec.Export.Resource.Size() == 0 {
			r.Spec.Export.Resource = new(v1.ResourceRequirements)
			*r.Spec.Export.Resource = defaultExportResource()
		}
	}

	if r.Spec.Storage.StorageClass == "" {
		r.Spec.Storage.StorageClass = cfg.STORAGE_CLASS_NAME
	}

	if r.Spec.Storage.Size == "" {
		r.Spec.Storage.Size = "2Gi"
	}
	func() {
		m := r.Spec.Resource.Requests.Memory().Size()
		if m <= 0 {
			return
		}
		specEnv = configs.SetEnvIfUnset(specEnv, configs.JVM_XMX, strconv.FormatInt(int64(m), 10)+"m")
		specEnv = configs.SetEnvIfUnset(specEnv, configs.JVM_XMS, strconv.FormatInt(int64(m)/4, 10)+"m")
	}()

	specEnv = configs.MergeEnv(specEnv, cfg.InstanceEnv)
}

func defaultBrokerResource() v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("500m"),
			v1.ResourceMemory: resource.MustParse("1Gi"),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("1000m"),
			v1.ResourceMemory: resource.MustParse("2Gi"),
		},
	}
}

func defaultExportResource() v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("200m"),
			v1.ResourceMemory: resource.MustParse("0.2Gi"),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("200m"),
			v1.ResourceMemory: resource.MustParse("0.2Gi"),
		},
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-rocketmq-daocloud-io-v1-dledgerbroker,mutating=false,failurePolicy=fail,groups=rocketmq.daocloud.io,resources=dledgerbrokers,versions=v1,name=vdledgerbroker.kb.io

var _ webhook.Validator = &DledgerBroker{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *DledgerBroker) ValidateCreate() error {
	dledgerbrokerlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *DledgerBroker) ValidateUpdate(old runtime.Object) error {
	dledgerbrokerlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DledgerBroker) ValidateDelete() error {
	dledgerbrokerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
