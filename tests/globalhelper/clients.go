package globalhelper

import (
	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	ctrlLogger "sigs.k8s.io/controller-runtime/pkg/log"
)

var egiClient *egiClients.Settings

func GetEcoGoinfraClient() *egiClients.Settings {
	if egiClient != nil {
		return egiClient
	}

	glog.Info("Creating new eco-goinfra k8s go-client with GinkgoLogr")

	ctrlLogger.SetLogger(GinkgoLogr)

	// ToDo: use kubeconfig from cli/config file.
	egiClient = egiClients.New("")

	return egiClient
}
