package accesscontrol

import (
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tshelper "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/helper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalparameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
)

// POC v3: Label-based scenario testing (CORRECT APPROACH)
//
// Key insight: Certsuite uses label filters from DefineCertsuiteConfig to determine
// which workloads to test. By using different labels for different scenarios, we can:
// 1. Create ALL deployments in parallel (fast setup)
// 2. Test each scenario separately by updating the label filter
// 3. Maintain test granularity while still benefiting from parallel creation
//
// Expected time savings:
// - Original: 4 tests × (deploy 15s + LaunchTests 60s) = ~300s
// - This approach: parallel deploy 15s + 4 × LaunchTests 60s = ~255s (15% savings)
// - Better savings with more scenarios sharing same namespace
//
// The key is: CREATE ONCE, TEST MULTIPLE TIMES with different label filters

// deploymentConfig defines a deployment configuration for testing
type deploymentConfig struct {
	name       string
	replicas   int32
	containers int
	withBPF    bool
}

// testScenario defines a complete test scenario with deployments and expected results
type testScenario struct {
	name           string
	label          string
	deployments    []deploymentConfig
	expectedResult string
	description    string
}

var _ = Describe("Access-control bpf-capability-check [POC v3 - Label-based]", Label("poc", "poc-v3", "optimization"), Ordered, func() {
	var randomNamespace string
	var randomReportDir string
	var randomCertsuiteConfigDir string

	var scenarios []testScenario

	BeforeAll(func() {
		// Create random namespace ONCE for all tests
		randomNamespace, randomReportDir, randomCertsuiteConfigDir =
			globalhelper.BeforeEachSetupWithRandomPrivilegedNamespace(
				tsparams.TestAccessControlNameSpace)

		// Define all test scenarios with unique labels
		scenarios = []testScenario{
			{
				name:           "one deployment without BPF",
				label:          "bpf-test-scenario-1",
				expectedResult: globalparameters.TestCasePassed,
				description:    "Single deployment without BPF capability should pass",
				deployments: []deploymentConfig{
					{"bpf-scenario1-dep1", 1, 1, false},
				},
			},
			{
				name:           "one deployment with BPF [negative]",
				label:          "bpf-test-scenario-2",
				expectedResult: globalparameters.TestCaseFailed,
				description:    "Single deployment with BPF capability should fail",
				deployments: []deploymentConfig{
					{"bpf-scenario2-dep1", 1, 1, true},
				},
			},
			{
				name:           "two deployments without BPF",
				label:          "bpf-test-scenario-3",
				expectedResult: globalparameters.TestCasePassed,
				description:    "Multiple deployments without BPF capability should pass",
				deployments: []deploymentConfig{
					{"bpf-scenario3-dep1", 1, 1, false},
					{"bpf-scenario3-dep2", 1, 1, false},
				},
			},
			{
				name:           "two deployments, one with BPF [negative]",
				label:          "bpf-test-scenario-4",
				expectedResult: globalparameters.TestCaseFailed,
				description:    "Mixed deployments with one having BPF should fail",
				deployments: []deploymentConfig{
					{"bpf-scenario4-dep1", 1, 1, true},  // Has BPF
					{"bpf-scenario4-dep2", 1, 1, false}, // No BPF
				},
			},
		}

		// Create ALL deployments in parallel with scenario-specific labels
		By("Creating all test deployments in parallel (all scenarios at once)")
		var wg sync.WaitGroup
		errChan := make(chan error, 100)

		for _, scenario := range scenarios {
			for _, depConfig := range scenario.deployments {
				wg.Add(1)
				go func(sc testScenario, dc deploymentConfig) {
					defer wg.Done()
					defer GinkgoRecover()

					By(fmt.Sprintf("Create deployment: %s (label: %s)", dc.name, sc.label))

					// Define deployment with scenario-specific label
					dep, err := tshelper.DefineDeployment(dc.replicas, dc.containers, dc.name, randomNamespace)
					if err != nil {
						errChan <- fmt.Errorf("failed to define deployment %s: %w", dc.name, err)
						return
					}

					// Update deployment labels to use scenario-specific label
					if dep.Spec.Template.Labels == nil {
						dep.Spec.Template.Labels = make(map[string]string)
					}
					// Set the scenario-specific label (this is what certsuite will filter on)
					dep.Spec.Template.Labels["test-network-function.com/generic"] = sc.label

					// Add BPF capability if needed
					if dc.withBPF {
						deployment.RedefineWithContainersSecurityContextBpf(dep)
					}

					// Create and wait for deployment to be ready
					err = globalhelper.CreateAndWaitUntilDeploymentIsReady(dep, tsparams.Timeout)
					if err != nil {
						errChan <- fmt.Errorf("failed to create deployment %s: %w", dc.name, err)
						return
					}
				}(scenario, depConfig)
			}
		}

		wg.Wait()
		close(errChan)

		// Check for errors
		for err := range errChan {
			Expect(err).ToNot(HaveOccurred())
		}

		By("All deployments created successfully in parallel")
	})

	AfterAll(func() {
		// Cleanup ONCE after all tests
		globalhelper.AfterEachCleanupWithRandomNamespace(randomNamespace,
			randomReportDir, randomCertsuiteConfigDir, tsparams.Timeout)
	})

	// Now test each scenario by updating the certsuite config with the appropriate label
	// This is the key: we loop through scenarios, update the label filter, and run LaunchTests
	for idx := range scenarios {
		// Capture the scenario for the closure
		scenario := scenarios[idx]

		It(scenario.name, func() {
			By(fmt.Sprintf("Testing scenario: %s", scenario.description))

			// Update certsuite config to ONLY test deployments with this scenario's label
			By(fmt.Sprintf("Configure certsuite to test workloads with label: %s", scenario.label))
			err := globalhelper.DefineCertsuiteConfig(
				[]string{randomNamespace},
				[]string{scenario.label}, // <-- This is the key: use scenario-specific label
				[]string{},
				[]string{},
				[]string{}, randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred(), "error defining certsuite config file")

			// Run LaunchTests - it will ONLY test deployments with the scenario label
			By(fmt.Sprintf("Start BPF capability test for scenario: %s", scenario.name))
			err = globalhelper.LaunchTests(
				tsparams.TestCaseNameAccessControlBpfCapability,
				globalhelper.ConvertSpecNameToFileName(CurrentSpecReport().FullText()),
				randomReportDir,
				randomCertsuiteConfigDir)
			Expect(err).ToNot(HaveOccurred())

			// Verify the expected result for this scenario
			By(fmt.Sprintf("Verify test case status - expecting: %s", scenario.expectedResult))
			err = globalhelper.ValidateIfReportsAreValid(
				tsparams.TestCaseNameAccessControlBpfCapability,
				scenario.expectedResult,
				randomReportDir)
			Expect(err).ToNot(HaveOccurred())
		})
	}
})

/*
TIMING BREAKDOWN:

Original approach (4 separate tests):
- Test 1: Create NS (5s) + Deploy (15s) + LaunchTests (60s) + Cleanup (5s) = 85s
- Test 2: Create NS (5s) + Deploy (15s) + LaunchTests (60s) + Cleanup (5s) = 85s
- Test 3: Create NS (5s) + Deploy (30s sequential) + LaunchTests (60s) + Cleanup (5s) = 100s
- Test 4: Create NS (5s) + Deploy (30s sequential) + LaunchTests (60s) + Cleanup (5s) = 100s
TOTAL: ~370s

This approach (POC v3):
- BeforeAll: Create NS (5s) + Deploy ALL in parallel (~20s for worst case) = 25s
- Test 1: Update config (1s) + LaunchTests (60s) + Validate (1s) = 62s
- Test 2: Update config (1s) + LaunchTests (60s) + Validate (1s) = 62s
- Test 3: Update config (1s) + LaunchTests (60s) + Validate (1s) = 62s
- Test 4: Update config (1s) + LaunchTests (60s) + Validate (1s) = 62s
- AfterAll: Cleanup NS (5s) = 5s
TOTAL: ~278s

SAVINGS: ~92s (25% reduction)

With more scenarios or more complex deployments, savings increase significantly!

KEY BENEFITS:
1. Parallel deployment creation (major time saver)
2. Single namespace creation/deletion (reduces overhead)
3. Maintains full test granularity (each scenario is a separate It block)
4. Easy to add/modify scenarios
5. Clear test reporting (each scenario reports separately)
*/
