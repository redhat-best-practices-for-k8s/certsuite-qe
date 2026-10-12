package helper

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnnotatePodTemplate(t *testing.T) {
	t.Parallel()

	dep := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{},
		},
	}

	AnnotatePodTemplate(dep, "certsuite.redhat.com/non-tls-ports", "8080")
	assert.Equal(t, "8080", dep.Spec.Template.Annotations["certsuite.redhat.com/non-tls-ports"])

	AnnotatePodTemplate(dep, "certsuite.redhat.com/non-tls-ports", "8080,8443")
	assert.Equal(t, "8080,8443", dep.Spec.Template.Annotations["certsuite.redhat.com/non-tls-ports"])
}

func TestUniqueNodeNamesFromPods(t *testing.T) {
	t.Parallel()

	pods := []corev1.Pod{
		{Spec: corev1.PodSpec{NodeName: "worker-0"}},
		{Spec: corev1.PodSpec{NodeName: "worker-1"}},
		{Spec: corev1.PodSpec{NodeName: "worker-0"}},
		{Spec: corev1.PodSpec{NodeName: ""}},
	}

	names := uniqueNodeNamesFromPods(pods)
	assert.Len(t, names, 2)
	assert.ElementsMatch(t, []string{"worker-0", "worker-1"}, names)
}

func TestGenerateSelfSignedCert(t *testing.T) {
	t.Parallel()

	certPEM, keyPEM, err := generateSelfSignedCert()
	assert.NoError(t, err)

	certBlock, _ := pem.Decode(certPEM)
	assert.NotNil(t, certBlock)
	assert.Equal(t, "CERTIFICATE", certBlock.Type)

	keyBlock, _ := pem.Decode(keyPEM)
	assert.NotNil(t, keyBlock)
	assert.Equal(t, "RSA PRIVATE KEY", keyBlock.Type)

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, "unsecured-ports-tls", cert.Subject.CommonName)
	assert.True(t, bytes.Contains(certPEM, []byte("BEGIN CERTIFICATE")))
}

func TestSetTCPReadinessProbe(t *testing.T) {
	t.Parallel()

	dep := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "nginx"}},
				},
			},
		},
	}

	SetTCPReadinessProbe(dep, 8080)
	assert.NotNil(t, dep.Spec.Template.Spec.Containers[0].ReadinessProbe)
	assert.NotNil(t, dep.Spec.Template.Spec.Containers[0].ReadinessProbe.TCPSocket)
	assert.Equal(t, int32(8080), dep.Spec.Template.Spec.Containers[0].ReadinessProbe.TCPSocket.Port.IntVal)
}

func TestEmptyDeploymentHelpers(t *testing.T) {
	t.Parallel()

	dep := &appsv1.Deployment{}
	SetTCPReadinessProbe(dep, 8080)
	setContainerPort(dep, 8080)
	mountNginxTLS(dep)

	assert.Empty(t, dep.Spec.Template.Spec.Containers)
	assert.Len(t, dep.Spec.Template.Spec.Volumes, 2)
}

func TestDefineHTTPNginxDeployment(t *testing.T) {
	t.Parallel()

	dep := DefineHTTPNginxDeployment("http-dep", "ns", 1)
	assert.NotNil(t, dep)
	assert.Equal(t, "http-dep", dep.Name)
	assert.Len(t, dep.Spec.Template.Spec.Containers, 1)
	assert.Nil(t, dep.Spec.Template.Spec.Containers[0].Command)
}

func TestAnnotatePodTemplateNilSafe(t *testing.T) {
	t.Parallel()

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "d"},
	}
	AnnotatePodTemplate(dep, "k", "v")
	assert.Equal(t, "v", dep.Spec.Template.Annotations["k"])
}
