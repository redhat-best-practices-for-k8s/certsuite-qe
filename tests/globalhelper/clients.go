package globalhelper

import (
	. "github.com/onsi/ginkgo/v2"
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	klog "k8s.io/klog/v2"
	ctrlLogger "sigs.k8s.io/controller-runtime/pkg/log"
)

var egiClient *egiClients.Settings

func GetEcoGoinfraClient() *egiClients.Settings {
	if egiClient != nil {
		return egiClient
	}

	klog.Info("Creating new eco-goinfra k8s go-client with GinkgoLogr")

	ctrlLogger.SetLogger(GinkgoLogr)

	// ToDo: use kubeconfig from cli/config file.
	egiClient = egiClients.New("")

	return egiClient
}
