//go:build !utest

package operator

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/operator/tests"
)

func TestOperator(t *testing.T) {
	globalhelper.RunSuite(t, "CNFCert operator tests")
}

var _ = SynchronizedBeforeSuite(func() {
	if globalhelper.IsKindCluster() {
		Skip("Skipping operator tests on kind cluster")
	}

	// Safeguard against running the operator tests on a cluster without catalog sources
	if !globalhelper.IsKindCluster() {
		By("Create catalog sources and wait for them to become ready")
		err := globalhelper.CreateAndValidateCatalogSources(true)
		Expect(err).ToNot(HaveOccurred(), "All necessary catalog sources are not available")
	}
}, func() {})
