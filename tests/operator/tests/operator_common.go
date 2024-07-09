package operator

import (
	"log"

	. "github.com/onsi/gomega"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/operator/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/operator/parameters"
)

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

		if err != nil {
			log.Printf("Error getting csv: %s", err)

			return false
		}

		return false
	}, tsparams.Timeout, tsparams.PollingInterval).Should(Equal(true),
		csvPrefix+" is not ready.")

	return err
}
