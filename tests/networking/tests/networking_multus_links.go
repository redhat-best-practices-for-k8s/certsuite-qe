package tests

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/parameters"
)

var _ = Describe("Networking custom namespace,", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(tsparams.TestNetworkingNameSpace)

		By("Define certsuite config file")
		err := globalhelper.DefineCertsuiteConfig(
			[]string{randomNamespace},
			[]string{tsparams.TestPodLabel},
			[]string{},
			[]string{},
			[]string{}, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.WaitingTime)
	})

	// expectNetworksAnnotation asserts that the multus networks annotation exists and
	// contains exactly the expected NAD names (order-insensitive).
	expectNetworksAnnotation := func(annotations map[string]string, expectedNadNames []string) {
		Expect(annotations).To(HaveKey("k8s.v1.cni.cncf.io/networks"))
		raw := annotations["k8s.v1.cni.cncf.io/networks"]
		var nets []struct {
			Name string `json:"name"`
		}
		Expect(json.Unmarshal([]byte(raw), &nets)).To(Succeed())
		got := make([]string, 0, len(nets))
		for _, n := range nets {
			got = append(got, n.Name)
		}
		// Order-insensitive comparison without varargs conversion complexity
		Expect(got).To(HaveLen(len(expectedNadNames)))
		for _, expected := range expectedNadNames {
			Expect(got).To(ContainElement(expected))
		}
	}

	// 48328
	It("custom deployment 3 pods, 1 NAD, connectivity via Multus secondary interface", func() {
		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		// Assert deployment has multus annotation
		runningDep48328, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(runningDep48328.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48330
	It("2 custom deployments 3 pods, 1 NAD, connectivity via Multus secondary interface", func() {
		// The NetworkAttachmentDefinition (mcvlan) created for this TC uses the default interface that is connecting
		// all worker/master nodes so that pods have connectivity irrespective of the node they are scheduled on
		// see https://github.com/redhat-best-practices-for-k8s/certsuite-qe/pull/263

		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define first deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		depA48330, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depA48330.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		depB48330, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentBName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depB48330.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48331
	It("custom deployment and daemonset 3 pods, 2 NADs, connectivity via Multus secondary interfaces", func() {
		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, randomNamespace, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define first deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		depA48331, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depA48331.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Define second deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, randomNamespace, []string{tsparams.TestNadNameB}, 3)
		Expect(err).ToNot(HaveOccurred())

		depB48331, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentBName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depB48331.Spec.Template.Annotations, []string{tsparams.TestNadNameB})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48334
	It("custom deployment 3 pods, 1 NAD missing IP, connectivity via Multus secondary interface[skip]", func() {
		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameA, randomNamespace, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		dep48334, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(dep48334.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48338
	It("custom deployments 3 pods and 1 pod, standalone IP, connectivity via Multus secondary interface[skip]", func() {

		By("Define and create Network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameB, randomNamespace, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment-a and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 1)
		Expect(err).ToNot(HaveOccurred())

		depA48338, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depA48338.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Define deployment-b and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentBName, randomNamespace, []string{tsparams.TestNadNameB}, 3)
		Expect(err).ToNot(HaveOccurred())

		depB48338, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentBName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depB48338.Spec.Template.Annotations, []string{tsparams.TestNadNameB})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48343
	It("custom deployment and daemonset 3 pods, daemonset missing ip, 2 NADs, connectivity via Multus "+
		"secondary interface", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(tsparams.TestNadNameB, randomNamespace, "")
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		dep48343, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(dep48343.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Define daemonset and create it on cluster")

		err = tshelper.DefineAndCreateDaemonsetWithMultusOnCluster(tsparams.TestNadNameB, randomNamespace, "ds1")
		Expect(err).ToNot(HaveOccurred())

		// time.Sleep(10 * time.Minute)

		// Assert daemonset is created and has multus annotation
		daemonset, err := globalhelper.GetRunningDaemonsetByName("ds1", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(daemonset.Spec.Template.Annotations, []string{tsparams.TestNadNameB})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48580
	It("custom daemonset 3 pods with skip label [skip]", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA, randomNamespace, "ds2")
		Expect(err).ToNot(HaveOccurred())

		ds2, err := globalhelper.GetRunningDaemonsetByName("ds2", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(ds2.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48582
	It("custom deployment and daemonset 3 pods with skip label[skip]", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA, randomNamespace, "ds3")
		Expect(err).ToNot(HaveOccurred())

		ds3, err := globalhelper.GetRunningDaemonsetByName("ds3", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(ds3.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		dep48382skip, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(dep48382skip.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCaseSkipped, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48582
	FIt("custom deployment and daemonSet 3 pods, daemonSet has skip label", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusAndSkipLabelOnCluster(tsparams.TestNadNameA, randomNamespace, "ds4")
		Expect(err).ToNot(HaveOccurred())

		ds4, err := globalhelper.GetRunningDaemonsetByName("ds4", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(ds4.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		dep48382pass, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(dep48382pass.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		// Verify all of the multus interfaces are setup
		daemonset, err := globalhelper.GetRunningDaemonsetByName("ds4", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		for _, container := range daemonset.Spec.Template.Spec.Containers {
			for _, volume := range container.VolumeMounts {
				Expect(volume.Name).To(Equal("multus-cni-network"))
			}
		}
		deployment, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		for _, container := range deployment.Spec.Template.Spec.Containers {
			for _, volume := range container.VolumeMounts {
				Expect(volume.Name).To(Equal("multus-cni-network"))
			}
		}

		// By("Sleep for 10 minutes")

		// The certsuite requires the network-status annotation to be set in the deployment and daemonset before the test is run.

		time.Sleep(10 * time.Minute)

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48582
	It("custom deployment 3 pods, 2 NADs, multiple Multus interfaces on deployment", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, randomNamespace, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA, tsparams.TestNadNameB}, 3)
		Expect(err).ToNot(HaveOccurred())

		depMulti48582, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depMulti48582.Spec.Template.Annotations, []string{tsparams.TestNadNameA, tsparams.TestNadNameB})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())

	})

	// 48346
	It("custom deployment 3 pods,1 NAD,no connectivity via Multus secondary interface[negative]", func() {

		By("Define and create Network-attachment-definition")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		depNeg48346, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depNeg48346.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Put one deployment's pod interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48347
	It("custom deployment and daemonset 3 pods, 2 NADs, No connectivity on daemonset via Multus secondary "+
		"interface[negative]", func() {

		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, randomNamespace, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		depNeg48347, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depNeg48347.Spec.Template.Annotations, []string{tsparams.TestNadNameA})

		By("Put one deployment's pod interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusOnCluster(tsparams.TestNadNameB, randomNamespace, "ds5")
		Expect(err).ToNot(HaveOccurred())

		ds5, err := globalhelper.GetRunningDaemonsetByName("ds5", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(ds5.Spec.Template.Annotations, []string{tsparams.TestNadNameB})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

	// 48590
	It("custom deployment and daemonset 3 pods, 2 NADs, multiple Multus interfaces on deployment no "+
		"connectivity via secondary interface[negative]", func() {

		By("Define and create network-attachment-definitions")
		err := tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameA, randomNamespace, tsparams.TestIPamIPNetworkA)
		Expect(err).ToNot(HaveOccurred())

		err = tshelper.DefineAndCreateNadOnCluster(
			tsparams.TestNadNameB, randomNamespace, tsparams.TestIPamIPNetworkB)
		Expect(err).ToNot(HaveOccurred())

		By("Define deployment and create it on cluster")
		err = tshelper.DefineAndCreateDeploymentWithMultusOnCluster(
			tsparams.TestDeploymentAName, randomNamespace, []string{tsparams.TestNadNameB, tsparams.TestNadNameA}, 3)
		Expect(err).ToNot(HaveOccurred())

		depNeg48590, err := globalhelper.GetRunningDeployment(randomNamespace, tsparams.TestDeploymentAName)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(depNeg48590.Spec.Template.Annotations, []string{tsparams.TestNadNameA, tsparams.TestNadNameB})

		By("Put one deployment's pod interface down")
		err = tshelper.ExecCmdOnOnePodInNamespace([]string{"ip", "link", "set", "net1", "down"}, randomNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Define daemonset and create it on cluster")
		err = tshelper.DefineAndCreateDaemonsetWithMultusOnCluster(tsparams.TestNadNameB, randomNamespace, "ds6")
		Expect(err).ToNot(HaveOccurred())

		ds6, err := globalhelper.GetRunningDaemonsetByName("ds6", randomNamespace)
		Expect(err).ToNot(HaveOccurred())
		expectNetworksAnnotation(ds6.Spec.Template.Annotations, []string{tsparams.TestNadNameB})

		By("Start tests")
		err = globalhelper.LaunchTests(
			tsparams.CertsuiteMultusIpv4TcName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Verify test case status in Claim report")
		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.CertsuiteMultusIpv4TcName,
			globalparameters.TestCaseFailed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})

})
