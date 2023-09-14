package tests

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
	The setup needs to be as mentioned here :
	https://developers.redhat.com/articles/2021/08/27/using-virtual-functions-dpdk-red-hat-openshift

	Assumption : One DPDK pod is running in default
*/

var _ = Describe("Networking dpdk-cpu-pinning-exec-probe,", func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	var dpdkPod *v1.Pod

	BeforeEach(func() {

		By("Check if DPDK pod is running")
		foundDpdkPod, err := globalhelper.GetAPIClient().Pods("default").Get(context.TODO(), "dpdk-pod", metav1.GetOptions{})
		if k8serrors.IsNotFound(err) {
			Skip("Setup for running DPPDK is not satisfied. Hence, skipping")
		}
		Expect(err).ToNot(HaveOccurred())
		dpdkPod = foundDpdkPod

		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNetworkingNameSpace)

		By("Define TNF config file")
		err = globalhelper.DefineTnfConfig(
			[]string{"default"},
			[]string{randomNamespace},
			[]string{},
			[]string{},
			[]string{})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.WaitingTime)
	})

	It("one dpdk pod with exec probe", func() {

		By("Deploy dpdk pod in test namespace")
		dpdkPod.Namespace = randomNamespace

		_, err := globalhelper.GetAPIClient().Pods(randomNamespace).Create(context.TODO(),
			dpdkPod, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDpdkCPUPinningExecProbe,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDpdkCPUPinningExecProbe,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	/*It("one dpdk pod with no exec probe", func() {

		By("Check if DPDK pod is running")
		dpdk, err := globalhelper.GetAPIClient().Pods("default").Get(context.TODO(), "dpdk-pod", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		By("Check if liveness probe exists")
		Expect(dpdk.Spec.Containers[0].LivenessProbe).NotTo(nil)

		dpdk.Spec.Containers[0].LivenessProbe = nil

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.TnfDpdkCPUPinningExecProbe,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfDpdkCPUPinningExecProbe,
			globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

	})*/
})
