package nad

import (
	"fmt"
	"os"
	"strings"

	netattdefv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefineNad returns basic network-attachment-definition manifest.
func DefineNad(name string, namespace string) *netattdefv1.NetworkAttachmentDefinition {
	// Optionally include a master interface for macvlan if provided via env
	// This helps environments (e.g., microshift) where macvlan requires explicit master
	master := os.Getenv("CERTSUITE_NAD_MASTER_IF")
	baseCfg := fmt.Sprintf(`{"cniVersion": "0.4.0", "name": "%s", "type": "macvlan", "mode": "bridge"}`, name)
	if master != "" {
		// inject master right after the type/mode block
		baseCfg = strings.TrimRight(baseCfg, "}") + fmt.Sprintf(", \"master\": \"%s\"}", master)
	}

	return &netattdefv1.NetworkAttachmentDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: netattdefv1.NetworkAttachmentDefinitionSpec{
			Config: baseCfg,
		},
	}
}

// RedefineNadWithWhereaboutsIpam updates nad with whereabouts ipam config.
func RedefineNadWithWhereaboutsIpam(
	nad *netattdefv1.NetworkAttachmentDefinition, network string) {
	nad.Spec.Config = strings.Trim(nad.Spec.Config, `}`)
	// Ensure master is present if provided via env and not already set
	if !strings.Contains(nad.Spec.Config, `"master"`) {
		if m := os.Getenv("CERTSUITE_NAD_MASTER_IF"); m != "" {
			nad.Spec.Config = fmt.Sprintf(`%s, "master": "%s"`, nad.Spec.Config, m)
		}
	}
	nad.Spec.Config = fmt.Sprintf(
		"%s, %s}", nad.Spec.Config,
		fmt.Sprintf(`"ipam":{ "type": "whereabouts", "range": "%s"}`, network),
	)
}

// AddMasterToNad injects a macvlan master if not already present
func AddMasterToNad(n *netattdefv1.NetworkAttachmentDefinition, master string) {
	if master == "" || strings.Contains(n.Spec.Config, `"master"`) {
		return
	}
	n.Spec.Config = strings.TrimRight(n.Spec.Config, "}") + fmt.Sprintf(`, "master": "%s"}`, master)
}
