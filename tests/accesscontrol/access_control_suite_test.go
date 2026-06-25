//go:build !utest

package accesscontrol

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/tests"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
)

func TestAccessControl(t *testing.T) {
	globalhelper.RunSuite(t, "CNFCert access-control tests")
}

var _ = SynchronizedBeforeSuite(func() {
	err := globalhelper.AllowAuthenticatedUsersRunPrivilegedContainers()
	Expect(err).ToNot(HaveOccurred(), "Error creating namespace")
}, func() {})
