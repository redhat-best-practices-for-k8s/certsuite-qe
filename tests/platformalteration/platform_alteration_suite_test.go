//go:build !utest

package platformalteration

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	klog "k8s.io/klog/v2"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/tests"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/cluster"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"

	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
)

func TestPlatformAlteration(t *testing.T) {
	globalhelper.RunSuite(t, "CNFCert platform-alteration tests")
}

var _ = SynchronizedBeforeSuite(func() {
	configSuite, err := config.NewConfig()
	if err != nil {
		klog.Fatalf("can not load config file: %v", err)
	}

	By("Validate that cluster is Schedulable")
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, tsparams.WaitingTime, tsparams.RetryInterval*time.Second).Should(BeTrue())

	By("Validate that all nodes are Ready")
	err = nodes.WaitForNodesReady(globalhelper.GetAPIClient().Nodes(), tsparams.WaitingTime, tsparams.RetryInterval)
	Expect(err).ToNot(HaveOccurred())

	By("Ensure all nodes are labeled with 'worker-cnf' label")
	err = nodes.EnsureAllNodesAreLabeled(configSuite.General.CnfNodeLabel)
	Expect(err).ToNot(HaveOccurred())

	By("Set rbac policy which allows authenticated users to run privileged containers")
	err = globalhelper.AllowAuthenticatedUsersRunPrivilegedContainers()
	Expect(err).ToNot(HaveOccurred())
}, func() {})
