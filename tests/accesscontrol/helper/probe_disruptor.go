package helper

import (
	"context"
	"time"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klog "k8s.io/klog/v2"
)

// DisruptNodeLocalCertsuiteProbe repeatedly memory-pressures the certsuite-probe
// pod on nodeName so kubelet OOMKills it, matching lab worker-0 CrashLoopBackOff.
func DisruptNodeLocalCertsuiteProbe(ctx context.Context, nodeName string) {
	ticker := time.NewTicker(tsparams.ProbeDisruptInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			disruptProbePodsOnNode(ctx, nodeName)
		}
	}
}

func disruptProbePodsOnNode(ctx context.Context, nodeName string) {
	podList, err := globalhelper.GetAPIClient().Pods("").List(ctx, metav1.ListOptions{
		LabelSelector: tsparams.CertsuiteProbePodLabel,
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		klog.V(5).Infof("failed to list certsuite-probe pods: %v", err)

		return
	}

	for i := range podList.Items {
		probePod := podList.Items[i]

		_, execErr := globalhelper.ExecCommand(probePod, probeMemoryPressureCommand())
		if execErr != nil {
			klog.V(5).Infof("probe disrupt exec on %s/%s: %v", probePod.Namespace, probePod.Name, execErr)
		}
	}
}

func probeMemoryPressureCommand() []string {
	// Allocate until the 100M daemonset memory limit OOMKills the container.
	// Fall back to killing PID 1 if the image has no python/dd.
	const script = `python3 -c 'x=[]
while True:
    x.append(" "*1024*1024)' 2>/dev/null || \
dd if=/dev/zero of=/dev/shm/fill bs=1M count=200 2>/dev/null; kill -9 1`

	return []string{"/bin/sh", "-c", script}
}
