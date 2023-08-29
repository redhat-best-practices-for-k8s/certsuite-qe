package globalhelper

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"k8s.io/client-go/kubernetes"
)

var (
	apiclient *testclient.ClientSet
	conf      *config.Config
)

func SetTestK8sAPIClient(client kubernetes.Interface) {
	if apiclient == nil {
		apiclient = &testclient.ClientSet{}
	}

	//nolint:godox
	// TODO: Add more client interfaces as needed
	apiclient.CoreV1Interface = client.CoreV1()
	apiclient.AppsV1Interface = client.AppsV1()
	apiclient.RbacV1Interface = client.RbacV1()
	apiclient.NodeV1Interface = client.NodeV1()
	apiclient.K8sClient = client
}

func UnsetTestK8sAPIClient() {
	apiclient = nil
}

func GetAPIClient() *testclient.ClientSet {
	if apiclient != nil {
		return apiclient
	}

	var err error
	apiclient, err = config.DefineClients()

	if err != nil {
		glog.Fatal(fmt.Sprintf("can not load api client. Please check KUBECONFIG env var - %s", err))
	}

	return apiclient
}

func GetConfiguration() *config.Config {
	if conf != nil {
		return conf
	}

	var err error
	conf, err = config.NewConfig()

	if err != nil {
		glog.Fatal(fmt.Sprintf("can not load configuration - %s", err))
	}

	return conf
}

func GetOriginalTNFPaths() (string, string) {
	return GetConfiguration().General.TnfReportDir, GetConfiguration().General.TnfConfigDir
}

func OverrideDirectories(randomStr string) {
	reportDir := GetConfiguration().General.TnfReportDir + "/" + randomStr
	OverrideReportDir(reportDir)

	configDir := GetConfiguration().General.TnfConfigDir + "/" + randomStr
	OverrideTnfConfigDir(configDir)
}

func RestoreOriginalTNFPaths(reportDir, configDir string) {
	GetConfiguration().General.TnfReportDir = reportDir
	GetConfiguration().General.TnfConfigDir = configDir
}

func BeforeEachSetupWithRandomNamespace(incomingNamespace string) (randomNamespace, origReportDir, origConfigDir string) {
	randomNamespace = incomingNamespace + "-" + GenerateRandomString(10)

	By(fmt.Sprintf("Create %s namespace", randomNamespace))
	err := namespaces.Create(randomNamespace, GetAPIClient())
	Expect(err).ToNot(HaveOccurred())

	origReportDir, origConfigDir = GetOriginalTNFPaths()

	By("Override directories")
	OverrideDirectories(randomNamespace)

	return randomNamespace, origReportDir, origConfigDir
}

func DeleteTnfConfig() error {
	// Delete the tnf config directory
	return os.RemoveAll(GetConfiguration().General.TnfConfigDir)
}

func AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origConfigDir string, waitingTime time.Duration) {
	By("Remove reports from report directory")

	err := RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())

	By(fmt.Sprintf("Remove %s namespace", randomNamespace))
	err = namespaces.DeleteAndWait(
		GetAPIClient().CoreV1Interface,
		randomNamespace,
		waitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Delete TNF config")

	err = DeleteTnfConfig()
	Expect(err).ToNot(HaveOccurred())

	By("Restore directories")
	RestoreOriginalTNFPaths(origReportDir, origConfigDir)
}
