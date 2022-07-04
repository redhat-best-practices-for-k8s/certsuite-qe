package tests

// import (
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
// 	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
// 	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
// )

// var _ = Describe("platform-alteration-is-redhat-release", func() {

// 	BeforeEach(func() {
// 		By("Clean namespace before each test")
// 		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
// 		Expect(err).ToNot(HaveOccurred())
// 	})

// 	// 51319
// 	It("One deployment, one pod, several containers, all running Red Hat release", func() {

// 		By("Define deployment")
// 		deployment := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.PlatformAlterationNamespace,
// 			globalhelper.Configuration.General.TestImage, tsparams.TnfTargetPodLabels)

// 		globalhelper.AppendContainersToDeployment(deployment, 3, globalhelper.Configuration.General.TestImage)

// 		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
// 		Expect(err).ToNot(HaveOccurred())

// 		By("Start platform-alteration-is-redhat-release test")
// 		err = globalhelper.LaunchTests(
// 			tsparams.TnfIsRedHatReleaseName,
// 			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
// 		Expect(err).ToNot(HaveOccurred())

// 		By("Verify test case status in Junit and Claim reports")
// 		err = globalhelper.ValidateIfReportsAreValid(
// 			tsparams.TnfIsRedHatReleaseName,
// 			globalparameters.TestCasePassed)
// 		Expect(err).ToNot(HaveOccurred())
// 	})

// 	// 51320
// 	It("One daemonSet that is running Red Hat release", func() {

// 		By("Define daemonSet")
// 		daemonSet := daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace,
// 			globalhelper.Configuration.General.TestImage,
// 			tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)

// 		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
// 		Expect(err).ToNot(HaveOccurred())

// 		By("Start platform-alteration-is-redhat-release test")
// 		err = globalhelper.LaunchTests(
// 			tsparams.TnfIsRedHatReleaseName,
// 			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
// 		Expect(err).ToNot(HaveOccurred())

// 		By("Verify test case status in Junit and Claim reports")
// 		err = globalhelper.ValidateIfReportsAreValid(
// 			tsparams.TnfIsRedHatReleaseName,
// 			globalparameters.TestCasePassed)
// 		Expect(err).ToNot(HaveOccurred())
// 	})

// 	// 51321
// 	It("One deployment, one pod, 2 containers, one running Red Hat release, other is not [negative]", func() {

// 		By("Define deployment")
// 		deployment := deployment.DefineDeployment(tsparams.TestDeploymentName, tsparams.PlatformAlterationNamespace,
// 			tsparams.NotRedHatRelease, tsparams.TnfTargetPodLabels)

// 		globalhelper.AppendContainersToDeployment(deployment, 1, globalhelper.Configuration.General.TestImage)

// 		err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deployment, tsparams.WaitingTime)
// 		Expect(err).ToNot(HaveOccurred())

// 		By("Start platform-alteration-is-redhat-release test")
// 		err = globalhelper.LaunchTests(
// 			tsparams.TnfIsRedHatReleaseName,
// 			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
// 		Expect(err).To(HaveOccurred())

// 		By("Verify test case status in Junit and Claim reports")
// 		err = globalhelper.ValidateIfReportsAreValid(
// 			tsparams.TnfIsRedHatReleaseName,
// 			globalparameters.TestCaseFailed)
// 		Expect(err).ToNot(HaveOccurred())
// 	})

// 	// 51326
// 	It("One statefulSet, one pod that is not running Red Hat release [negative]", func() {

// 		By("Define statefulSet")
// 		statefulSet := statefulset.DefineStatefulSet(tsparams.TestStatefulSetName, tsparams.PlatformAlterationNamespace,
// 			tsparams.NotRedHatRelease, tsparams.TnfTargetPodLabels)

// 		err := globalhelper.CreateAndWaitUntilStatefulSetIsReady(statefulSet, tsparams.WaitingTime)
// 		Expect(err).ToNot(HaveOccurred())

// 		By("Start platform-alteration-is-redhat-release test")
// 		err = globalhelper.LaunchTests(
// 			tsparams.TnfIsRedHatReleaseName,
// 			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
// 		Expect(err).To(HaveOccurred())

// 		By("Verify test case status in Junit and Claim reports")
// 		err = globalhelper.ValidateIfReportsAreValid(
// 			tsparams.TnfIsRedHatReleaseName,
// 			globalparameters.TestCaseFailed)
// 		Expect(err).ToNot(HaveOccurred())
// 	})
// })
