#!/bin/bash
# Performance Test Diagnostic Script for Certsuite QE
# Compares cluster configuration across OCP versions to identify why
# performance-exclusive-cpu-pool-rt-scheduling-policy and
# performance-shared-cpu-pool-non-rt-scheduling-policy tests are failing.
#
# Usage: ./diagnose-performance-tests.sh [output-file]
# Run this on each cluster (4.14, 4.16, 4.17) and compare results.

set -o pipefail

OUTPUT_FILE="${1:-/tmp/certsuite-perf-diag-$(date +%Y%m%d-%H%M%S).txt}"
TEST_NS="certsuite-diag-$$"
TEST_POD_EXCLUSIVE="diag-exclusive-cpu"
TEST_POD_SHARED="diag-shared-cpu"

# Colors for terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "$1" | tee -a "$OUTPUT_FILE"
}

section() {
    log "\n${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    log "${BLUE}  $1${NC}"
    log "${BLUE}═══════════════════════════════════════════════════════════════${NC}\n"
}

check() {
    log "${YELLOW}▶ $1${NC}"
}

success() {
    log "${GREEN}✓ $1${NC}"
}

error() {
    log "${RED}✗ $1${NC}"
}

info() {
    log "  $1"
}

cleanup() {
    log "\nCleaning up test resources..."
    oc delete namespace "$TEST_NS" --ignore-not-found --wait=false 2>/dev/null
}

trap cleanup EXIT

# Initialize output file
echo "Certsuite Performance Test Diagnostic Report" > "$OUTPUT_FILE"
echo "Generated: $(date)" >> "$OUTPUT_FILE"
echo "=============================================" >> "$OUTPUT_FILE"

section "1. CLUSTER INFORMATION"

check "Cluster Version"
CLUSTER_VERSION=$(oc get clusterversion version -o jsonpath='{.status.desired.version}' 2>/dev/null)
if [ -n "$CLUSTER_VERSION" ]; then
    success "OCP Version: $CLUSTER_VERSION"
else
    error "Could not determine cluster version"
fi

check "Cluster ID"
CLUSTER_ID=$(oc get clusterversion version -o jsonpath='{.spec.clusterID}' 2>/dev/null)
info "Cluster ID: ${CLUSTER_ID:-unknown}"

check "API Server URL"
API_URL=$(oc whoami --show-server 2>/dev/null)
info "API: ${API_URL:-unknown}"

section "2. NODE CONFIGURATION"

check "Node Count and Roles"
oc get nodes -o wide 2>/dev/null | tee -a "$OUTPUT_FILE"

WORKER_COUNT=$(oc get nodes -l node-role.kubernetes.io/worker --no-headers 2>/dev/null | wc -l)
MASTER_COUNT=$(oc get nodes -l node-role.kubernetes.io/master --no-headers 2>/dev/null | wc -l)
info "Workers: $WORKER_COUNT, Masters: $MASTER_COUNT"

if [ "$WORKER_COUNT" -eq 0 ]; then
    error "No dedicated worker nodes found - this may be a compact cluster"
    info "Compact clusters run workloads on master nodes which may affect test behavior"
fi

check "Node Allocatable Resources"
oc get nodes -o custom-columns='NAME:.metadata.name,CPU:.status.allocatable.cpu,MEMORY:.status.allocatable.memory' 2>/dev/null | tee -a "$OUTPUT_FILE"

section "3. PERFORMANCE PROFILE & CPU MANAGER"

check "PerformanceProfile CRD"
if oc get crd performanceprofiles.performance.openshift.io &>/dev/null; then
    success "PerformanceProfile CRD exists"

    check "PerformanceProfile Resources"
    PP_COUNT=$(oc get performanceprofiles --no-headers 2>/dev/null | wc -l)
    if [ "$PP_COUNT" -gt 0 ]; then
        success "Found $PP_COUNT PerformanceProfile(s)"
        oc get performanceprofiles -o wide 2>/dev/null | tee -a "$OUTPUT_FILE"

        log "\nPerformanceProfile Details:"
        oc get performanceprofiles -o yaml 2>/dev/null | grep -A 50 "spec:" | head -60 | tee -a "$OUTPUT_FILE"
    else
        error "No PerformanceProfiles configured"
        info "This is required for exclusive CPU pool tests"
    fi
else
    error "PerformanceProfile CRD not found"
    info "Node Tuning Operator with performance profile support may not be installed"
fi

check "Node Tuning Operator"
oc get csv -n openshift-cluster-node-tuning-operator 2>/dev/null | tee -a "$OUTPUT_FILE"

check "CPU Manager State (from first worker node)"
WORKER_NODE=$(oc get nodes -l node-role.kubernetes.io/worker -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
if [ -z "$WORKER_NODE" ]; then
    WORKER_NODE=$(oc get nodes -l node-role.kubernetes.io/master -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    info "Using master node for checks: $WORKER_NODE"
else
    info "Using worker node: $WORKER_NODE"
fi

if [ -n "$WORKER_NODE" ]; then
    CPU_MANAGER_STATE=$(oc debug node/"$WORKER_NODE" -- chroot /host cat /var/lib/kubelet/cpu_manager_state 2>/dev/null)
    if [ -n "$CPU_MANAGER_STATE" ]; then
        info "CPU Manager State:"
        echo "$CPU_MANAGER_STATE" | head -5 | tee -a "$OUTPUT_FILE"

        if echo "$CPU_MANAGER_STATE" | grep -q '"policyName":"static"'; then
            success "CPU Manager Policy: static (required for exclusive CPUs)"
        elif echo "$CPU_MANAGER_STATE" | grep -q '"policyName":"none"'; then
            error "CPU Manager Policy: none (exclusive CPUs not available)"
        fi
    else
        error "Could not read CPU manager state"
    fi
fi

section "4. KERNEL & RUNTIME CONFIGURATION"

check "Kernel Version"
if [ -n "$WORKER_NODE" ]; then
    KERNEL=$(oc debug node/"$WORKER_NODE" -- chroot /host uname -r 2>/dev/null)
    info "Kernel: $KERNEL"

    if echo "$KERNEL" | grep -qi "rt"; then
        success "Real-time kernel detected"
    else
        info "Standard kernel (not RT)"
    fi
fi

check "CRI-O Version"
if [ -n "$WORKER_NODE" ]; then
    CRIO_VERSION=$(oc debug node/"$WORKER_NODE" -- chroot /host crio --version 2>/dev/null | head -1)
    info "CRI-O: ${CRIO_VERSION:-unknown}"
fi

check "chrt command availability"
if [ -n "$WORKER_NODE" ]; then
    CHRT_CHECK=$(oc debug node/"$WORKER_NODE" -- chroot /host which chrt 2>/dev/null)
    if [ -n "$CHRT_CHECK" ]; then
        success "chrt available at: $CHRT_CHECK"
    else
        error "chrt command not found on node"
    fi
fi

section "5. LIVE POD TESTS"

check "Creating test namespace: $TEST_NS"
oc create namespace "$TEST_NS" 2>/dev/null
oc label namespace "$TEST_NS" pod-security.kubernetes.io/enforce=privileged --overwrite 2>/dev/null

# Test 1: Exclusive CPU Pool Pod
check "Creating exclusive CPU pool test pod"
cat <<EOF | oc apply -n "$TEST_NS" -f - 2>/dev/null
apiVersion: v1
kind: Pod
metadata:
  name: $TEST_POD_EXCLUSIVE
  labels:
    test: exclusive-cpu
spec:
  containers:
  - name: test
    image: quay.io/testnetworkfunction/k8s-best-practices-debug:latest
    command: ["sleep", "infinity"]
    resources:
      requests:
        cpu: "1"
        memory: "256Mi"
      limits:
        cpu: "1"
        memory: "256Mi"
    securityContext:
      privileged: true
  terminationGracePeriodSeconds: 0
EOF

# Test 2: Shared CPU Pool Pod
check "Creating shared CPU pool test pod"
cat <<EOF | oc apply -n "$TEST_NS" -f - 2>/dev/null
apiVersion: v1
kind: Pod
metadata:
  name: $TEST_POD_SHARED
  labels:
    test: shared-cpu
spec:
  containers:
  - name: test
    image: registry.access.redhat.com/ubi8/ubi-micro:latest
    command: ["sleep", "infinity"]
    resources:
      requests:
        cpu: "100m"
        memory: "64Mi"
      limits:
        cpu: "200m"
        memory: "128Mi"
  terminationGracePeriodSeconds: 0
EOF

log "\nWaiting for pods to be ready (up to 120s)..."
oc wait --for=condition=Ready pod/"$TEST_POD_EXCLUSIVE" -n "$TEST_NS" --timeout=120s 2>/dev/null
EXCLUSIVE_READY=$?
oc wait --for=condition=Ready pod/"$TEST_POD_SHARED" -n "$TEST_NS" --timeout=120s 2>/dev/null
SHARED_READY=$?

log "\n--- Exclusive CPU Pool Pod Analysis ---"
if [ $EXCLUSIVE_READY -eq 0 ]; then
    success "Exclusive CPU pod is running"

    check "QoS Class"
    QOS=$(oc get pod "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" -o jsonpath='{.status.qosClass}' 2>/dev/null)
    if [ "$QOS" == "Guaranteed" ]; then
        success "QoS Class: $QOS (correct for exclusive CPU)"
    else
        error "QoS Class: $QOS (should be Guaranteed)"
    fi

    check "Resource Configuration"
    oc get pod "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" -o jsonpath='CPU Requests: {.spec.containers[0].resources.requests.cpu}, Limits: {.spec.containers[0].resources.limits.cpu}' 2>/dev/null | tee -a "$OUTPUT_FILE"
    log ""

    check "CPU Assignment (cgroup)"
    CPUSET=$(oc exec "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" -- cat /sys/fs/cgroup/cpuset/cpuset.cpus 2>/dev/null || \
             oc exec "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" -- cat /sys/fs/cgroup/cpuset.cpus.effective 2>/dev/null)
    info "Assigned CPUs: ${CPUSET:-unknown}"

    check "Scheduling Policy (PID 1)"
    CHRT_OUTPUT=$(oc exec "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" -- chrt -p 1 2>/dev/null)
    if [ -n "$CHRT_OUTPUT" ]; then
        info "$CHRT_OUTPUT"

        POLICY=$(echo "$CHRT_OUTPUT" | grep "scheduling policy" | awk '{print $NF}')
        PRIORITY=$(echo "$CHRT_OUTPUT" | grep "scheduling priority" | awk '{print $NF}')

        log "\n  Parsed: policy=$POLICY, priority=$PRIORITY"

        # Check against certsuite requirements
        if [ "$PRIORITY" == "0" ]; then
            success "Priority 0 - PASSES exclusive CPU pool check (SCHED_OTHER allowed)"
        elif [ "$PRIORITY" -lt 10 ] && [[ "$POLICY" == "SCHED_FIFO" || "$POLICY" == "SCHED_RR" ]]; then
            success "RT policy with priority < 10 - PASSES exclusive CPU pool check"
        else
            error "Policy=$POLICY, Priority=$PRIORITY - FAILS exclusive CPU pool check"
            info "Expected: priority==0 OR (priority<10 AND policy in [SCHED_RR, SCHED_FIFO])"
        fi
    else
        error "Could not execute chrt command in pod"
    fi

    check "All Processes Scheduling"
    PS_OUTPUT=$(oc exec "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" -- ps -eo pid,comm 2>/dev/null | tail -n +2)
    if [ -n "$PS_OUTPUT" ]; then
        log "  Checking scheduling policy for all processes:"
        while IFS= read -r line; do
            PID=$(echo "$line" | awk '{print $1}')
            COMM=$(echo "$line" | awk '{print $2}')
            if [ -n "$PID" ] && [ "$PID" != "PID" ]; then
                PROC_CHRT=$(oc exec "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" -- chrt -p "$PID" 2>/dev/null)
                PROC_PRIORITY=$(echo "$PROC_CHRT" | grep "priority" | awk '{print $NF}')
                PROC_POLICY=$(echo "$PROC_CHRT" | grep "policy" | awk '{print $NF}')
                info "  PID $PID ($COMM): policy=$PROC_POLICY priority=$PROC_PRIORITY"
            fi
        done <<< "$PS_OUTPUT"
    fi
else
    error "Exclusive CPU pod failed to start"
    oc describe pod "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" 2>/dev/null | tail -20 | tee -a "$OUTPUT_FILE"
fi

log "\n--- Shared CPU Pool Pod Analysis ---"
if [ $SHARED_READY -eq 0 ]; then
    success "Shared CPU pod is running"

    check "QoS Class"
    QOS=$(oc get pod "$TEST_POD_SHARED" -n "$TEST_NS" -o jsonpath='{.status.qosClass}' 2>/dev/null)
    if [ "$QOS" == "Burstable" ]; then
        success "QoS Class: $QOS (correct for shared CPU pool)"
    else
        info "QoS Class: $QOS"
    fi

    check "Resource Configuration"
    oc get pod "$TEST_POD_SHARED" -n "$TEST_NS" -o jsonpath='CPU Requests: {.spec.containers[0].resources.requests.cpu}, Limits: {.spec.containers[0].resources.limits.cpu}' 2>/dev/null | tee -a "$OUTPUT_FILE"
    log ""

    # Note: ubi-micro doesn't have chrt, so we can't check scheduling policy directly
    info "Note: ubi-micro image doesn't have chrt - certsuite uses probe pod for this check"
else
    error "Shared CPU pod failed to start"
    oc describe pod "$TEST_POD_SHARED" -n "$TEST_NS" 2>/dev/null | tail -20 | tee -a "$OUTPUT_FILE"
fi

section "6. CERTSUITE PROBE POD SIMULATION"

check "Testing scheduling policy check via node debug (simulates certsuite probe)"
if [ -n "$WORKER_NODE" ] && [ $EXCLUSIVE_READY -eq 0 ]; then
    # Get the container's PID on the host
    POD_UID=$(oc get pod "$TEST_POD_EXCLUSIVE" -n "$TEST_NS" -o jsonpath='{.metadata.uid}' 2>/dev/null)
    info "Pod UID: $POD_UID"

    # Find container PID via crictl
    CONTAINER_ID=$(oc debug node/"$WORKER_NODE" -- chroot /host crictl ps --name test -q 2>/dev/null | head -1)
    if [ -n "$CONTAINER_ID" ]; then
        info "Container ID: $CONTAINER_ID"

        CONTAINER_PID=$(oc debug node/"$WORKER_NODE" -- chroot /host crictl inspect "$CONTAINER_ID" 2>/dev/null | grep -m1 '"pid":' | awk -F: '{print $2}' | tr -d ' ,')
        if [ -n "$CONTAINER_PID" ]; then
            info "Container PID on host: $CONTAINER_PID"

            # Check scheduling from host perspective
            HOST_CHRT=$(oc debug node/"$WORKER_NODE" -- chroot /host chrt -p "$CONTAINER_PID" 2>/dev/null)
            if [ -n "$HOST_CHRT" ]; then
                info "Host view of container scheduling:"
                echo "$HOST_CHRT" | tee -a "$OUTPUT_FILE"
            fi
        fi
    else
        info "Could not find container via crictl (this is informational only)"
    fi
fi

section "7. MACHINE CONFIG & TUNING"

check "MachineConfigPools"
oc get mcp 2>/dev/null | tee -a "$OUTPUT_FILE"

check "Performance-related MachineConfigs"
oc get mc 2>/dev/null | grep -iE "performance|kubelet|crio" | tee -a "$OUTPUT_FILE"

check "Tuned Profiles"
oc get tuned -A 2>/dev/null | tee -a "$OUTPUT_FILE"

section "8. POTENTIAL ISSUES SUMMARY"

log "Based on the diagnostics above, check for these common issues:\n"

# Check 1: No performance profile
PP_COUNT=$(oc get performanceprofiles --no-headers 2>/dev/null | wc -l)
if [ "$PP_COUNT" -eq 0 ]; then
    error "ISSUE: No PerformanceProfile configured"
    info "  → Exclusive CPU pool tests require a PerformanceProfile"
    info "  → Install Node Tuning Operator and create a PerformanceProfile"
fi

# Check 2: Compact cluster
if [ "$WORKER_COUNT" -eq 0 ]; then
    error "ISSUE: Compact cluster detected (no dedicated workers)"
    info "  → Tests may behave differently on compact clusters"
    info "  → Master nodes may have different CPU configurations"
fi

# Check 3: CPU manager policy
if [ -n "$CPU_MANAGER_STATE" ]; then
    if echo "$CPU_MANAGER_STATE" | grep -q '"policyName":"none"'; then
        error "ISSUE: CPU Manager policy is 'none'"
        info "  → Exclusive CPUs cannot be assigned"
        info "  → Configure static CPU manager policy via PerformanceProfile"
    fi
fi

section "9. CLEANUP"

log "Deleting test namespace: $TEST_NS"
# Cleanup happens via trap

section "REPORT COMPLETE"

log "Full diagnostic report saved to: $OUTPUT_FILE"
log "\nTo compare across clusters:"
log "  1. Run this script on 4.14, 4.16, 4.17 clusters"
log "  2. Save outputs with cluster version in filename"
log "  3. Compare the 'Scheduling Policy' and 'CPU Assignment' sections"
log ""
log "Key things to compare:"
log "  - PerformanceProfile presence and configuration"
log "  - CPU Manager policy (static vs none)"
log "  - Scheduling policy output from chrt"
log "  - QoS class of test pods"
