package tests

import (
	. "github.com/onsi/ginkgo"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/execute"
)

var _ = Describe("Affiliated-certification container certification,", func() {

	execute.BeforeAll(func() {

	})

	BeforeEach(func() {

	})

	// 46562
	It("one container to test, container is certified", func() {

	})

	// 46563
	It("one container to test, container is not certified [negative]", func() {

	})

	// 46564
	It("two containers to test, both are certified", func() {

	})

	// 46565
	It("two containers to test, one is certified, one is not [negative]", func() {

	})

	// 46566
	It("certifiedcontainerinfo field exists in tnf_config but has no value [skip]", func() {

	})

	// 46567
	It("certifiedcontainerinfo field does not exist in tnf_config [skip]", func() {

	})

	// 46578
	It("name and repository fields exist in certifiedcontainerinfo field but are empty [skip]", func() {

	})

	// 46579
	It("name field in certifiedcontainerinfo field is populated but repository field is not [skip]", func() {

	})

	// 46580
	It("repository field in certifiedcontainerinfo field is populated but name field is not [skip]", func() {

	})

	// 46581
	It("two containers listed in certifiedcontainerinfo field, one is certified, one has empty name and "+
		"repository fields", func() {

	})

})
