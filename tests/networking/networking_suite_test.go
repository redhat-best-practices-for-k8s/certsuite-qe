//go:build !utest

package networking

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	klog "k8s.io/klog/v2"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/tests"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/cluster"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/parameters"
)

func TestNetworking(t *testing.T) {
	globalhelper.RunSuite(t, "CNFCert networking tests")
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

	By("Set rbac policy which allows authenticated users to run privileged containers")
	err = globalhelper.AllowAuthenticatedUsersRunPrivilegedContainers()
	Expect(err).ToNot(HaveOccurred())

	By("Ensure all nodes are labeled with 'worker-cnf' label")
	err = nodes.EnsureAllNodesAreLabeled(configSuite.General.CnfNodeLabel)
	Expect(err).ToNot(HaveOccurred())
}, func() {})
