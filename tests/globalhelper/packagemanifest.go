package globalhelper

import (
	"fmt"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiOLM "github.com/openshift-kni/eco-goinfra/pkg/olm"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// QueryPackageManifest queries the package manifest for the operator.
func QueryPackageManifestForVersion(operatorName, operatorNamespace, channel string) (string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", err
	}

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			fmt.Printf("Comparing %s with %s\n", item.Object.Status.DefaultChannel, channel)

			// skip if the default channel is not the one we are looking for
			if item.Object.Status.DefaultChannel != channel {
				continue
			}

			for _, chanObj := range item.Object.Status.Channels {
				fmt.Printf("Comparing name %s with channel %s\n", chanObj.Name, channel)

				// skip if the name of the channel is not the one we are looking for
				if chanObj.Name != channel {
					continue
				}

				return chanObj.CurrentCSVDesc.Version.String(), nil
			}
		}
	}

	return "not found", nil
}

func QueryPackageManifestForDefaultChannel(operatorName, operatorNamespace string) (string, error) {
	pkgManifest, err := egiOLM.ListPackageManifest(egiClients.New(""), operatorNamespace, client.ListOptions{})

	if err != nil {
		return "", err
	}

	// Example of how to get the default channel
	// oc get packagemanifest cluster-logging -o json | jq .status.defaultChannel

	for _, item := range pkgManifest {
		if item.Object.GetName() == operatorName {
			return item.Object.Status.DefaultChannel, nil
		}
	}

	return "not found", nil
}
