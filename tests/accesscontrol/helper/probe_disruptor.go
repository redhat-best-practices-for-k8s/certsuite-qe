package helper

import (
	"context"
	"time"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	klog "k8s.io/klog/v2"
)

// DisruptNodeLocalCertsuiteProbe waits until the certsuite-probe DaemonSet is
// fully Ready, pauses so certsuite's WaitDaemonsetReady poll can observe that
// state, then SIGKILLs the Running probe on nodeName until ctx is cancelled.
//
// The probe DaemonSet is created during LaunchTests. Killing a single Ready
// pod too early drops NumberReady and can make WaitDaemonsetReady time out,
// which skips ssh-daemons instead of failing it as a probe-exec outage.
func DisruptNodeLocalCertsuiteProbe(ctx context.Context, nodeName string) {
	err := wait.PollUntilContextCancel(ctx, tsparams.ProbeDisruptPollInterval, true,
		func(ctx context.Context) (bool, error) {
			ready, err := probeDaemonSetReady(ctx)
			if err != nil {
				klog.V(5).Infof("failed to get certsuite-probe daemonset: %v", err)

				return false, nil
			}

			return ready, nil
		})
	if err != nil {
		klog.Infof("gave up waiting for certsuite-probe daemonset: %v", err)

		return
	}

	klog.Infof("certsuite-probe daemonset is Ready; waiting %s before disrupting node %s",
		tsparams.ProbeDisruptReadyGracePeriod, nodeName)

	timer := time.NewTimer(tsparams.ProbeDisruptReadyGracePeriod)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
	}

	klog.Infof("starting node-local probe disrupt loop on %s", nodeName)

	ticker := time.NewTicker(tsparams.ProbeDisruptPollInterval)
	defer ticker.Stop()

	disruptRunningProbeOnNode(ctx, nodeName)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			disruptRunningProbeOnNode(ctx, nodeName)
		}
	}
}

func disruptRunningProbeOnNode(ctx context.Context, nodeName string) {
	pod, err := findProbePodOnNode(ctx, nodeName)
	if err != nil {
		klog.V(5).Infof("failed to list certsuite-probe pods: %v", err)

		return
	}

	if pod == nil {
		return
	}

	klog.V(5).Infof("disrupting certsuite-probe %s/%s on node %s", pod.Namespace, pod.Name, nodeName)

	if _, err := globalhelper.ExecCommand(*pod, probeDisruptCommand()); err != nil {
		klog.V(5).Infof("probe disrupt exec on %s/%s: %v", pod.Namespace, pod.Name, err)
	}
}

func probeDaemonSetReady(ctx context.Context) (bool, error) {
	daemonSet, err := globalhelper.GetAPIClient().DaemonSets(tsparams.CertsuiteProbeDaemonSetNamespace).Get(
		ctx, tsparams.CertsuiteProbeDaemonSetName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return isProbeDaemonSetReady(&daemonSet.Status), nil
}

func isProbeDaemonSetReady(status *appsv1.DaemonSetStatus) bool {
	if status == nil || status.DesiredNumberScheduled <= 0 {
		return false
	}

	desired := status.DesiredNumberScheduled

	return status.CurrentNumberScheduled == desired &&
		status.NumberAvailable == desired &&
		status.NumberReady == desired &&
		status.NumberMisscheduled == 0
}

func findProbePodOnNode(ctx context.Context, nodeName string) (*corev1.Pod, error) {
	podList, err := globalhelper.GetAPIClient().Pods(tsparams.CertsuiteProbeDaemonSetNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: tsparams.CertsuiteProbePodLabel,
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return nil, err
	}

	return firstMatchingProbe(podList.Items, nodeName), nil
}

func firstMatchingProbe(pods []corev1.Pod, nodeName string) *corev1.Pod {
	for i := range pods {
		pod := &pods[i]
		if pod.Spec.NodeName != nodeName || pod.Status.Phase != corev1.PodRunning {
			continue
		}

		return pod
	}

	return nil
}

func probeDisruptCommand() []string {
	// The probe image has no python3 and QE does not set a memory limit, so
	// allocate-until-OOM never trips. SIGKILL PID 1 after a short sleep is
	// reliable on UBI; the Go-side DaemonSet-ready grace already covers
	// WaitDaemonsetReady.
	const script = `sleep 1; kill -9 1`

	return []string{"/bin/sh", "-c", script}
}
