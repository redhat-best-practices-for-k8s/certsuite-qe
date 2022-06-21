package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func DeleteNamespaces(nsToDelete []string, clientSet *client.ClientSet, timeout time.Duration) error {
	var failedNs []string

	for _, namespace := range nsToDelete {
		err := namespaces.DeleteAndWait(
			clientSet,
			namespace,
			timeout,
		)
		if err != nil {
			failedNs = append(failedNs, namespace)
		}
	}

	if len(failedNs) > 0 {
		return fmt.Errorf("Failed to remove the following namespaces: " +
			strings.Join(failedNs, ", "))
	}

	return nil
}
