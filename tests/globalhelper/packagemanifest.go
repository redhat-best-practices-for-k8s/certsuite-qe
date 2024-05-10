package globalhelper

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// QueryPackageManifest queries the package manifest for the operator.
func QueryPackageManifestForVersion(operatorName, operatorNamespace string) (string, error) {
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
			channelsObj, _, err := unstructured.NestedSlice(item.Object, "status", "channels")

			if err != nil {
				return "", err
			}

			for _, channel := range channelsObj {
				// check if version exists
				if _, ok := channel.(map[string]interface{})["currentCSVDesc"].(map[string]interface{})["version"]; !ok {
					continue
				}

				// type assertion
				if _, ok := channel.(map[string]interface{})["currentCSVDesc"].(map[string]interface{})["version"].(string); !ok {
					continue
				}

				//nolint:forcetypeassert
				return channel.(map[string]interface{})["currentCSVDesc"].(map[string]interface{})["version"].(string), nil
			}
		}
	}

	return "not found", nil
}
