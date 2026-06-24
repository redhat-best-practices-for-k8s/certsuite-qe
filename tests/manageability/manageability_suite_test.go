//go:build !utest

package performance

import (
	"testing"

	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/manageability/tests"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
)

func TestManageability(t *testing.T) {
	globalhelper.RunSuite(t, "CNFCert performance tests")
}
