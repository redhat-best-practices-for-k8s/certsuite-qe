package globalhelper

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// QueryPackageManifest queries the package manifest for the operator.
//
//nolint:gocognit
func QueryPackageManifestForVersion(operatorName, operatorNamespace, channel string) (string, error) {
	gvr := schema.GroupVersionResource{
		Group:    "packages.operators.coreos.com",
		Version:  "v1",
		Resource: "packagemanifests",
	}

	// Query the package manifest for the operator
	pkgManifest, err := GetAPIClient().DynamicClient.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return "", err
	}

	for _, item := range pkgManifest.Items {
		if item.GetName() == operatorName {
			// check if defaultChannel exists
			if _, ok := item.Object["status"].(map[string]interface{})["defaultChannel"]; !ok {
				continue
			}

			// type assertion
			if _, ok := item.Object["status"].(map[string]interface{})["defaultChannel"].(string); !ok {
				continue
			}

			//nolint:forcetypeassert
			fmt.Printf("Comparing %s with %s\n", item.Object["status"].(map[string]interface{})["defaultChannel"].(string), channel)

			//nolint:forcetypeassert
			if item.Object["status"].(map[string]interface{})["defaultChannel"].(string) != channel {
				continue
			}

			channelsObj, _, err := unstructured.NestedSlice(item.Object, "status", "channels")

			if err != nil {
				return "", err
			}

			for _, chanObj := range channelsObj {
				// verify the name of the channel is the same as the one we are looking for
				if _, ok := chanObj.(map[string]interface{})["name"]; !ok {
					continue
				}

				// type assertion
				if _, ok := chanObj.(map[string]interface{})["name"].(string); !ok {
					continue
				}

				//nolint:forcetypeassert
				fmt.Printf("Comparing name %s with channel %s\n", chanObj.(map[string]interface{})["name"].(string), channel)

				// skip this channel because it does not match the one we are looking for
				//nolint:forcetypeassert
				if chanObj.(map[string]interface{})["name"].(string) != channel {
					continue
				}

				// check if version exists
				if _, ok := chanObj.(map[string]interface{})["currentCSVDesc"].(map[string]interface{})["version"]; !ok {
					continue
				}

				// type assertion
				if _, ok := chanObj.(map[string]interface{})["currentCSVDesc"].(map[string]interface{})["version"].(string); !ok {
					continue
				}

				//nolint:forcetypeassert
				return chanObj.(map[string]interface{})["currentCSVDesc"].(map[string]interface{})["version"].(string), nil
			}
		}
	}

	return "not found", nil
}

func QueryPackageManifestForDefaultChannel(operatorName, operatorNamespace string) (string, error) {
	gvr := schema.GroupVersionResource{
		Group:    "packages.operators.coreos.com",
		Version:  "v1",
		Resource: "packagemanifests",
	}

	// Query the package manifest for the operator
	pkgManifest, err := GetAPIClient().DynamicClient.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return "", err
	}

	// Example of how to get the default channel
	// oc get packagemanifest cluster-logging -o json | jq .status.defaultChannel

	for _, item := range pkgManifest.Items {
		if item.GetName() == operatorName {
			// check if defaultChannel exists
			if _, ok := item.Object["status"].(map[string]interface{})["defaultChannel"]; !ok {
				continue
			}

			// type assertion
			if _, ok := item.Object["status"].(map[string]interface{})["defaultChannel"].(string); !ok {
				continue
			}

			//nolint:forcetypeassert
			return item.Object["status"].(map[string]interface{})["defaultChannel"].(string), nil
		}
	}

	return "not found", nil
}
