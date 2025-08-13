package globalhelper

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	testclient "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/client"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
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
		glog.Fatalf("can not load api client. Please check KUBECONFIG env var - %s", err)
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
		glog.Fatalf("can not load configuration - %s", err)
	}

	// Apply optional memory/GC tuning to the current test process (with defaults)
	if conf.General.GoMemLimit != "" {
		_ = os.Setenv("GOMEMLIMIT", conf.General.GoMemLimit)
	}
	if conf.General.GoGC != "" {
		_ = os.Setenv("GOGC", conf.General.GoGC)
	} else if os.Getenv("GOGC") == "" {
		// Default to a moderately conservative GC to reduce memory footprint
		_ = os.Setenv("GOGC", "75")
	}

	return conf
}

func GenerateDirectories(randomStr string) (reportDir, configDir string) {
	reportDir = GetConfiguration().General.CertsuiteReportDir + "/" + randomStr
	configDir = GetConfiguration().General.CertsuiteConfigDir + "/" + randomStr

	err := os.MkdirAll(reportDir, os.ModePerm)
	if err != nil {
		glog.Error("could not create dest directory= %s, err=%s", reportDir, err)
	}

	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		glog.Error("could not create dest directory= %s, err=%s", reportDir, err)
	}

	return reportDir, configDir
}

func BeforeEachSetupWithRandomNamespace(incomingNamespace string) (randomNamespace, randomReportDir, randomConfigDir string) {
	randomNamespace = incomingNamespace + "-" + GenerateRandomString(10)

	By(fmt.Sprintf("Create %s namespace", randomNamespace))
	err := CreateNamespace(randomNamespace)
	Expect(err).ToNot(HaveOccurred())

	By("Generate directories")

	randomReportDir, randomConfigDir = GenerateDirectories(randomNamespace)

	return randomNamespace, randomReportDir, randomConfigDir
}

// BeforeEachSetupWithRandomPrivilegedNamespace creates a random namespace with privileged Pod Security Standards
// for tests that require hostIPC, hostPID, or other privileged operations (like access control tests).
func BeforeEachSetupWithRandomPrivilegedNamespace(incomingNamespace string) (randomNamespace, randomReportDir,
	randomConfigDir string) {
	randomNamespace = incomingNamespace + "-" + GenerateRandomString(10)

	By(fmt.Sprintf("Create %s privileged namespace", randomNamespace))
	err := CreatePrivilegedNamespace(randomNamespace)
	Expect(err).ToNot(HaveOccurred())

	By("Generate directories")

	randomReportDir, randomConfigDir = GenerateDirectories(randomNamespace)

	return randomNamespace, randomReportDir, randomConfigDir
}

func AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomConfigDir string, waitingTime time.Duration) {
	// logfile := "certsuite.log"
	// By("Print logs")
	// myFile, err := os.ReadFile(randomReportDir + "/" + logfile)
	// if err != nil {
	// 	glog.Errorf("can not read file %s - %s", logfile, err)
	// }
	// fmt.Println(string(myFile))
	// By(fmt.Sprintf("Remove reports from report directory: %s", randomReportDir))
	err := RemoveContentsFromReportDir(randomReportDir)
	Expect(err).ToNot(HaveOccurred())

	By(fmt.Sprintf("Remove configs from config directory: %s", randomConfigDir))

	err = RemoveContentsFromConfigDir(randomConfigDir)
	Expect(err).ToNot(HaveOccurred())

	By(fmt.Sprintf("Remove %s namespace", randomNamespace))
	err = DeleteNamespaceAndWait(randomNamespace, waitingTime)
	Expect(err).ToNot(HaveOccurred())
}
