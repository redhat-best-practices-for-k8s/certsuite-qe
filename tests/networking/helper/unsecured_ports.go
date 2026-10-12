package helper

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/networking/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	tlsSecretName      = "unsecured-ports-tls"
	nginxConfigMapName = "unsecured-ports-nginx"
	tlsCertMountPath   = "/etc/nginx/tls"
	nginxTLSConf       = `server {
    listen 8443 ssl;
    listen [::]:8443 ssl;
    ssl_certificate /etc/nginx/tls/tls.crt;
    ssl_certificate_key /etc/nginx/tls/tls.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    location / {
        return 200 'tls\n';
        add_header Content-Type text/plain;
    }
}
`
)

// AnnotatePodTemplate sets an annotation on the deployment pod template.
func AnnotatePodTemplate(dep *appsv1.Deployment, key, value string) {
	if dep.Spec.Template.Annotations == nil {
		dep.Spec.Template.Annotations = map[string]string{}
	}

	dep.Spec.Template.Annotations[key] = value
}

// SetTCPReadinessProbe adds a TCP readiness probe on the first container.
func SetTCPReadinessProbe(dep *appsv1.Deployment, port int32) {
	if len(dep.Spec.Template.Spec.Containers) == 0 {
		return
	}

	dep.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt32(port),
			},
		},
		InitialDelaySeconds: 1,
		PeriodSeconds:       2,
	}
}

// UniquePodNodeNames returns distinct node names for pods in a namespace.
func UniquePodNodeNames(namespace, labelSelector string) ([]string, error) {
	podList, err := globalhelper.GetAPIClient().Pods(namespace).List(
		context.TODO(),
		metav1.ListOptions{LabelSelector: labelSelector},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in %s: %w", namespace, err)
	}

	return uniqueNodeNamesFromPods(podList.Items), nil
}

func uniqueNodeNamesFromPods(pods []corev1.Pod) []string {
	nodesSeen := map[string]struct{}{}

	for i := range pods {
		nodeName := pods[i].Spec.NodeName
		if nodeName == "" {
			continue
		}

		nodesSeen[nodeName] = struct{}{}
	}

	names := make([]string, 0, len(nodesSeen))
	for nodeName := range nodesSeen {
		names = append(names, nodeName)
	}

	return names
}

// DefineHTTPNginxDeployment returns an nginx deployment listening on plaintext HTTP.
func DefineHTTPNginxDeployment(name, namespace string, replicaNumber int32) *appsv1.Deployment {
	dep := defineNginxBaseDeployment(name, namespace, replicaNumber)
	setContainerPort(dep, tsparams.HTTPListenPort)
	SetTCPReadinessProbe(dep, tsparams.HTTPListenPort)

	return dep
}

// DefineTLSNginxDeployment creates TLS materials and returns an nginx TLS deployment.
func DefineTLSNginxDeployment(name, namespace string, replicaNumber int32) (*appsv1.Deployment, error) {
	err := createNginxTLSMaterials(namespace)
	if err != nil {
		return nil, err
	}

	dep := defineNginxBaseDeployment(name, namespace, replicaNumber)
	setContainerPort(dep, tsparams.TLSListenPort)
	SetTCPReadinessProbe(dep, tsparams.TLSListenPort)
	mountNginxTLS(dep)

	return dep, nil
}

func defineNginxBaseDeployment(name, namespace string, replicaNumber int32) *appsv1.Deployment {
	dep := deployment.DefineDeployment(name, namespace, tsparams.NginxUnprivilegedImage, tsparams.TestDeploymentLabels)
	deployment.RedefineWithReplicaNumber(dep, replicaNumber)

	// Use the image entrypoint (nginx) instead of the default sleep INF command.
	_ = deployment.RedefineContainerCommand(dep, 0, nil)

	return dep
}

func setContainerPort(dep *appsv1.Deployment, port int32) {
	if len(dep.Spec.Template.Spec.Containers) == 0 {
		return
	}

	dep.Spec.Template.Spec.Containers[0].Ports = []corev1.ContainerPort{
		{
			ContainerPort: port,
			Protocol:      corev1.ProtocolTCP,
		},
	}
}

func mountNginxTLS(dep *appsv1.Deployment) {
	dep.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: tlsSecretName,
				},
			},
		},
		{
			Name: "conf",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: nginxConfigMapName,
					},
				},
			},
		},
	}

	if len(dep.Spec.Template.Spec.Containers) == 0 {
		return
	}

	dep.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "tls",
			MountPath: tlsCertMountPath,
			ReadOnly:  true,
		},
		{
			Name:      "conf",
			MountPath: "/etc/nginx/conf.d",
			ReadOnly:  true,
		},
	}
}

func createNginxTLSMaterials(namespace string) error {
	certPEM, keyPEM, err := generateSelfSignedCert()
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tlsSecretName,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			corev1.TLSCertKey:       certPEM,
			corev1.TLSPrivateKeyKey: keyPEM,
		},
	}

	_, err = globalhelper.GetAPIClient().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create TLS secret: %w", err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxConfigMapName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"default.conf": nginxTLSConf,
		},
	}

	_, err = globalhelper.GetAPIClient().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create nginx TLS configmap: %w", err)
	}

	return nil
}

func generateSelfSignedCert() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "unsecured-ports-tls",
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return certPEM, keyPEM, nil
}
