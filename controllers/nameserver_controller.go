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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	"rocketmq-operator-v2/api/v1"
	rocketmqv1 "rocketmq-operator-v2/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NameserverReconciler reconciles a Nameserver object
type NameserverReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	v1.Nameserver
}

// +kubebuilder:rbac:groups=rocketmq.daocloud.io,resources=nameservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rocketmq.daocloud.io,resources=nameservers/status,verbs=get;update;patch

//以deployment的形式创建namesever，并且为创建的这个nameserver创建一个无头服务

func (r *NameserverReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("nameserver", req.NamespacedName)

	var NameSever v1.Nameserver
	if err := r.Get(context, req.NamespacedName, &NameSever); err != nil {
		log.Error(err, "unable to fetch NameSever")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploy := &appsv1.Deployment{
		TypeMeta:   NameSever.TypeMeta,
		ObjectMeta: NameSever.ObjectMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(rune(r.Nameserver.Spec.NameserverNumber)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "rocketmq-nameSvr",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nameSever",
					Labels: map[string]string{
						"app":     "rocketmq-nameSvr",
						"version": "v1",
					},
				},
				Spec: corev1.PodSpec{
					HostAliases:      NameSever.Spec.PodSpec.HostAliases,
					RestartPolicy:    NameSever.Spec.PodSpec.RestartPolicy,
					NodeSelector:     NameSever.Spec.PodSpec.NodeSelector,
					SecurityContext:  NameSever.Spec.PodSpec.SecurityContext,
					Affinity:         NameSever.Spec.PodSpec.Affinity,
					Tolerations:      NameSever.Spec.PodSpec.Tolerations,
					ImagePullSecrets: NameSever.Spec.Image.ImagePullSecret,
					Containers: []corev1.Container{
						{
							Name:            "rocketmq-nameSrv",
							Image:           NameSever.Spec.Image.Image,
							ImagePullPolicy: NameSever.Spec.Image.ImagePullPolicy,
						},
					},
				},
			},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    pointer.Int32Ptr(10),
			Paused:                  false,
			ProgressDeadlineSeconds: pointer.Int32Ptr(600),
		},
	}

	if err := r.Client.Create(context, deploy, &client.CreateOptions{DryRun: []string{"666"}}); err != nil {
		log.Error(err, "unable to create deployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var NameSeverLists v1.NameserverList
	if err := r.List(context, &NameSeverLists, client.InNamespace(req.Namespace), client.MatchingFields{"nameSever": req.Name}); err != nil {
		log.Error(err, "unable to list nameSevers")
		return ctrl.Result{}, err
	}

	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: NameSever.Spec.ServiceAccountName,
			Labels: map[string]string{
				"apps": "nameSrv-headless",
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(deploy, schema.GroupVersionKind{
					Group: "rocketmq.daocloud.io",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeNodePort,
			ClusterIP: "none",
			Selector: map[string]string{
				"app": "rocketmq-nameSvc",
			},
		},
	}

	if err := r.Client.Create(context, svc, &client.CreateOptions{DryRun: []string{"666"}}); err != nil {
		log.Error(err, "unable to create service")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return ctrl.Result{}, nil
}

func (r *NameserverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rocketmqv1.Nameserver{}).
		Complete(r)
}
