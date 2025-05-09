package tests

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"golang.org/x/mod/semver"
)

const (
	// Change this to the minimum OCP version that went EOL.
	// 
	minimumOCPVersionEOL = "v4.14"
)

var _ = Describe("platform-alteration-ocp-lifecycle", func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	BeforeEach(func() {
		// Create random namespace and keep original report and certsuite config directories
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomNamespace(
				tsparams.PlatformAlterationNamespace)

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

	FIt("OCP version should be supported", func() {

		// Function to normalize the version string to a standard format
		// e.g., "4.14" -> "v4.14.0", "4.14.0" -> "v4.14.0", "4.14.0-ec.0" -> "v4.14.0"
		normalizeSemver := func(version string) string {
			if !strings.HasPrefix(version, "v") {
				version = "v" + version
			}
			parts := strings.SplitN(version[1:], ".", 3)
			switch len(parts) {
			case 1:
				return "v" + parts[0] + ".0.0"
			case 2:
				return "v" + parts[0] + "." + parts[1] + ".0"
			default:
				return version
			}
		}

		if globalhelper.IsKindCluster() {
			Skip("OCP version is not applicable for Kind cluster")
		}

		By("Start platform-alteration-ocp-lifecycle test")
		err := globalhelper.LaunchTests(tsparams.CertsuiteOCPLifecycleName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		By("Check OCP version and validate expected result")
		ocpVersion, err := globalhelper.GetClusterVersion()
		Expect(err).ToNot(HaveOccurred())

		// OCP version is in the format 4.14.0-ec.0, we need to trim it to 4.14
		ocpVersion = normalizeSemver(ocpVersion[:4])
		minRequiredVersion := normalizeSemver(minimumOCPVersionEOL)

		// Check if the OCP version is less than the minimum required version
		if semver.Compare(ocpVersion, minRequiredVersion) > 0 {
			err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteOCPLifecycleName, globalparameters.TestCasePassed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		} else {
			err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteOCPLifecycleName, globalparameters.TestCaseFailed, randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		}
	})
})
