package crd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	crscaleoperator "github.com/test-network-function/cr-scale-operator/api/v1"
)

func DefineCustomResourceDefinition(names apiextv1.CustomResourceDefinitionNames,
	group string, addStatusSubresource bool) *apiextv1.CustomResourceDefinition {
	// Helper object for the fake "v1" version schema.
	version := apiextv1.CustomResourceDefinitionVersion{
		Served:  true,
		Storage: true,
		Name:    "v1",
		Schema: &apiextv1.CustomResourceValidation{
			OpenAPIV3Schema: &apiextv1.JSONSchemaProps{
				Type: "object",
				Properties: map[string]apiextv1.JSONSchemaProps{
					"spec": {
						Type: "object",
						Properties: map[string]apiextv1.JSONSchemaProps{
							"specProperty1": {
								Description: "Fake spec property",
								Type:        "string",
							},
						},
					},
				},
			},
		},
	}

	// Add the status schema property and the status subresource.
	if addStatusSubresource {
		version.Schema.OpenAPIV3Schema.Properties["status"] = apiextv1.JSONSchemaProps{
			Type: "object",
			Properties: map[string]apiextv1.JSONSchemaProps{
				"statusProperty1": {
					Description: "Fake status property",
					Type:        "string",
				},
			},
		}

		version.Subresources = &apiextv1.CustomResourceSubresources{
			Status: &apiextv1.CustomResourceSubresourceStatus{},
		}
	}

	return &apiextv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: names.Plural + "." + group,
		},
		Spec: apiextv1.CustomResourceDefinitionSpec{
			Group:    group,
			Names:    names,
			Scope:    "Namespaced",
			Versions: []apiextv1.CustomResourceDefinitionVersion{version},
		},
	}
}

func EnsureCrdExists(name string) (bool, error) {
	_, err := globalhelper.GetAPIClient().CustomResourceDefinitions().Get(context.TODO(),
		name, metav1.GetOptions{})

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Define a custom resource.
func DefineCustomResource(name, namespace, operatorLabels string, operatorLabelsMap map[string]string) *crscaleoperator.Memcached {
	return &crscaleoperator.Memcached{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cache.example.com/v1",
			Kind:       "Memcached",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    operatorLabelsMap,
		},
		Spec: crscaleoperator.MemcachedSpec{
			Size: 1,
		},
		Status: crscaleoperator.MemcachedStatus{
			Selector: operatorLabels,
		},
	}
}

func RedefineCustomResourceWithReplica(aCustomResource *crscaleoperator.Memcached, replicas int) {
	aCustomResource.Spec.Size = int32(replicas)
}

func CreateCustomResourceScale(name, namespace, operatorLabels string, operatorLabelsMap map[string]string) (string, error) {
	aCustomResource := DefineCustomResource(name, namespace, operatorLabels, operatorLabelsMap)

	body, err := json.Marshal(aCustomResource)

	if err != nil {
		return "", fmt.Errorf("error during marshaling the custom resource definition: %w", err)
	}

	data, err := globalhelper.GetAPIClient().CoreV1Interface.RESTClient().
		Post().AbsPath("/apis/cache.example.com/v1/namespaces/" + namespace + "/memcacheds").
		Body(body).DoRaw(context.TODO())

	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return "success", nil
		}

		return "", fmt.Errorf("return data %v and err %w", data, err)
	}

	return "success", nil
}
