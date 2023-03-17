package nad

import (
	"fmt"
	"strings"

	netattdefv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefineNad returns basic network-attachment-definition manifest.
func DefineNad(name string, namespace string) *netattdefv1.NetworkAttachmentDefinition {
	return &netattdefv1.NetworkAttachmentDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: netattdefv1.NetworkAttachmentDefinitionSpec{
			Config: fmt.Sprintf(
				`{"cniVersion": "0.4.0", "name": "%s", "type": "macvlan", "mode": "bridge"}`,
				name),
		},
	}
}

// RedefineNadWithWhereaboutsIpam updates nad with whereabouts ipam config.
func RedefineNadWithWhereaboutsIpam(
	nad *netattdefv1.NetworkAttachmentDefinition, network string) *netattdefv1.NetworkAttachmentDefinition {
	nad.Spec.Config = strings.Trim(nad.Spec.Config, `}`)
	nad.Spec.Config = fmt.Sprintf(
		"%s, %s}", nad.Spec.Config,
		fmt.Sprintf(`"ipam":{ "type": "whereabouts", "range": "%s"}`, network),
	)

	return nad
}
