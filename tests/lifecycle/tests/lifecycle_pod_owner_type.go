package tests

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/config"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

var _ = Describe("lifecycle-pod-owner-type", func() {

	configSuite, err := config.NewConfig()
	if err != nil {
		glog.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	BeforeEach(func() {
		err := tshelper.WaitUntilClusterIsStable()
		Expect(err).ToNot(HaveOccurred())

		By("Clean namespace before each test")
		err = namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		By("Ensure all nodes are labeled with 'worker-cnf' label")
		err = nodes.EnsureAllNodesAreLabeled(globalhelper.GetAPIClient().CoreV1Interface, configSuite.General.CnfNodeLabel)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.LifecycleNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	// 47409
	It("One ReplicaSet, several pods", func() {

		By("Define ReplicaSet with replica number")
		replicaSet := tshelper.DefineReplicaSet(tsparams.TestReplicaSetName)
		replicaset.RedefineWithReplicaNumber(replicaSet, 3)

		err := globalhelper.CreateAndWaitUntilReplicaSetIsReady(replicaSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.TnfPodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodOwnerTypeTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47424
	It("Two deployments, several pods", func() {

		By("Define deployments")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentb, err := tshelper.DefineDeployment(2, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.TnfPodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodOwnerTypeTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47426
	It("StatefulSet pod", func() {

		By("Define statefulSet")
		statefulSet := tshelper.DefineStatefulSet(tsparams.TestStatefulSetName)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.TnfPodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodOwnerTypeTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47429
	It("One pod, not part of any workload resource [negative]", func() {

		By("Define pod")
		put := tshelper.DefinePod(tsparams.TestPodName)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.TnfPodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodOwnerTypeTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 47430
	It("Two deployments, one pod not related to any resource [negative]", func() {

		By("Define deployments")
		deploymenta, err := tshelper.DefineDeployment(2, 1, tsparams.TestDeploymentName)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymenta, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		deploymentb, err := tshelper.DefineDeployment(2, 1, "lifecycle-dpb")
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentb, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Define pod")
		put := tshelper.DefinePod(tsparams.TestPodName)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start lifecycle-pod-owner-type test")
		err = globalhelper.LaunchTests(tsparams.TnfPodOwnerTypeTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodOwnerTypeTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

})
