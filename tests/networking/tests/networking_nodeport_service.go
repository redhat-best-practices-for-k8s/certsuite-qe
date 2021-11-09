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

	// 45447
	It("2 custom pods, no service installed, service Should not have type of nodePort", func() {

	})

	// 45481
	It("2 custom pods, service installed without NodePort, service Should not have type of nodePort", func() {

	})

	// 45482
	It("2 custom pods, multiple services installed without NodePort, service Should not have type of nodePort",
		func() {

		})

	// 45483
	It("2 custom pods, service installed with NodePort, service Should not have type of nodePort [negative]",
		func() {

		})

	// 45484
	It("2 custom pods, multiple services installed and one has NodePort, service Should not have type of "+
		"nodePort [negative]", func() {

	})

})
