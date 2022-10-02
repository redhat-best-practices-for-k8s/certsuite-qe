package runtimeclass

import (
	v1 "k8s.io/api/node/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineRunTimeClass(name string) *v1.RuntimeClass {
	return &v1.RuntimeClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Handler: "runc",
	}
}
