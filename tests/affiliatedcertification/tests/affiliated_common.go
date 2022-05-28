package tests

import (
	"fmt"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcerthelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/affiliatedcertification/affiliatedcertparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	utils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/operator"
)

func preConfigureAffiliatedCertificationEnvironment() {
	By("Clean test namespace")

	err := namespaces.Clean(affiliatedcertparameters.TestCertificationNameSpace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred(),
		"Error cleaning namespace "+affiliatedcertparameters.TestCertificationNameSpace)

	By("Ensure default catalog source is enabled")

	catalogEnabled, err := affiliatedcerthelper.IsCatalogSourceEnabled(
		affiliatedcertparameters.CertifiedOperatorGroup,
		affiliatedcertparameters.OperatorSourceNamespace,
		affiliatedcertparameters.CertifiedOperatorDisplayName)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

	if !catalogEnabled {
		Expect(
			affiliatedcerthelper.EnableCatalogSource(affiliatedcertparameters.CertifiedOperatorGroup)).ToNot(
			HaveOccurred())
		Eventually(func() bool {
			catalogEnabled, err = affiliatedcerthelper.IsCatalogSourceEnabled(
				affiliatedcertparameters.CertifiedOperatorGroup,
				affiliatedcertparameters.OperatorSourceNamespace,
				affiliatedcertparameters.CertifiedOperatorDisplayName)
			Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("can not collect catalogSource object due to %s", err))

			return catalogEnabled
		}, affiliatedcertparameters.TimeoutLabelCsv, affiliatedcertparameters.PollingInterval).Should(BeTrue(),
			fmt.Sprintf("Default catalog source %s is not enabled",
				affiliatedcertparameters.CertifiedOperatorGroup))
	}

	By("Deploy OperatorGroup if not already deployed")

	if affiliatedcerthelper.IsOperatorGroupInstalled(affiliatedcertparameters.OperatorGroupName,
		affiliatedcertparameters.TestCertificationNameSpace) != nil {
		err = affiliatedcerthelper.DeployOperatorGroup(affiliatedcertparameters.TestCertificationNameSpace,
			utils.DefineOperatorGroup(affiliatedcertparameters.OperatorGroupName,
				affiliatedcertparameters.TestCertificationNameSpace,
				[]string{affiliatedcertparameters.TestCertificationNameSpace}),
		)
		Expect(err).ToNot(HaveOccurred(), "Error deploying operatorgroup")
	}

	By("Define config file " + globalparameters.DefaultTnfConfigFileName)

	err = globalhelper.DefineTnfConfig(
		[]string{affiliatedcertparameters.TestCertificationNameSpace},
		[]string{affiliatedcertparameters.TestPodLabel},
		[]string{})
	Expect(err).ToNot(HaveOccurred(), "Error defining tnf config file")
}

func waitUntilOperatorIsReady(csvPrefix, namespace string) error {
	var err error

	var csv *v1alpha1.ClusterServiceVersion

	Eventually(func() bool {
		csv, err = affiliatedcerthelper.GetCsvByPrefix(csvPrefix, namespace)
		if csv != nil && csv.Status.Phase != v1alpha1.CSVPhaseNone {
			return csv.Status.Phase != "InstallReady" &&
				csv.Status.Phase != "Deleting" &&
				csv.Status.Phase != "Replacing" &&
				csv.Status.Phase != "Unknown"
		}

		return false
	}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
		csvPrefix+" is not ready.")

	return err
}

func approveInstallPlanWhenReady(csvName, namespace string) {
	Eventually(func() bool {
		installPlan, err := affiliatedcerthelper.GetInstallPlanByCSV(namespace, csvName)
		if err != nil {
			return false
		}

		if installPlan.Spec.Approval == v1alpha1.ApprovalAutomatic {
			return true
		}

		if installPlan.Status.Phase == v1alpha1.InstallPlanPhaseRequiresApproval {
			err = affiliatedcerthelper.ApproveInstallPlan(affiliatedcertparameters.TestCertificationNameSpace,
				installPlan)

			return err == nil
		}

		return false
	}, affiliatedcertparameters.Timeout, affiliatedcertparameters.PollingInterval).Should(Equal(true),
		csvName+" install plan is not ready.")
}
