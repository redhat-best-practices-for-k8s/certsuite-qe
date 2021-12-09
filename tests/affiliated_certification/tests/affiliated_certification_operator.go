package tests

import (
	. "github.com/onsi/ginkgo"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Affiliated-certification operator certification,", func() {

	execute.BeforeAll(func() {

	})

	BeforeEach(func() {

	})

	// 46582
	It("one operator to test, operator is certified", func() {

	})

	// 46695
	It("one operator to test, operator is not certified [negative]", func() {

	})

	// 46696
	It("two operators to test, both are certified", func() {

	})

	// 46697
	It("two operators to test, one is certified, one is not [negative]", func() {

	})

	// 46698
	It("certifiedoperatorinfo field exists in tnf_config but has no value [skip]", func() {

	})

	// 46699
	It("certifiedoperatorinfo field does not exist in tnf_config [skip]", func() {

	})

	// 46700
	It("name and organization fields exist in certifiedoperatorinfo but are empty [skip]", func() {

	})

	// 46702
	It("name field in certifiedoperatorinfo field is populated but organization field is not [skip]", func() {

	})

	// 46704
	It("organization field in certifiedoperatorinfo field is populated but name field is not [skip]", func() {

	})

	// 46706
	It("two operators to test, one is certified, one has empty name and organization fields", func() {

	})

})
