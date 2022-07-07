package helper

import (
	"fmt"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func DeleteNamespaces(nsToDelete []string, clientSet *client.ClientSet, timeout time.Duration) error {
	failedNs := make(map[string]error)

	for _, namespace := range nsToDelete {
		err := namespaces.DeleteAndWait(
			clientSet,
			namespace,
			timeout,
		)
		if err != nil {
			failedNs[namespace] = err
		}
	}

	if len(failedNs) > 0 {
		return fmt.Errorf("failed to remove the following namespaces: %v", failedNs)
	}

	return nil
}
