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
	var randomReportDir string
	var randomTnfConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and TNF config directories
		randomNamespace, randomReportDir, randomTnfConfigDir = globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNamespace)

		By("Define TNF config file")
		err := globalhelper.DefineTnfConfig(
			[]string{randomNamespace},
			tshelper.GetTnfTargetPodLabelsSlice(),
			[]string{},
			[]string{},
			[]string{tsparams.CrdSuffix1, tsparams.CrdSuffix2}, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace, randomReportDir, randomTnfConfigDir, tsparams.CrdDeployTimeoutMins)
	})

	// 56635
	It("One deployment, pod disruption budget minAvailable value meet requirements", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has one replica")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(1)))

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56636
	It("One deployment, pod disruption budget maxUnavailable value meet requirements", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has two replicas")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(2)))

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56637
	It("One statefulSet, pod disruption budget minAvailable value is zero [negative]", func() {
		By("Create statefulSet")
		myStatefulSet := statefulset.DefineStatefulSet(tsparams.TestStatefulSetBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		statefulset.RedefineWithReplicaNumber(myStatefulSet, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(myStatefulSet, tsparams.StatefulSetDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert statefulSet has one replica")
		runningStatefulSet, err := globalhelper.GetRunningStatefulSet(myStatefulSet.Namespace, myStatefulSet.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningStatefulSet.Spec.Replicas).To(Equal(int32(1)))

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(0), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56638
	It("One deployment, pod disruption budget maxUnavailable equals to replica number [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has two replicas")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(2)))

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(2), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56746
	It("One deployment, pod disruption budget maxUnavailable is bigger than the replica number [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 2)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has two replicas")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(2)))

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMaxUnAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(3), tsparams.TnfTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment, pod disruption budget matchLabels does not match deployment label [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has one replica")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(1)))

		By("Create pod disruption budget")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.TnfUnknownPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment, no pod disruption budget [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has one replica")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(1)))

		By("Start TNF " + tsparams.TnfPodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.TnfPodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomTnfConfigDir)
		Expect(err).To(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfPodDisruptionBudgetTcName, globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})
