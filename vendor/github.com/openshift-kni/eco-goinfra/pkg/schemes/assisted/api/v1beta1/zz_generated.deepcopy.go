//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2020.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"github.com/openshift-kni/eco-goinfra/pkg/schemes/assisted/api/common"
	"github.com/openshift-kni/eco-goinfra/pkg/schemes/assisted/models"
	"github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Agent) DeepCopyInto(out *Agent) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Agent.
func (in *Agent) DeepCopy() *Agent {
	if in == nil {
		return nil
	}
	out := new(Agent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Agent) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentClassification) DeepCopyInto(out *AgentClassification) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentClassification.
func (in *AgentClassification) DeepCopy() *AgentClassification {
	if in == nil {
		return nil
	}
	out := new(AgentClassification)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AgentClassification) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentClassificationList) DeepCopyInto(out *AgentClassificationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AgentClassification, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentClassificationList.
func (in *AgentClassificationList) DeepCopy() *AgentClassificationList {
	if in == nil {
		return nil
	}
	out := new(AgentClassificationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AgentClassificationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentClassificationSpec) DeepCopyInto(out *AgentClassificationSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentClassificationSpec.
func (in *AgentClassificationSpec) DeepCopy() *AgentClassificationSpec {
	if in == nil {
		return nil
	}
	out := new(AgentClassificationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentClassificationStatus) DeepCopyInto(out *AgentClassificationStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentClassificationStatus.
func (in *AgentClassificationStatus) DeepCopy() *AgentClassificationStatus {
	if in == nil {
		return nil
	}
	out := new(AgentClassificationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentDeprovisionInfo) DeepCopyInto(out *AgentDeprovisionInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentDeprovisionInfo.
func (in *AgentDeprovisionInfo) DeepCopy() *AgentDeprovisionInfo {
	if in == nil {
		return nil
	}
	out := new(AgentDeprovisionInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentList) DeepCopyInto(out *AgentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Agent, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentList.
func (in *AgentList) DeepCopy() *AgentList {
	if in == nil {
		return nil
	}
	out := new(AgentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AgentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentServiceConfig) DeepCopyInto(out *AgentServiceConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentServiceConfig.
func (in *AgentServiceConfig) DeepCopy() *AgentServiceConfig {
	if in == nil {
		return nil
	}
	out := new(AgentServiceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AgentServiceConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentServiceConfigList) DeepCopyInto(out *AgentServiceConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AgentServiceConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentServiceConfigList.
func (in *AgentServiceConfigList) DeepCopy() *AgentServiceConfigList {
	if in == nil {
		return nil
	}
	out := new(AgentServiceConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AgentServiceConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentServiceConfigSpec) DeepCopyInto(out *AgentServiceConfigSpec) {
	*out = *in
	in.FileSystemStorage.DeepCopyInto(&out.FileSystemStorage)
	in.DatabaseStorage.DeepCopyInto(&out.DatabaseStorage)
	if in.ImageStorage != nil {
		in, out := &in.ImageStorage, &out.ImageStorage
		*out = new(corev1.PersistentVolumeClaimSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.MirrorRegistryRef != nil {
		in, out := &in.MirrorRegistryRef, &out.MirrorRegistryRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.OSImages != nil {
		in, out := &in.OSImages, &out.OSImages
		*out = make([]OSImage, len(*in))
		copy(*out, *in)
	}
	if in.MustGatherImages != nil {
		in, out := &in.MustGatherImages, &out.MustGatherImages
		*out = make([]MustGatherImage, len(*in))
		copy(*out, *in)
	}
	if in.UnauthenticatedRegistries != nil {
		in, out := &in.UnauthenticatedRegistries, &out.UnauthenticatedRegistries
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.OSImageCACertRef != nil {
		in, out := &in.OSImageCACertRef, &out.OSImageCACertRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.OSImageAdditionalParamsRef != nil {
		in, out := &in.OSImageAdditionalParamsRef, &out.OSImageAdditionalParamsRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(Ingress)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentServiceConfigSpec.
func (in *AgentServiceConfigSpec) DeepCopy() *AgentServiceConfigSpec {
	if in == nil {
		return nil
	}
	out := new(AgentServiceConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentServiceConfigStatus) DeepCopyInto(out *AgentServiceConfigStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentServiceConfigStatus.
func (in *AgentServiceConfigStatus) DeepCopy() *AgentServiceConfigStatus {
	if in == nil {
		return nil
	}
	out := new(AgentServiceConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentSpec) DeepCopyInto(out *AgentSpec) {
	*out = *in
	if in.ClusterDeploymentName != nil {
		in, out := &in.ClusterDeploymentName, &out.ClusterDeploymentName
		*out = new(ClusterReference)
		**out = **in
	}
	if in.IgnitionEndpointTokenReference != nil {
		in, out := &in.IgnitionEndpointTokenReference, &out.IgnitionEndpointTokenReference
		*out = new(IgnitionEndpointTokenReference)
		**out = **in
	}
	if in.IgnitionEndpointHTTPHeaders != nil {
		in, out := &in.IgnitionEndpointHTTPHeaders, &out.IgnitionEndpointHTTPHeaders
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.NodeLabels != nil {
		in, out := &in.NodeLabels, &out.NodeLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentSpec.
func (in *AgentSpec) DeepCopy() *AgentSpec {
	if in == nil {
		return nil
	}
	out := new(AgentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentStatus) DeepCopyInto(out *AgentStatus) {
	*out = *in
	in.Inventory.DeepCopyInto(&out.Inventory)
	in.Progress.DeepCopyInto(&out.Progress)
	if in.NtpSources != nil {
		in, out := &in.NtpSources, &out.NtpSources
		*out = make([]HostNTPSources, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.DebugInfo = in.DebugInfo
	if in.ValidationsInfo != nil {
		in, out := &in.ValidationsInfo, &out.ValidationsInfo
		*out = make(common.ValidationsStatus, len(*in))
		for key, val := range *in {
			var outVal []common.ValidationResult
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(common.ValidationResults, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.DeprovisionInfo != nil {
		in, out := &in.DeprovisionInfo, &out.DeprovisionInfo
		*out = new(AgentDeprovisionInfo)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentStatus.
func (in *AgentStatus) DeepCopy() *AgentStatus {
	if in == nil {
		return nil
	}
	out := new(AgentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BootArtifacts) DeepCopyInto(out *BootArtifacts) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BootArtifacts.
func (in *BootArtifacts) DeepCopy() *BootArtifacts {
	if in == nil {
		return nil
	}
	out := new(BootArtifacts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterReference) DeepCopyInto(out *ClusterReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterReference.
func (in *ClusterReference) DeepCopy() *ClusterReference {
	if in == nil {
		return nil
	}
	out := new(ClusterReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DebugInfo) DeepCopyInto(out *DebugInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DebugInfo.
func (in *DebugInfo) DeepCopy() *DebugInfo {
	if in == nil {
		return nil
	}
	out := new(DebugInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostBoot) DeepCopyInto(out *HostBoot) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostBoot.
func (in *HostBoot) DeepCopy() *HostBoot {
	if in == nil {
		return nil
	}
	out := new(HostBoot)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostCPU) DeepCopyInto(out *HostCPU) {
	*out = *in
	if in.Flags != nil {
		in, out := &in.Flags, &out.Flags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostCPU.
func (in *HostCPU) DeepCopy() *HostCPU {
	if in == nil {
		return nil
	}
	out := new(HostCPU)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostDisk) DeepCopyInto(out *HostDisk) {
	*out = *in
	in.InstallationEligibility.DeepCopyInto(&out.InstallationEligibility)
	out.IoPerf = in.IoPerf
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostDisk.
func (in *HostDisk) DeepCopy() *HostDisk {
	if in == nil {
		return nil
	}
	out := new(HostDisk)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostIOPerf) DeepCopyInto(out *HostIOPerf) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostIOPerf.
func (in *HostIOPerf) DeepCopy() *HostIOPerf {
	if in == nil {
		return nil
	}
	out := new(HostIOPerf)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostInstallationEligibility) DeepCopyInto(out *HostInstallationEligibility) {
	*out = *in
	if in.NotEligibleReasons != nil {
		in, out := &in.NotEligibleReasons, &out.NotEligibleReasons
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostInstallationEligibility.
func (in *HostInstallationEligibility) DeepCopy() *HostInstallationEligibility {
	if in == nil {
		return nil
	}
	out := new(HostInstallationEligibility)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostInterface) DeepCopyInto(out *HostInterface) {
	*out = *in
	if in.IPV6Addresses != nil {
		in, out := &in.IPV6Addresses, &out.IPV6Addresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.IPV4Addresses != nil {
		in, out := &in.IPV4Addresses, &out.IPV4Addresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Flags != nil {
		in, out := &in.Flags, &out.Flags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostInterface.
func (in *HostInterface) DeepCopy() *HostInterface {
	if in == nil {
		return nil
	}
	out := new(HostInterface)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostInventory) DeepCopyInto(out *HostInventory) {
	*out = *in
	if in.ReportTime != nil {
		in, out := &in.ReportTime, &out.ReportTime
		*out = (*in).DeepCopy()
	}
	out.Memory = in.Memory
	in.Cpu.DeepCopyInto(&out.Cpu)
	if in.Interfaces != nil {
		in, out := &in.Interfaces, &out.Interfaces
		*out = make([]HostInterface, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Disks != nil {
		in, out := &in.Disks, &out.Disks
		*out = make([]HostDisk, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Boot = in.Boot
	out.SystemVendor = in.SystemVendor
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostInventory.
func (in *HostInventory) DeepCopy() *HostInventory {
	if in == nil {
		return nil
	}
	out := new(HostInventory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostMemory) DeepCopyInto(out *HostMemory) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostMemory.
func (in *HostMemory) DeepCopy() *HostMemory {
	if in == nil {
		return nil
	}
	out := new(HostMemory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostNTPSources) DeepCopyInto(out *HostNTPSources) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostNTPSources.
func (in *HostNTPSources) DeepCopy() *HostNTPSources {
	if in == nil {
		return nil
	}
	out := new(HostNTPSources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostProgressInfo) DeepCopyInto(out *HostProgressInfo) {
	*out = *in
	if in.ProgressStages != nil {
		in, out := &in.ProgressStages, &out.ProgressStages
		*out = make([]models.HostStage, len(*in))
		copy(*out, *in)
	}
	if in.StageStartTime != nil {
		in, out := &in.StageStartTime, &out.StageStartTime
		*out = (*in).DeepCopy()
	}
	if in.StageUpdateTime != nil {
		in, out := &in.StageUpdateTime, &out.StageUpdateTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostProgressInfo.
func (in *HostProgressInfo) DeepCopy() *HostProgressInfo {
	if in == nil {
		return nil
	}
	out := new(HostProgressInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostSystemVendor) DeepCopyInto(out *HostSystemVendor) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostSystemVendor.
func (in *HostSystemVendor) DeepCopy() *HostSystemVendor {
	if in == nil {
		return nil
	}
	out := new(HostSystemVendor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HypershiftAgentServiceConfig) DeepCopyInto(out *HypershiftAgentServiceConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HypershiftAgentServiceConfig.
func (in *HypershiftAgentServiceConfig) DeepCopy() *HypershiftAgentServiceConfig {
	if in == nil {
		return nil
	}
	out := new(HypershiftAgentServiceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HypershiftAgentServiceConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HypershiftAgentServiceConfigList) DeepCopyInto(out *HypershiftAgentServiceConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HypershiftAgentServiceConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HypershiftAgentServiceConfigList.
func (in *HypershiftAgentServiceConfigList) DeepCopy() *HypershiftAgentServiceConfigList {
	if in == nil {
		return nil
	}
	out := new(HypershiftAgentServiceConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HypershiftAgentServiceConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HypershiftAgentServiceConfigSpec) DeepCopyInto(out *HypershiftAgentServiceConfigSpec) {
	*out = *in
	in.AgentServiceConfigSpec.DeepCopyInto(&out.AgentServiceConfigSpec)
	out.KubeconfigSecretRef = in.KubeconfigSecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HypershiftAgentServiceConfigSpec.
func (in *HypershiftAgentServiceConfigSpec) DeepCopy() *HypershiftAgentServiceConfigSpec {
	if in == nil {
		return nil
	}
	out := new(HypershiftAgentServiceConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HypershiftAgentServiceConfigStatus) DeepCopyInto(out *HypershiftAgentServiceConfigStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HypershiftAgentServiceConfigStatus.
func (in *HypershiftAgentServiceConfigStatus) DeepCopy() *HypershiftAgentServiceConfigStatus {
	if in == nil {
		return nil
	}
	out := new(HypershiftAgentServiceConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IgnitionEndpointTokenReference) DeepCopyInto(out *IgnitionEndpointTokenReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IgnitionEndpointTokenReference.
func (in *IgnitionEndpointTokenReference) DeepCopy() *IgnitionEndpointTokenReference {
	if in == nil {
		return nil
	}
	out := new(IgnitionEndpointTokenReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InfraEnv) DeepCopyInto(out *InfraEnv) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InfraEnv.
func (in *InfraEnv) DeepCopy() *InfraEnv {
	if in == nil {
		return nil
	}
	out := new(InfraEnv)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InfraEnv) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InfraEnvDebugInfo) DeepCopyInto(out *InfraEnvDebugInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InfraEnvDebugInfo.
func (in *InfraEnvDebugInfo) DeepCopy() *InfraEnvDebugInfo {
	if in == nil {
		return nil
	}
	out := new(InfraEnvDebugInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InfraEnvList) DeepCopyInto(out *InfraEnvList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InfraEnv, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InfraEnvList.
func (in *InfraEnvList) DeepCopy() *InfraEnvList {
	if in == nil {
		return nil
	}
	out := new(InfraEnvList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InfraEnvList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InfraEnvSpec) DeepCopyInto(out *InfraEnvSpec) {
	*out = *in
	if in.Proxy != nil {
		in, out := &in.Proxy, &out.Proxy
		*out = new(Proxy)
		**out = **in
	}
	if in.AdditionalNTPSources != nil {
		in, out := &in.AdditionalNTPSources, &out.AdditionalNTPSources
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PullSecretRef != nil {
		in, out := &in.PullSecretRef, &out.PullSecretRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.AgentLabels != nil {
		in, out := &in.AgentLabels, &out.AgentLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.NMStateConfigLabelSelector.DeepCopyInto(&out.NMStateConfigLabelSelector)
	if in.ClusterRef != nil {
		in, out := &in.ClusterRef, &out.ClusterRef
		*out = new(ClusterReference)
		**out = **in
	}
	if in.KernelArguments != nil {
		in, out := &in.KernelArguments, &out.KernelArguments
		*out = make([]KernelArgument, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InfraEnvSpec.
func (in *InfraEnvSpec) DeepCopy() *InfraEnvSpec {
	if in == nil {
		return nil
	}
	out := new(InfraEnvSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InfraEnvStatus) DeepCopyInto(out *InfraEnvStatus) {
	*out = *in
	if in.CreatedTime != nil {
		in, out := &in.CreatedTime, &out.CreatedTime
		*out = (*in).DeepCopy()
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.AgentLabelSelector.DeepCopyInto(&out.AgentLabelSelector)
	out.InfraEnvDebugInfo = in.InfraEnvDebugInfo
	out.BootArtifacts = in.BootArtifacts
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InfraEnvStatus.
func (in *InfraEnvStatus) DeepCopy() *InfraEnvStatus {
	if in == nil {
		return nil
	}
	out := new(InfraEnvStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
	if in.ClassName != nil {
		in, out := &in.ClassName, &out.ClassName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Interface) DeepCopyInto(out *Interface) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Interface.
func (in *Interface) DeepCopy() *Interface {
	if in == nil {
		return nil
	}
	out := new(Interface)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KernelArgument) DeepCopyInto(out *KernelArgument) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KernelArgument.
func (in *KernelArgument) DeepCopy() *KernelArgument {
	if in == nil {
		return nil
	}
	out := new(KernelArgument)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MustGatherImage) DeepCopyInto(out *MustGatherImage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MustGatherImage.
func (in *MustGatherImage) DeepCopy() *MustGatherImage {
	if in == nil {
		return nil
	}
	out := new(MustGatherImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NMStateConfig) DeepCopyInto(out *NMStateConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NMStateConfig.
func (in *NMStateConfig) DeepCopy() *NMStateConfig {
	if in == nil {
		return nil
	}
	out := new(NMStateConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NMStateConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NMStateConfigList) DeepCopyInto(out *NMStateConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NMStateConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NMStateConfigList.
func (in *NMStateConfigList) DeepCopy() *NMStateConfigList {
	if in == nil {
		return nil
	}
	out := new(NMStateConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NMStateConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NMStateConfigSpec) DeepCopyInto(out *NMStateConfigSpec) {
	*out = *in
	if in.Interfaces != nil {
		in, out := &in.Interfaces, &out.Interfaces
		*out = make([]*Interface, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Interface)
				**out = **in
			}
		}
	}
	in.NetConfig.DeepCopyInto(&out.NetConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NMStateConfigSpec.
func (in *NMStateConfigSpec) DeepCopy() *NMStateConfigSpec {
	if in == nil {
		return nil
	}
	out := new(NMStateConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetConfig) DeepCopyInto(out *NetConfig) {
	*out = *in
	if in.Raw != nil {
		in, out := &in.Raw, &out.Raw
		*out = make(RawNetConfig, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetConfig.
func (in *NetConfig) DeepCopy() *NetConfig {
	if in == nil {
		return nil
	}
	out := new(NetConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OSImage) DeepCopyInto(out *OSImage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OSImage.
func (in *OSImage) DeepCopy() *OSImage {
	if in == nil {
		return nil
	}
	out := new(OSImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Proxy) DeepCopyInto(out *Proxy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Proxy.
func (in *Proxy) DeepCopy() *Proxy {
	if in == nil {
		return nil
	}
	out := new(Proxy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in RawNetConfig) DeepCopyInto(out *RawNetConfig) {
	{
		in := &in
		*out = make(RawNetConfig, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RawNetConfig.
func (in RawNetConfig) DeepCopy() RawNetConfig {
	if in == nil {
		return nil
	}
	out := new(RawNetConfig)
	in.DeepCopyInto(out)
	return *out
}
