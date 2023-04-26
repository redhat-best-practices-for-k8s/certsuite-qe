//go:build !utest

package performance

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/performance/parameters"
	_ "github.com/test-network-function/cnfcert-tests-verification/tests/performance/tests"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPerformance(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	_ = flag.Lookup("logtostderr").Value.Set("true")
	_ = flag.Lookup("v").Value.Set(globalhelper.Configuration.General.VerificationLogLevel)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = globalhelper.Configuration.GetReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "CNFCert performance tests", reporterConfig)
}

var _ = BeforeSuite(func() {

	By("Create namespace")
	err := namespaces.Create(tsparams.PerformanceNamespace, globalhelper.APIClient)
	Expect(err).ToNot(HaveOccurred())

	// Create service account and roles and roles binding
	err = ConfigurePrivilegedServiceAccount(tsparams.PerformanceNamespace)
	Expect(err).ToNot(HaveOccurred())

	By("Define TNF config file")
	err = globalhelper.DefineTnfConfig(
		[]string{tsparams.PerformanceNamespace},
		[]string{tsparams.TestPodLabel},
		[]string{},
		[]string{})
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {

	By(fmt.Sprintf("Remove %s namespace", tsparams.PerformanceNamespace))
	err := namespaces.DeleteAndWait(
		globalhelper.APIClient,
		tsparams.PerformanceNamespace,
		tsparams.WaitingTime,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Remove reports from reports directory")
	err = globalhelper.RemoveContentsFromReportDir()
	Expect(err).ToNot(HaveOccurred())
})

func ConfigurePrivilegedServiceAccount(namespace string) error {
	aRole, aRoleBinding, aServiceAccount := getPrivilegedServiceAccountObjects(namespace)
	// create role
	_, err := globalhelper.APIClient.RbacV1Interface.Roles(namespace).Create(context.TODO(), &aRole, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating role, err=%w", err)
	}

	// create rolebinding
	_, err = globalhelper.APIClient.RbacV1Interface.RoleBindings(namespace).Create(context.TODO(), &aRoleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating role bindings, err=%w", err)
	}

	// create service account
	_, err = globalhelper.APIClient.CoreV1Interface.ServiceAccounts(namespace).Create(context.TODO(),
		&aServiceAccount, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating service account, err=%w", err)
	}

	return nil
}

func getPrivilegedServiceAccountObjects(namespace string) (aRole rbacv1.Role,
	aRoleBinding rbacv1.RoleBinding, aServiceAccount corev1.ServiceAccount) {
	aRole = rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsparams.PriviledgedRoleName,
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{{
			APIGroups: []string{"*"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		},
		},
	}

	aRoleBinding = rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsparams.PriviledgedRoleName,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      tsparams.PriviledgedRoleName,
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     tsparams.PriviledgedRoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	aServiceAccount = corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsparams.PriviledgedRoleName,
			Namespace: namespace,
		},
	}

	return aRole, aRoleBinding, aServiceAccount
}
