package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/platformalteration/parameters"
)

const (
	// rhcosVersionMapRelativePath is the path to the rhcos_version_map file relative to the certsuite repo root.
	rhcosVersionMapRelativePath = "tests/platform/operatingsystem/files/rhcos_version_map"
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

	It("OCP version should be supported", func() {
		if globalhelper.IsKindCluster() {
			Skip("OCP version is not applicable for Kind cluster")
		}

		By("Verify RHCOS versions exist in certsuite rhcos_version_map")
		err := verifyRHCOSVersionsInCertsuite()
		Expect(err).ToNot(HaveOccurred(), "RHCOS version validation failed")

		By("Start platform-alteration-ocp-lifecycle test")
		err = globalhelper.LaunchTests(tsparams.CertsuiteOCPLifecycleName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()), randomReportDir, randomCertsuiteConfigDir)
		Expect(err).ToNot(HaveOccurred())

		err = globalhelper.ValidateIfReportsAreValid(tsparams.CertsuiteOCPLifecycleName, globalparameters.TestCasePassed, randomReportDir)
		Expect(err).ToNot(HaveOccurred())
	})
})

// verifyRHCOSVersionsInCertsuite verifies that the RHCOS version from each node's OSImage
// exists in the certsuite rhcos_version_map file.
func verifyRHCOSVersionsInCertsuite() error {
	// Get all nodes
	nodesList, err := globalhelper.GetAPIClient().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	// Get certsuite repo path from config
	certsuiteRepoPath := globalhelper.GetConfiguration().General.CertsuiteRepoPath
	if certsuiteRepoPath == "" {
		return fmt.Errorf("CERTSUITE_REPO_PATH is not configured")
	}

	// Path to rhcos_version_map file - split the relative path and join with certsuite repo path
	pathParts := strings.Split(rhcosVersionMapRelativePath, "/")
	pathComponents := append([]string{certsuiteRepoPath}, pathParts...)
	rhcosVersionMapPath := filepath.Join(pathComponents...)

	// Read the rhcos_version_map file
	rhcosVersionMapData, err := os.ReadFile(rhcosVersionMapPath)
	if err != nil {
		return fmt.Errorf("failed to read rhcos_version_map file at %s: %w", rhcosVersionMapPath, err)
	}

	return verifyRHCOSVersions(nodesList.Items, string(rhcosVersionMapData))
}

// verifyRHCOSVersions is a parameterized function that verifies RHCOS versions from nodes
// against the rhcos_version_map content. This function is designed to be testable.
func verifyRHCOSVersions(nodes []corev1.Node, rhcosVersionMapContent string) error {
	const rhcosName = "Red Hat Enterprise Linux CoreOS"

	if len(nodes) == 0 {
		return fmt.Errorf("no nodes provided for verification")
	}

	if rhcosVersionMapContent == "" {
		return fmt.Errorf("rhcos_version_map content is empty")
	}

	// Check each node's RHCOS version
	for _, node := range nodes {
		osImage := node.Status.NodeInfo.OSImage

		// Extract RHCOS version from OSImage
		// e.g., "Red Hat Enterprise Linux CoreOS 410.84.202205031645-0 (Ootpa)" -> "410.84.202205031645-0"
		if !strings.Contains(osImage, rhcosName) {
			return fmt.Errorf("node %s has unexpected OS image format: %s (does not contain %s)",
				node.Name, osImage, rhcosName)
		}

		splitStr := strings.Split(osImage, rhcosName)
		if len(splitStr) < 2 {
			return fmt.Errorf("node %s has unexpected OS image format: %s (cannot split by %s)",
				node.Name, osImage, rhcosName)
		}

		longVersionSplit := strings.Split(strings.TrimSpace(splitStr[1]), " ")
		if len(longVersionSplit) == 0 || longVersionSplit[0] == "" {
			return fmt.Errorf("node %s has unexpected OS image format: %s (cannot extract version)",
				node.Name, osImage)
		}

		rhcosVersion := longVersionSplit[0]

		// Check if the version exists in the rhcos_version_map file
		if !strings.Contains(rhcosVersionMapContent, rhcosVersion) {
			return fmt.Errorf("RHCOS version %s from node %s not found in rhcos_version_map",
				rhcosVersion, node.Name)
		}

		fmt.Printf("✓ Node %s RHCOS version %s found in rhcos_version_map\n", node.Name, rhcosVersion)
	}

	return nil
}
