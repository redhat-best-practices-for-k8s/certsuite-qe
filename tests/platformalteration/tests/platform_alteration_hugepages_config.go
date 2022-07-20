package tests

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/crd"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

var _ = Describe("platform-alteration-hugepages-config", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

	})

	// 51308
	It("Hugepages config unchanged configuration", func() {

		crdExists, err := crd.EnsureCrdExists(tsparams.PerformanceProfileCrd)

		if !crdExists || err != nil {
			Skip("performance profile does not exist.")
		}

		// cluster should be set with kernel hugepages = MC hugepages configuration by performance profile.
		By("Start platform-alteration-hugepages-config test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePagesConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(
			tsparams.TnfHugePagesConfigName,
			globalparameters.TestCasePassed)
		Expect(err).ToNot(HaveOccurred())
	})

	// 51309
	It("Change Hugepages config manually [negative]", func() {

		By("Create daemonSet")
		daemonset := daemonset.RedefineWithPriviledgedContainer(
			daemonset.RedefineWithVolumeMount(
				daemonset.DefineDaemonSet(tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
					tsparams.TnfTargetPodLabels, tsparams.TestDaemonSetName)))

		err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonset, tsparams.WaitingTime)
		Expect(err).ToNot(HaveOccurred())

		podList, err := globalhelper.GetListOfPodsInNamespace(tsparams.PlatformAlterationNamespace)
		Expect(err).ToNot(HaveOccurred())

		By("Get first hugepages file")
		nrHugepagesFiles, err := globalhelper.ExecCommand(
			podList.Items[0], []string{"/bin/bash", "-c", tsparams.FindHugePagesFiles})
		Expect(err).ToNot(HaveOccurred())

		hugePagesPaths := strings.Split(nrHugepagesFiles.String(), "\r\n")

		By("Get hugepages config")
		cmd := fmt.Sprintf("cat %s", hugePagesPaths[0])
		buf, err := globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", cmd})
		Expect(err).ToNot(HaveOccurred())

		hugepagesNumber, err := strconv.Atoi(strings.Split(buf.String(), "\r\n")[0])
		Expect(err).ToNot(HaveOccurred())

		By("Change hugepages config")
		cmd = fmt.Sprintf("echo %d > %s", hugepagesNumber+1, hugePagesPaths[0])
		_, err = globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", cmd})
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-hugepages-config test")
		err = globalhelper.LaunchTests(tsparams.TnfHugePagesConfigName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfHugePagesConfigName, globalparameters.TestCaseFailed)
		Expect(err).ToNot(HaveOccurred())

		By("Fix hugepages config")
		cmd = fmt.Sprintf("echo %d > %s", hugepagesNumber, hugePagesPaths[0])
		_, err = globalhelper.ExecCommand(podList.Items[0], []string{"/bin/bash", "-c", cmd})
		Expect(err).ToNot(HaveOccurred())
	})

})
