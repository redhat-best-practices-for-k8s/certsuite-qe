//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	semver "github.com/blang/semver/v4"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	packagemanifestv1 "github.com/openshift-kni/eco-goinfra/pkg/schemes/olm/package-server/operators/v1"
	packageversion "github.com/openshift-kni/eco-goinfra/pkg/schemes/olm/version"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var coverageOuterT *testing.T

func TestGlobalHelperGinkgoSuite(t *testing.T) {
	coverageOuterT = t

	gomega.RegisterFailHandler(ginkgo.Fail)
	coverageFixture(t)
	conf = &config.Config{}
	conf.General.CertsuiteReportDir = filepath.Join(t.TempDir(), "reports")
	conf.General.CertsuiteConfigDir = filepath.Join(t.TempDir(), "configs")
	conf.General.ReportDirAbsPath = t.TempDir()
	ginkgo.RunSpecs(t, "globalhelper coverage")
}

var _ = ginkgo.Describe("Ginkgo backed global helpers", func() {
	ginkgo.It("sets up and cleans a test namespace and report directories", func() {
		namespace, reportDir, configDir := BeforeEachSetupWithRandomNamespace("coverage")
		gomega.Expect(namespace).To(gomega.HavePrefix("coverage-"))
		gomega.Expect(reportDir).To(gomega.BeADirectory())
		gomega.Expect(configDir).To(gomega.BeADirectory())
		AfterEachCleanupWithRandomNamespace(namespace, reportDir, configDir, time.Second)

		originalRunSpecs := runSpecs
		runSpecs = func(ginkgo.GinkgoTestingT, string, ...any) bool { return true }

		RunSuite(coverageOuterT, "nested suite")
		runSpecs = originalRunSpecs

		LogCheckDetails(&CheckDetails{
			CompliantObjectsOut:    []*ReportObject{{ObjectType: "Pod", ObjectFieldsKeys: []string{ReasonForCompliance}, ObjectFieldsValues: []string{"good"}}},
			NonCompliantObjectsOut: []*ReportObject{{ObjectType: "Pod", ObjectFieldsKeys: []string{ReasonForNonCompliance}, ObjectFieldsValues: []string{"bad"}}},
		}, nil)
		LogCheckDetails(nil, coverageError{})
	})

	ginkgo.It("creates a privileged test namespace and cleans its resources", func() {
		namespace, reportDir, configDir := BeforeEachSetupWithRandomPrivilegedNamespace("coverage-privileged")
		ns, err := GetAPIClient().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(ns.Labels["pod-security.kubernetes.io/enforce"]).To(gomega.Equal("privileged"))
		AfterEachCleanupWithRandomNamespace(namespace, reportDir, configDir, time.Second)
	})

	ginkgo.It("checks for operators with the Ginkgo skip wrappers", func() {
		manifest := &packagemanifestv1.PackageManifest{
			ObjectMeta: metav1.ObjectMeta{Name: "coverage-operator", Namespace: "marketplace"},
			Status: packagemanifestv1.PackageManifestStatus{
				CatalogSource: "community-operators", DefaultChannel: "stable",
				Channels: []packagemanifestv1.PackageChannel{{
					Name: "stable", CurrentCSV: "coverage-operator.v1.0.0",
					CurrentCSVDesc: packagemanifestv1.CSVDescription{Version: packageversion.OperatorVersion{Version: semver.Version{Major: 1}}},
				}},
			},
		}
		setCoverageCtrlObjects(manifest)
		name, source := CheckOperatorExistsOrSkip("coverage-operator", "marketplace")
		gomega.Expect(name).To(gomega.Equal("coverage-operator"))
		gomega.Expect(source).To(gomega.Equal("community-operators"))
		channel, version, csv := CheckOperatorChannelAndVersionOrSkip("coverage-operator", "marketplace")
		gomega.Expect(channel).To(gomega.Equal("stable"))
		gomega.Expect(version).To(gomega.Equal("1.0.0"))
		gomega.Expect(csv).To(gomega.Equal("coverage-operator.v1.0.0"))
	})
})

type coverageError struct{}

func (coverageError) Error() string { return "ignored" }
