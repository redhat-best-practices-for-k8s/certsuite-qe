package tests

import (
	. "github.com/onsi/ginkgo"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Networking custom namespace, custom deployment,", func() {

	execute.BeforeAll(func() {

	})

	BeforeEach(func() {

	})

	// 45440
	It("3 custom pods on Default network networking-icmpv4-connectivity", func() {

	})

	// 45441
	It("custom daemonset, 4 custom pods on Default network", func() {

	})

	// 45442
	It("3 custom pods on Default network networking-icmpv4-connectivity fail when one pod is "+
		"disconnected [negative]", func() {

	})

	// 45443
	It("2 custom pods on Default network networking-icmpv4-connectivity fail when there is no ping binary "+
		"[negative]", func() {

	})

	// 45444
	It("2 custom pods on Default network networking-icmpv4-connectivity skip when label "+
		"test-network-function.com/skip_connectivity_tests is set in deployment [skip]", func() {

	})

	// 45445
	It("custom daemonset, 4 custom pods on Default network networking-icmpv4-connectivity pass when label "+
		"test-network-function.com/skip_connectivity_tests is set in deployment only", func() {

	})

	// 45446
	It("2 custom pods on Default network networking-icmpv4-connectivity skip when there is no ip binary [skip]",
		func() {

		})

})
