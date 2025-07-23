package tests

import (
	"fmt"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/affiliatedcertification/parameters"
	utils "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/operator"
)

func preConfigureAffiliatedCertificationEnvironment(namespace, configDir string) {
	By("Clean test namespace")

	err := globalhelper.CleanNamespace(namespace)
	Expect(err).ToNot(HaveOccurred(),
		"Error cleaning namespace "+namespace)
	By("Ensure default catalog source is enabled")

	catalogEnabled, err := globalhelper.IsCatalogSourceEnabled(
		tsparams.CertifiedOperatorGroup,
		tsparams.OperatorSourceNamespace,
		tsparams.CertifiedOperatorDisplayName)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

	if !catalogEnabled {
		Expect(
			globalhelper.EnableCatalogSource(tsparams.CertifiedOperatorGroup)).ToNot(
			HaveOccurred())
		Eventually(func() bool {
			catalogEnabled, err = globalhelper.IsCatalogSourceEnabled(
				tsparams.CertifiedOperatorGroup,
				tsparams.OperatorSourceNamespace,
				tsparams.CertifiedOperatorDisplayName)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

			return catalogEnabled
		}, tsparams.TimeoutLabelCsv, tsparams.PollingInterval).Should(BeTrue(),
			fmt.Sprintf("Default catalog source %s is not enabled",
				tsparams.CertifiedOperatorGroup))
	}

	By("Deploy OperatorGroup if not already deployed")

	if globalhelper.IsOperatorGroupInstalled(tsparams.OperatorGroupName,
		namespace) != nil {
		err = globalhelper.DeployOperatorGroup(namespace,
			utils.DefineOperatorGroup(tsparams.OperatorGroupName,
				namespace,
				[]string{namespace}),
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operatorgroup")
	}

	By("Define config file " + globalparameters.DefaultCertsuiteConfigFileName)

	err = globalhelper.DefineCertsuiteConfig(
		[]string{namespace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{},
		[]string{}, configDir)
	Expect(err).ToNot(HaveOccurred(), "Error defining certsuite config file")
}

func approveInstallPlanWhenReady(csvName, namespace string) {
	Eventually(func() bool {
		installPlan, err := globalhelper.GetInstallPlanByCSV(namespace, csvName)
		if err != nil {
			return false
		}

		if installPlan.Spec.Approval == v1alpha1.ApprovalAutomatic {
			return true
		}

		if installPlan.Status.Phase == v1alpha1.InstallPlanPhaseRequiresApproval {
			err = globalhelper.ApproveInstallPlan(namespace,
				installPlan)

			return err == nil
		}

		return false
	}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
		csvName+" install plan is not ready.")
}
