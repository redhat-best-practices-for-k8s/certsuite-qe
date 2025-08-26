package accesscontrol

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/resourcequota"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Access-control namespace-resource-quota,", func() {
	var randomNamespace string
	var randomNamespace2 string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.TestAccessControlNameSpace)

		By("Create additional namespace for deployment2")
		randomNamespace2 = tsparams.AdditionalNamespaceForResourceQuotas + "-" + globalhelper.GenerateRandomString(5)
		err := globalhelper.CreateNamespace(randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		By("Define certsuite config file")
		err = globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)

		By("Delete additional namespace for deployment2")
		err := globalhelper.DeleteNamespaceAndWait(randomNamespace2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56469
	It("one deployment, one pod in a namespace with resource quota", func() {
		By("Define deployment")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateDeploymentNoWait(dep)
		Expect(err).ToNot(HaveOccurred())

		By("Define resource quota")
		resourceQuota := resourcequota.DefineResourceQuota("quota1", randomNamespace, tsparams.CPURequest,
			tsparams.MemoryRequest, tsparams.CPULimit, tsparams.MemoryLimit)

		By("Create resource quota")
		err = globalhelper.CreateResourceQuota(resourceQuota)
		Expect(err).ToNot(HaveOccurred())

		By("Assert ResourceQuota exists with expected hard limits")
		resourceQuotaObj, err := globalhelper.GetAPIClient().ResourceQuotas(randomNamespace).Get(
			context.TODO(), "quota1", metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		cpuReq := resourceQuotaObj.Spec.Hard[corev1.ResourceRequestsCPU]
		memReq := resourceQuotaObj.Spec.Hard[corev1.ResourceRequestsMemory]
		cpuLim := resourceQuotaObj.Spec.Hard[corev1.ResourceLimitsCPU]
		memLim := resourceQuotaObj.Spec.Hard[corev1.ResourceLimitsMemory]
		Expect((&cpuReq).String()).To(Equal(tsparams.CPURequest))
		Expect((&memReq).String()).To(Equal(tsparams.MemoryRequest))
		Expect((&cpuLim).String()).To(Equal(tsparams.CPULimit))
		Expect((&memLim).String()).To(Equal(tsparams.MemoryLimit))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56470
	It("one deployment, one pod in a namespace without resource quota [negative]", func() {
		By("Define deployment")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment")
		err = globalhelper.CreateDeploymentNoWait(dep)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56471
	It("two deployments, one pod each, both in a namespace with resource quota", func() {
		By("Define deployment 1")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 1")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Define resource quota")
		resourceQuota := resourcequota.DefineResourceQuota("quota1", randomNamespace, tsparams.CPURequest,
			tsparams.MemoryRequest, tsparams.CPULimit, tsparams.MemoryLimit)

		By("Create resource quota")
		err = globalhelper.CreateResourceQuota(resourceQuota)
		Expect(err).ToNot(HaveOccurred())

		By("Assert ResourceQuota exists with expected hard limits in namespace 1")
		rq1, err := globalhelper.GetAPIClient().ResourceQuotas(randomNamespace).Get(
			context.TODO(), "quota1", metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		q1cpuReq := rq1.Spec.Hard[corev1.ResourceRequestsCPU]
		q1memReq := rq1.Spec.Hard[corev1.ResourceRequestsMemory]
		q1cpuLim := rq1.Spec.Hard[corev1.ResourceLimitsCPU]
		q1memLim := rq1.Spec.Hard[corev1.ResourceLimitsMemory]
		Expect((&q1cpuReq).String()).To(Equal(tsparams.CPURequest))
		Expect((&q1memReq).String()).To(Equal(tsparams.MemoryRequest))
		Expect((&q1cpuLim).String()).To(Equal(tsparams.CPULimit))
		Expect((&q1memLim).String()).To(Equal(tsparams.MemoryLimit))

		By("Define deployment 2")
		dep2, err := tshelper.DefineDeploymentWithNamespace(1, 1, "accesscontroldeployment2",
			randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep2, tsparams.Timeout)
		Expect(err).ToNot(HaveOccurred())

		By("Define resource quota")
		resourceQuota = resourcequota.DefineResourceQuota("quota1", randomNamespace2, tsparams.CPURequest,
			tsparams.MemoryRequest, tsparams.CPULimit, tsparams.MemoryLimit)

		By("Create resource quota")
		err = globalhelper.CreateResourceQuota(resourceQuota)
		Expect(err).ToNot(HaveOccurred())

		By("Assert ResourceQuota exists with expected hard limits in namespace 2")
		rq2, err := globalhelper.GetAPIClient().ResourceQuotas(randomNamespace2).Get(
			context.TODO(), "quota1", metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		q2cpuReq := rq2.Spec.Hard[corev1.ResourceRequestsCPU]
		q2memReq := rq2.Spec.Hard[corev1.ResourceRequestsMemory]
		q2cpuLim := rq2.Spec.Hard[corev1.ResourceLimitsCPU]
		q2memLim := rq2.Spec.Hard[corev1.ResourceLimitsMemory]
		Expect((&q2cpuReq).String()).To(Equal(tsparams.CPURequest))
		Expect((&q2memReq).String()).To(Equal(tsparams.MemoryRequest))
		Expect((&q2cpuLim).String()).To(Equal(tsparams.CPULimit))
		Expect((&q2memLim).String()).To(Equal(tsparams.MemoryLimit))

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 56472
	It("two deployments, one pod each, one in a namespace without resource quota [negative]", func() {
		By("Define deployment 1")
		dep, err := tshelper.DefineDeployment(1, 1, "accesscontroldeployment1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 1")
		err = globalhelper.CreateDeploymentNoWait(dep)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment 2")
		dep2, err := tshelper.DefineDeploymentWithNamespace(1, 1, "accesscontroldeployment2",
			randomNamespace2)
		Expect(err).ToNot(HaveOccurred())

		By("Create deployment 2")
		err = globalhelper.CreateDeploymentNoWait(dep2)
		Expect(err).ToNot(HaveOccurred())

		By("Define resource quota")
		resourceQuota := resourcequota.DefineResourceQuota("quota1", randomNamespace2, tsparams.CPURequest,
			tsparams.MemoryRequest, tsparams.CPULimit, tsparams.MemoryLimit)

		By("Create resource quota")
		err = globalhelper.CreateResourceQuota(resourceQuota)
		Expect(err).ToNot(HaveOccurred())

		By("Start test")
		err = globalhelper.LaunchTests(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TestCaseNameAccessControlNamespaceResourceQuota,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})
