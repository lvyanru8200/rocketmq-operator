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
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"rocketmq-operator-v2/pkg/logi"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	rocketmqv1 "rocketmq-operator-v2/api/v1"
)

var log = logi.GetSugaredLogger()
var dledgerBrokerFinalizerName = "dledgerbroker.finalizers.rocketmq.daocloud.io"

// DledgerBrokerReconciler reconciles a DledgerBroker object
type DledgerBrokerReconciler struct {
	client.Client
	Log    *zap.SugaredLogger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=rocketmq.daocloud.io,resources=dledgerbrokers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rocketmq.daocloud.io,resources=dledgerbrokers/status,verbs=get;update;patch

func (r *DledgerBrokerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLog := log.With(
		zap.String("Request.Namespace", req.Namespace),
		zap.String("Request.Name", req.Name),
	)
	reqLog.Info("Reconcile DledgerModel Broker")
	r.Log = reqLog
	instance := &rocketmqv1.DledgerBroker{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(instance.GetFinalizers(), dledgerBrokerFinalizerName) {
			controllerutil.AddFinalizer(instance, dledgerBrokerFinalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(instance.GetFinalizers(), dledgerBrokerFinalizerName) {
			controllerutil.RemoveFinalizer(instance, dledgerBrokerFinalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	//

	return ctrl.Result{}, nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func (r *DledgerBrokerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rocketmqv1.DledgerBroker{}).
		Complete(r)
}
