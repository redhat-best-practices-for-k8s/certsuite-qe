package crd

import (
	"context"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	_, err := globalhelper.GetAPIClient().CustomResourceDefinitions().Get(context.Background(),
		name, metav1.GetOptions{})

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
