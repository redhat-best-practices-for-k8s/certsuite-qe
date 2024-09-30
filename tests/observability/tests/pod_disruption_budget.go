package tests

import (
	"maps"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/observability/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/observability/parameters"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/poddisruptionbudget"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/statefulset"
	"k8s.io/apimachinery/pkg/util/intstr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe(tsparams.CertsuitePodDisruptionBudgetTcName, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNamespace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			tshelper.GetCertsuiteTargetPodLabelsSlice(),
			[]string{},
			[]string{},
			[]string{tsparams.CrdSuffix1, tsparams.CrdSuffix2}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.CrdDeployTimeoutMins)
	})

	// 56635
	It("One deployment, pod disruption budget minAvailable value meet requirements", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

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
			intstr.FromInt(1), tsparams.CertsuiteTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56636
	It("One deployment, pod disruption budget maxUnavailable value meet requirements", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

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
			intstr.FromInt(1), tsparams.CertsuiteTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56637
	It("One statefulSet, pod disruption budget minAvailable value is zero [negative]", func() {
		By("Create statefulSet")
		myStatefulSet := statefulset.DefineStatefulSet(tsparams.TestStatefulSetBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

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
			intstr.FromInt(0), tsparams.CertsuiteTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56638
	It("One deployment, pod disruption budget maxUnavailable equals to replica number [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

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
			intstr.FromInt(2), tsparams.CertsuiteTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56746
	It("One deployment, pod disruption budget maxUnavailable is bigger than the replica number [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

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
			intstr.FromInt(3), tsparams.CertsuiteTargetPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment, pod disruption budget matchLabels does not match deployment label [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

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
			intstr.FromInt(1), tsparams.CertsuiteUnknownPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment, no pod disruption budget [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has one replica")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(1)))

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment with two labels that match the labels present in the pod disruption budget matchLabels", func() {
		By("Define deployment with an extra label")
		extraPodLabels := maps.Clone(tsparams.CertsuiteTargetPodLabels)
		extraPodLabels[tsparams.TestPodExtraLabelKey] = tsparams.TestPodExtraLabelValue
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, extraPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has one replica")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(1)))

		By("Create pod disruption budget with an extra label")
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), extraPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment that misses one label present in the pod disruption budget matchLabels [negative]", func() {
		By("Define deployment")
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, tsparams.CertsuiteTargetPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has one replica")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(1)))

		By("Create pod disruption budget with an extra label")
		extraPodLabels := maps.Clone(tsparams.CertsuiteTargetPodLabels)
		extraPodLabels[tsparams.TestPodExtraLabelKey] = tsparams.TestPodExtraLabelValue
		pdb := poddisruptionbudget.DefinePodDisruptionBudgetMinAvailable(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), extraPodLabels)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment with two labels that match the labels in the pod disruption budget matchLabels and matchExpressions", func() {
		By("Define deployment with an extra label")
		extraPodLabels := maps.Clone(tsparams.CertsuiteTargetPodLabels)
		extraPodLabels[tsparams.TestPodExtraLabelKey] = tsparams.TestPodExtraLabelValue
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, extraPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has one replica")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(1)))

		By("Create pod disruption budget with an extra label in matchExpressions")
		expressions := []metav1.LabelSelectorRequirement{
			{
				Key:      tsparams.TestPodExtraLabelKey,
				Operator: "In",
				Values:   []string{tsparams.TestPodExtraLabelValue},
			},
		}
		pdb := poddisruptionbudget.DefinePDBMinAvailableWithMatchLabelsAndExpressions(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.CertsuiteTargetPodLabels, expressions)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("One deployment with two labels, one matches the label in the PDB's matchLabels but no label matches the matchExpressions [negative]", func() {
		By("Define deployment with an extra label")
		extraPodLabels := maps.Clone(tsparams.CertsuiteTargetPodLabels)
		extraPodLabels[tsparams.TestPodExtraLabelKey] = tsparams.TestPodExtraLabelValue
		dep := deployment.DefineDeployment(tsparams.TestDeploymentBaseName, randomNamespace,
			globalhelper.GetConfiguration().General.TestImage, extraPodLabels)

		deployment.RedefineWithReplicaNumber(dep, 1)

		By("Create deployment")
		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.DeploymentDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Assert deployment has one replica")
		runningDeployment, err := globalhelper.GetRunningDeployment(dep.Namespace, dep.Name)
		Expect(err).ToNot(HaveOccurred())
		Expect(*runningDeployment.Spec.Replicas).To(Equal(int32(1)))

		By("Create pod disruption budget with an extra label in matchExpressions")
		expressions := []metav1.LabelSelectorRequirement{
			{
				Key:      tsparams.UnknownKey,
				Operator: "In",
				Values:   []string{tsparams.UnknownValue},
			},
		}
		pdb := poddisruptionbudget.DefinePDBMinAvailableWithMatchLabelsAndExpressions(tsparams.TestPdbBaseName, randomNamespace,
			intstr.FromInt(1), tsparams.CertsuiteTargetPodLabels, expressions)

		err = globalhelper.CreatePodDisruptionBudget(pdb, tsparams.PdbDeployTimeoutMins)
		Expect(err).ToNot(HaveOccurred())

		By("Start Certsuite " + tsparams.CertsuitePodDisruptionBudgetTcName + " test case")
		err = globalhelper.LaunchTests(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuitePodDisruptionBudgetTcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})
