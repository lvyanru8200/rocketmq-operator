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

package main

import (
	"flag"
	"fmt"
	"os"
	"rocketmq-operator-v2/pkg/logi"

	"github.com/open-policy-agent/cert-controller/pkg/rotator"
	"k8s.io/apimachinery/pkg/types"

	"rocketmq-operator-v2/pkg/controller/common"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	rocketmqv1 "rocketmq-operator-v2/api/v1"
	"rocketmq-operator-v2/controllers"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var (
	log      = logi.GetSugaredLogger()
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
	dnsName  = fmt.Sprintf("%s.%s.svc", serviceName, common.GetOperatorNamespace())
	webhooks = []rotator.WebhookInfo{
		{
			Name: "rocketmq-operator-validating-webhook-config",
			Type: rotator.Validating,
		},
		{
			Name: "rocketmq-operator-mutating-webhook-config",
			Type: rotator.Mutating,
		},
	}
)

const (
	secretName     = "rocketmq-operator-webhook-cert"
	serviceName    = "rocketmq-operator-webhook"
	caName         = "rocketmq-operator-ca"
	caOrganization = "rocketmq-operator"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(rocketmqv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var disableCertRotation bool
	var certDir string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&certDir, "cert-dir", "/certs", "The directory where certs are stored, defaults to /certs")
	flag.BoolVar(&disableCertRotation, "disable-cert-rotation", false, "disable automatic generation and rotation of webhook TLS certificates/keys")
	flag.Parse()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "5c4daf29.daocloud.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	setupFinished := make(chan struct{})
	if !disableCertRotation {
		setupLog.Info("setting up cert rotation")
		err := rotator.AddRotator(mgr, &rotator.CertRotator{
			SecretKey: types.NamespacedName{
				Namespace: common.GetOperatorNamespace(),
				Name:      secretName,
			},
			CertDir:        certDir,
			CAName:         caName,
			CAOrganization: caOrganization,
			DNSName:        dnsName,
			IsReady:        setupFinished,
			Webhooks:       webhooks,
		})
		if err != nil {
			setupLog.Error(err, "unable to set up cert rotation")
			os.Exit(1)
		}
	} else {
		close(setupFinished)
	}

	go func() {
		<-setupFinished

		if err = (&controllers.DledgerBrokerReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "DledgerBroker")
			os.Exit(1)
		}
		if err = (&controllers.NameserverReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Nameserver")
			os.Exit(1)
		}

		if certDir != "" {
			if err = (&rocketmqv1.DledgerBroker{}).SetupWebhookWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create webhook", "webhook", "DledgerBroker")
				os.Exit(1)
			}
			if err = (&rocketmqv1.Nameserver{}).SetupWebhookWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create webhook", "webhook", "Nameserver")
				os.Exit(1)
			}
		}
	}()
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
