package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/poddisruptionbudget"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe(tsparams.TnfPodDisruptionBudgetTcName, func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.TestNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Clean namespace after each test")
		err := namespaces.Clean(tsparams.TestNamespace, globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())
	})

	const tnfTestCaseName = tsparams.TnfPodDisruptionBudgetTcName

	// 56635
	It("One deployment, pod disruption budget minAvailable value meet requirements", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, tsparams.TestNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, tsparams.TestNamespace,
			intstr.FromInt(1), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56636
	It("One deployment, pod disruption budget maxUnavailable value meet requirements", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, tsparams.TestNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, tsparams.TestNamespace,
			intstr.FromInt(1), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56637
	It("One statefulSet, pod disruption budget minAvailable value is zero [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create statefulSet")
		sf := statefulset.DefineStatefulSet(tsparams.TestStatefulSetBaseName, tsparams.TestNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		statefulset.RedefineWithReplicaNumber(sf, 1)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(sf, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, tsparams.TestNamespace,
			intstr.FromInt(0), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56638
	It("One deployment, pod disruption budget maxUnavailable equals to replica number [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, tsparams.TestNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, tsparams.TestNamespace,
			intstr.FromInt(2), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56746
	It("One deployment, pod disruption budget maxUnavailable is bigger than the replica number [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, tsparams.TestNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, tsparams.TestNamespace,
			intstr.FromInt(3), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tnfTestCaseName + " test case")
		err = globalhelper.LaunchTests(tnfTestCaseName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tnfTestCaseName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})
