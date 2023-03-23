package tests

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	tshelper "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/helper"
	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/platformalteration/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	istioNamespace = "istio-system"
	istioCR        = "installed-state"
	waitingTime    = 2 * time.Minute
)

/*
*

	Precondition :
		The Istio Operator needs to be pre-installed.
		We are checking if Istio Operator is installed or not.

*
*/

// Not tested yet
func createServiceMesh(filename string) (bool, error) {
	bytesInFile, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal(err)

		return false, err
	}

	log.Printf("%q \n", string(bytesInFile))

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(bytesInFile), 100)

	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		if len(rawObj.Raw) == 0 {
			// if the yaml object is empty just continue to the next one
			continue
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			log.Fatal(err)

			return false, err
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Fatal(err)

			return false, err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		groupResource, err := restmapper.GetAPIGroupResources(globalhelper.APIClient.K8sClient.Discovery())
		if err != nil {
			log.Fatal(err)

			return false, err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(groupResource)

		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatal(err)

			return false, err
		}

		var dri dynamic.ResourceInterface

		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}

			dri = globalhelper.APIClient.DynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = globalhelper.APIClient.DynamicClient.Resource(mapping.Resource)
		}

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			log.Fatal(err)

			return false, err
		}
	}

	return true, nil
}

var _ = Describe("platform-alteration-service-mesh-usage", func() {

	BeforeEach(func() {
		By("Clean namespace before each test")
		err := namespaces.Clean(tsparams.PlatformAlterationNamespace, globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		// createServiceMesh("istio.yaml")
	})

	// 56594
	FIt("istio is installed", func() {

		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.Configuration.General.TestImage)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, waitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		logrus.Error(err)
		Expect(err).ToNot(HaveOccurred())

		//err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCasePassed)
		//Expect(err).ToNot(HaveOccurred())
	})

	// 56596
	FIt("istio is installed but proxy containers does not exist [negative]", func() {

		By("Check Istio resource exists")
		gvr := schema.GroupVersionResource{Group: "install.istio.io", Version: "v1alpha1", Resource: "istiooperators"}

		cr, err := globalhelper.APIClient.DynamicClient.Resource(gvr).Namespace(istioNamespace).Get(context.TODO(),
			"example-istiocontrolplane", metav1.GetOptions{})

		Expect(err).ToNot(HaveOccurred())
		Expect(cr).NotTo(BeNil())

		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(put, waitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))

		Expect(err).To(HaveOccurred())

		//err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseFailed)
		//Expect(err).ToNot(HaveOccurred())
	})

	// 56597
	It("istio is installed but proxy container exist on one pod only [negative]", func() {

		By("Define first pod with instio container")
		put := pod.DefinePod(tsparams.TestPodName, tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)
		tshelper.AppendIstioContainerToPod(put, globalhelper.Configuration.General.TestImage)

		err := globalhelper.CreateAndWaitUntilPodIsReady(put, waitingTime)
		Expect(err).ToNot(HaveOccurred())

		putb := pod.DefinePod("lifecycle-putb", tsparams.PlatformAlterationNamespace, globalhelper.Configuration.General.TestImage,
			tsparams.TnfTargetPodLabels)

		err = globalhelper.CreateAndWaitUntilPodIsReady(putb, waitingTime)
		Expect(err).ToNot(HaveOccurred())

		By("Start platform-alteration-service-mesh-usage test")
		err = globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).To(HaveOccurred())

		//err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseFailed)
		//Expect(err).ToNot(HaveOccurred())
	})

	// TODO: Change this code
	// 56595
	It("istio is not installed", func() {
		By("Start platform-alteration-service-mesh-usage test")
		err := globalhelper.LaunchTests(tsparams.TnfServiceMeshUsageName,
			globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()))
		Expect(err).ToNot(HaveOccurred())

		//err = globalhelper.ValidateIfReportsAreValid(tsparams.TnfServiceMeshUsageName, globalparameters.TestCaseSkipped)
		//Expect(err).ToNot(HaveOccurred())
	})
})
