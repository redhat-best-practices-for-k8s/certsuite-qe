//go:build !utest

package performance

import (
	"testing"

	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/tests"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
)

func TestPerformance(t *testing.T) {
	globalhelper.RunSuite(t, "CNFCert performance tests")
}
