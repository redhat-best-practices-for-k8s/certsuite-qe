/*
Copyright 2021 Red Hat OpenShift Data Foundation.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

//+kubebuilder:validation:Enum=flashsystemcluster.odf.ibm.com/v1alpha1;storagecluster.ocs.openshift.io/v1

// StorageKind specifies a vendor's supported kinds
type StorageKind string

const (
	//ConditionStorageSystemInvalid communicates if storagesystem CR is invalid
	ConditionStorageSystemInvalid conditionsv1.ConditionType = "StorageSystemInvalid"

	// ConditionVendorCsvReady communicates if CSV is ready
	ConditionVendorCsvReady conditionsv1.ConditionType = "VendorCsvReady"

	// ConditionVendorSystemPresent communicates if backend storage CR is present in the cluster
	ConditionVendorSystemPresent conditionsv1.ConditionType = "VendorSystemPresent"
)

// StorageSystemSpec defines the desired state of StorageSystem
type StorageSystemSpec struct {
	// Important: Run "make generate" to regenerate code after modifying this file

	//+kubebuilder:validation:Required
	// Name describes the name of managed storage vendor CR
	Name string `json:"name"`

	//+kubebuilder:validation:Required
	// Namespace describes the namespace of managed storage vendor CR
	Namespace string `json:"namespace"`

	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=storagecluster.ocs.openshift.io/v1
	// Kind describes the kind of storage vendor
	Kind StorageKind `json:"kind,omitempty"`
}

// StorageSystemStatus defines the observed state of StorageSystem
type StorageSystemStatus struct {
	// Important: Run "make generate" to regenerate code after modifying this file

	// Conditions describes the state of the StorageSystem resource.
	// +optional
	Conditions []conditionsv1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=storsys
//+kubebuilder:printcolumn:JSONPath=".spec.kind",name=storage-system-kind,type=string
//+kubebuilder:printcolumn:JSONPath=".spec.name",name=storage-system-name,type=string
//+operator-sdk:csv:customresourcedefinitions:resources={{StorageCluster,v1,storageclusters.ocs.openshift.io},{FlashSystemCluster,v1alpha1,flashsystemclusters.odf.ibm.com}}

// StorageSystem is the Schema for the storagesystems API
type StorageSystem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageSystemSpec   `json:"spec,omitempty"`
	Status StorageSystemStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StorageSystemList contains a list of StorageSystem
type StorageSystemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageSystem `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StorageSystem{}, &StorageSystemList{})
}
