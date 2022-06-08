package tests

import (
	"fmt"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"

	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/parameters"
	utils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"
)

func preConfigureAffiliatedCertificationEnvironment() {
	By("Clean test namespace")

	err := namespaces.Clean(tsparams.TestCertificationNameSpace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred(),
		"Error cleaning namespace "+tsparams.TestCertificationNameSpace)

	By("Ensure default catalog source is enabled")

	catalogEnabled, err := tshelper.IsCatalogSourceEnabled(
		tsparams.CertifiedOperatorGroup,
		tsparams.OperatorSourceNamespace,
		tsparams.CertifiedOperatorDisplayName)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

	if !catalogEnabled {
		Expect(
			tshelper.EnableCatalogSource(tsparams.CertifiedOperatorGroup)).ToNot(
			HaveOccurred())
		Eventually(func() bool {
			catalogEnabled, err = tshelper.IsCatalogSourceEnabled(
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

	if tshelper.IsOperatorGroupInstalled(tsparams.OperatorGroupName,
		tsparams.TestCertificationNameSpace) != nil {
		err = tshelper.DeployOperatorGroup(tsparams.TestCertificationNameSpace,
			utils.DefineOperatorGroup(tsparams.OperatorGroupName,
				tsparams.TestCertificationNameSpace,
				[]string{tsparams.TestCertificationNameSpace}),
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operatorgroup")
	}

	By("Define config file " + globalparameters.DefaultTnfConfigFileName)

	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.TestCertificationNameSpace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")
}

func waitUntilOperatorIsReady(csvPrefix, namespace string) error {
	var err error

	var csv *v1alpha1.ClusterServiceVersion

	Eventually(func() bool {
		csv, err = tshelper.GetCsvByPrefix(csvPrefix, namespace)
		if csv != nil && csv.Status.Phase != v1alpha1.CSVPhaseNone {
			return csv.Status.Phase != "InstallReady" &&
				csv.Status.Phase != "Deleting" &&
				csv.Status.Phase != "Replacing" &&
				csv.Status.Phase != "Unknown"
		}

		return false
	}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
		csvPrefix+" is not ready.")

	return err
}

func approveInstallPlanWhenReady(csvName, namespace string) {
	Eventually(func() bool {
		installPlan, err := tshelper.GetInstallPlanByCSV(namespace, csvName)
		if err != nil {
			return false
		}

		if installPlan.Spec.Approval == v1alpha1.ApprovalAutomatic {
			return true
		}

		if installPlan.Status.Phase == v1alpha1.InstallPlanPhaseRequiresApproval {
			err = tshelper.ApproveInstallPlan(tsparams.TestCertificationNameSpace,
				installPlan)

			return err == nil
		}

		return false
	}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
		csvName+" install plan is not ready.")
}
