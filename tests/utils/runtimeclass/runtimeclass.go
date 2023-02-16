package runtimeclass

import (
	nodev1 "k8s.io/api/node/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineRunTimeClass(name string) *nodev1.RuntimeClass {
	return &nodev1.RuntimeClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Handler: "runc",
	}
}
