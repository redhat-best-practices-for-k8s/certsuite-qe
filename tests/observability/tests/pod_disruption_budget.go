package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/observability/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/poddisruptionbudget"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe(tsparams.TnfPodDisruptionBudgetTcName, func() {
	var randomNamespace string
	var origReportDir string
	var origTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, origReportDir, origTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			tshelper.GetTnfTargetPodLabelsSlice(),
			[]string{},
			[]string{},
			[]string{tsparams.CrdSuffix1, tsparams.CrdSuffix2})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, origReportDir, origTnfConfigDir, tsparams.CrdDeployTimeoutMins)
	})

	// 56635
	It("One deployment, pod disruption budget minAvailable value meet requirements", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56636
	It("One deployment, pod disruption budget maxUnavailable value meet requirements", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName, qeTcFileName)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56637
	It("One statefulSet, pod disruption budget minAvailable value is zero [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create statefulSet")
		sf := statefulset.DefineStatefulSet(tsparams.TestStatefulSetBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		statefulset.RedefineWithReplicaNumber(sf, 1)

		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(sf, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(0), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56638
	It("One deployment, pod disruption budget maxUnavailable equals to replica number [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(2), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56746
	It("One deployment, pod disruption budget maxUnavailable is bigger than the replica number [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(3), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment, pod disruption budget matchLabels does not match deployment label [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.TnfUnknownPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment, no pod disruption budget [negative]", func() {
		qeTcFileName := globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText())

		By("Create deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(globalhelper.GetAPIClient().K8sClient.AppsV1(),
			dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName, qeTcFileName)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Junit and Claim reports")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())
	})
})
