#!/bin/bash
#
# Check Operator Catalogs for Required Operators
#
# This script checks Red Hat operator catalogs for the operators that our QE tests depend on.
# It helps diagnose issues where operators are no longer available in certain OCP versions.
#
# Prerequisites:
#   - podman or docker (for registry authentication)
#   - opm (Operator Package Manager) - https://github.com/operator-framework/operator-registry
#   - jq (for JSON parsing) - https://jqlang.github.io/jq/
#   - Logged into registry.redhat.io (podman login registry.redhat.io)
#
# Usage:
#   ./check-operator-catalogs.sh                    # Check all default versions
#   ./check-operator-catalogs.sh 4.20               # Check specific version
#   ./check-operator-catalogs.sh 4.19 4.20          # Check multiple versions
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Verbose mode (set with -v flag)
VERBOSE="${VERBOSE:-false}"

# Fresh pull mode (set with -f/--fresh flag) - always re-pull catalog images
FRESH_PULL="${FRESH_PULL:-false}"

# Default OCP versions to check
DEFAULT_VERSIONS="4.14 4.16 4.17 4.18 4.19 4.20"

# Operators we depend on for QE tests, split by OCP version.
# Format: catalog_type:operator_package_name:description
#
# The certified operator differs by catalog availability:
#   - cockroachdb-certified: Used for OCP 4.14-4.19 (certified operator)
#   - mongodb-enterprise: Used for OCP 4.20+ (alternative certified operator)
# See tests/utils/operatorversions/ for the authoritative mapping.
COMMON_OPERATORS="
community-operators:grafana-operator:Used for operator tests
community-operators:prometheus-exporter-operator:Used as lightweight operator (all OCP versions)
redhat-operators:cluster-logging:Used for cluster-wide operator tests
"

PRE420_OPERATORS="
certified-operators:cockroachdb-certified:Used for affiliated certification tests (OCP 4.14-4.19)
"

POST420_OPERATORS="
certified-operators:mongodb-enterprise:Used for affiliated certification tests (OCP 4.20+)
"

# Operators required for the given OCP version, in the
# catalog_type:operator_package_name:description format.
# Keep the 4.20 cutoff in sync with IsVersion420OrLater in
# tests/utils/operatorversions/operatorversions.go. Versions that do not parse
# fall back to the pre-4.20 set, as that helper does.
required_operators_for_version() {
	local ocp_version="$1"
	local major minor

	major="${ocp_version%%.*}"
	minor="${ocp_version#*.}"
	minor="${minor%%.*}"

	printf '%s\n' "$COMMON_OPERATORS"

	if [[ "$major" =~ ^[0-9]+$ && "$minor" =~ ^[0-9]+$ ]]; then
		if [ "$major" -gt 4 ] || { [ "$major" -eq 4 ] && [ "$minor" -ge 20 ]; }; then
			printf '%s\n' "$POST420_OPERATORS"

			return 0
		fi
	fi

	printf '%s\n' "$PRE420_OPERATORS"
}

# Get catalog image URL for a catalog type
get_catalog_image() {
	local catalog_type="$1"
	case "$catalog_type" in
	certified-operators)
		echo "registry.redhat.io/redhat/certified-operator-index"
		;;
	community-operators)
		echo "registry.redhat.io/redhat/community-operator-index"
		;;
	redhat-operators)
		echo "registry.redhat.io/redhat/redhat-operator-index"
		;;
	*)
		echo ""
		;;
	esac
}

# Check for opm updates
check_opm_update() {
	local current_version="$1"

	echo -e "${BLUE}Checking for opm updates...${NC}"

	# Fetch latest release from GitHub API
	local latest_version
	latest_version=$(curl -sL "https://api.github.com/repos/operator-framework/operator-registry/releases/latest" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

	if [ -z "$latest_version" ]; then
		echo -e "  ${YELLOW}⚠${NC} Could not check for updates (GitHub API unavailable)"
		return
	fi

	echo -e "  Current version: ${current_version}"
	echo -e "  Latest version:  ${latest_version}"

	# Compare versions (strip 'v' prefix for comparison)
	local current_num="${current_version#v}"
	local latest_num="${latest_version#v}"

	if [ "$current_num" = "$latest_num" ]; then
		echo -e "  ${GREEN}✓${NC} opm is up to date!"
		return
	fi

	# Simple version comparison (works for semver)
	if [ "$(printf '%s\n' "$latest_num" "$current_num" | sort -V | tail -n1)" = "$latest_num" ] && [ "$latest_num" != "$current_num" ]; then
		echo ""
		echo -e "  ${YELLOW}⚠ A newer version of opm is available!${NC}"
		echo ""
		echo -e "  Release notes: https://github.com/operator-framework/operator-registry/releases/tag/${latest_version}"
		echo ""

		# Detect OS and architecture for download URL
		local os_type
		local arch_type
		case "$(uname -s)" in
		Linux*) os_type="linux" ;;
		Darwin*) os_type="darwin" ;;
		*) os_type="unknown" ;;
		esac
		case "$(uname -m)" in
		x86_64) arch_type="amd64" ;;
		aarch64) arch_type="arm64" ;;
		arm64) arch_type="arm64" ;;
		*) arch_type="unknown" ;;
		esac

		if [ "$os_type" != "unknown" ] && [ "$arch_type" != "unknown" ]; then
			local download_url="https://github.com/operator-framework/operator-registry/releases/download/${latest_version}/${os_type}-${arch_type}-opm"

			echo -n "  Would you like to update opm now? [y/N]: "
			read -r response

			if [[ "$response" =~ ^[Yy]$ ]]; then
				echo ""
				echo -e "  ${BLUE}Downloading opm ${latest_version}...${NC}"
				echo "  URL: $download_url"

				local temp_file="/tmp/opm-${latest_version}"
				if curl -sL "$download_url" -o "$temp_file"; then
					chmod +x "$temp_file"

					local opm_path
					opm_path=$(command -v opm)

					echo ""
					echo -e "  ${BLUE}Installing to ${opm_path}...${NC}"
					echo "  Running: sudo mv $temp_file $opm_path"

					if sudo mv "$temp_file" "$opm_path"; then
						echo -e "  ${GREEN}✓ opm updated successfully to ${latest_version}!${NC}"
						# Verify
						echo -e "  Verifying: $(opm version 2>&1 | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+')"
					else
						echo -e "  ${RED}Failed to install. You can manually run:${NC}"
						echo "    sudo mv $temp_file $opm_path"
					fi
					echo ""
				else
					echo -e "  ${RED}Failed to download opm${NC}"
				fi
			else
				echo -e "  Skipping update."
			fi
		fi
	fi
	echo ""
}

# Ensure containers policy.json exists (required by opm render)
ensure_policy_json() {
	local policy_path="$HOME/.config/containers/policy.json"

	if [ ! -f "$policy_path" ] && [ ! -f "/etc/containers/policy.json" ]; then
		echo -e "  ${YELLOW}Creating containers policy.json (required by opm)...${NC}"
		mkdir -p "$HOME/.config/containers"
		cat >"$policy_path" <<'POLICY_EOF'
{
    "default": [
        {
            "type": "insecureAcceptAnything"
        }
    ],
    "transports": {
        "docker-daemon": {
            "": [
                {
                    "type": "insecureAcceptAnything"
                }
            ]
        }
    }
}
POLICY_EOF
		echo -e "  ${GREEN}✓${NC} Created $policy_path"
	fi
}

# Check if opm is installed
check_prerequisites() {
	echo -e "${BLUE}Checking prerequisites...${NC}"

	# Ensure policy.json exists for opm
	ensure_policy_json

	if ! command -v opm &>/dev/null; then
		echo -e "${RED}Error: opm (Operator Package Manager) is not installed.${NC}"
		echo ""
		echo "Install opm from: https://github.com/operator-framework/operator-registry/releases"
		echo ""

		# Detect OS and architecture
		local os_type arch_type
		case "$(uname -s)" in
		Linux*) os_type="linux" ;;
		Darwin*) os_type="darwin" ;;
		*) os_type="linux" ;;
		esac
		case "$(uname -m)" in
		x86_64) arch_type="amd64" ;;
		aarch64 | arm64) arch_type="arm64" ;;
		*) arch_type="amd64" ;;
		esac

		# Get latest version
		local latest_version
		latest_version=$(curl -sL "https://api.github.com/repos/operator-framework/operator-registry/releases/latest" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
		latest_version="${latest_version:-v1.61.0}"

		echo "Quick install (${os_type} ${arch_type}):"
		echo "  curl -sLO https://github.com/operator-framework/operator-registry/releases/download/${latest_version}/${os_type}-${arch_type}-opm"
		echo "  chmod +x ${os_type}-${arch_type}-opm"
		echo "  sudo mv ${os_type}-${arch_type}-opm /usr/local/bin/opm"
		exit 1
	fi

	local current_version
	current_version=$(opm version 2>&1 | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+')

	echo -e "  ${GREEN}✓${NC} opm found at: $(command -v opm)"
	echo -e "  ${GREEN}✓${NC} opm version: ${current_version}"

	# Check for jq
	if ! command -v jq &>/dev/null; then
		echo -e "${RED}Error: jq is not installed.${NC}"
		echo ""
		echo "Install jq from: https://jqlang.github.io/jq/download/"
		echo ""
		echo "Quick install:"
		echo "  macOS: brew install jq"
		echo "  Ubuntu/Debian: sudo apt-get install jq"
		echo "  Fedora/RHEL: sudo dnf install jq"
		exit 1
	fi
	echo -e "  ${GREEN}✓${NC} jq found at: $(command -v jq)"
	echo ""

	# Check for updates
	check_opm_update "$current_version"

	echo -e "${BLUE}How this script works:${NC}"
	echo "  1. Uses 'opm render <catalog-image>' to pull and render catalog contents"
	echo "  2. The catalog images are from registry.redhat.io:"
	echo "     - certified-operator-index:v4.XX"
	echo "     - community-operator-index:v4.XX"
	echo "     - redhat-operator-index:v4.XX"
	echo "  3. Searches the rendered JSON output for each operator package name"
	echo ""
}

# Force pull a catalog image to ensure we have the latest
force_pull_image() {
	local catalog_image="$1"

	echo -e "    ${YELLOW}Pulling fresh image...${NC}"

	# Try podman first, then docker
	if command -v podman &>/dev/null; then
		podman pull --quiet "$catalog_image" 2>/dev/null || true
	elif command -v docker &>/dev/null; then
		docker pull --quiet "$catalog_image" 2>/dev/null || true
	fi
}

# Rendered catalogs are shared by every check in a run, so a catalog image is
# rendered at most once per OCP version. The directory is created on first use
# and removed by the EXIT trap.
RENDER_DIR=""

# Result of the last check_operator call: present, missing or render-failed.
CHECK_RESULT=""

ensure_render_dir() {
	if [ -z "$RENDER_DIR" ]; then
		RENDER_DIR=$(mktemp -d)
		trap 'rm -rf "$RENDER_DIR"' EXIT
	fi
}

# Path of the rendered catalog for an OCP version, and of its status file. The
# status file holds ok, empty or failed and doubles as the "already rendered"
# marker so a catalog is never rendered twice in one run.
render_file_path() {
	local ocp_version="$1"
	local catalog_type="$2"

	printf '%s/%s-v%s.json' "$RENDER_DIR" "$catalog_type" "$ocp_version"
}

render_status_path() {
	local ocp_version="$1"
	local catalog_type="$2"

	printf '%s/%s-v%s.status' "$RENDER_DIR" "$catalog_type" "$ocp_version"
}

catalog_image_for() {
	local catalog_type="$1"
	local ocp_version="$2"

	local catalog_base
	catalog_base=$(get_catalog_image "$catalog_type")
	if [ -z "$catalog_base" ]; then
		return 1
	fi

	printf '%s:v%s\n' "$catalog_base" "$ocp_version"
}

# Render a catalog image once, recording ok/empty/failed in its status file.
render_catalog() {
	local ocp_version="$1"
	local catalog_type="$2"

	ensure_render_dir

	local status_file
	status_file=$(render_status_path "$ocp_version" "$catalog_type")

	if [ -f "$status_file" ]; then
		return 0
	fi

	local catalog_image
	if ! catalog_image=$(catalog_image_for "$catalog_type" "$ocp_version"); then
		echo -e "  ${YELLOW}⚠ Unknown catalog type: $catalog_type${NC}"
		echo failed >"$status_file"

		return 0
	fi

	echo -e "  ${BLUE}Rendering catalog:${NC} $catalog_type"
	echo -e "    Image:   $catalog_image"

	if [ "$FRESH_PULL" = "true" ]; then
		force_pull_image "$catalog_image"
	fi

	if [ "$VERBOSE" = "true" ]; then
		echo -e "    ${YELLOW}Running: opm render $catalog_image${NC}"
	fi

	local render_file render_err
	render_file=$(render_file_path "$ocp_version" "$catalog_type")
	render_err="${render_file}.err"

	# opm render pulls the catalog image and outputs all packages/channels/bundles
	# as JSON. A failing render is an infrastructure problem (registry auth, opm
	# bugs), not a missing operator, so keep its stderr and status around instead
	# of collapsing it into "missing".
	if opm render "$catalog_image" >"$render_file" 2>"$render_err"; then
		if [ -s "$render_file" ]; then
			echo ok >"$status_file"
		else
			echo empty >"$status_file"
			echo -e "    ${YELLOW}⚠ opm render produced no content${NC}"
		fi
	else
		echo failed >"$status_file"
		echo -e "    ${YELLOW}⚠ opm render failed${NC}"
		sed -n '1,5p' "$render_err" | sed 's/^/      /'

		# The usual cause is a missing registry login, so point at the fix here
		# instead of paying for a separate preflight render of the catalog.
		echo "      If this is an authentication error, run: podman login registry.redhat.io"
		echo "      You'll need a Red Hat account. Get one at: https://access.redhat.com"
	fi

	return 0
}

# Check whether an operator package exists in an already rendered catalog.
operator_in_catalog() {
	local operator_name="$1"
	local render_file="$2"

	# Match only olm.package schema entries: grep would produce false positives
	# by matching relatedImages.
	jq -e "select(.schema == \"olm.package\" and .name == \"$operator_name\")" "$render_file" >/dev/null 2>&1
}

# Check a single operator in a catalog. Sets CHECK_RESULT to present, missing or
# render-failed; it reports through a global rather than stdout so callers can
# count results without a subshell swallowing the increments.
check_operator() {
	local ocp_version="$1"
	local catalog_type="$2"
	local operator_name="$3"
	local description="$4"

	local catalog_image
	catalog_image=$(catalog_image_for "$catalog_type" "$ocp_version") || catalog_image="unknown"

	echo -e "  ${BLUE}Checking:${NC} $operator_name"
	echo -e "    Catalog: $catalog_type"
	echo -e "    Image:   $catalog_image"

	local status
	status=$(cat "$(render_status_path "$ocp_version" "$catalog_type")" 2>/dev/null || echo failed)

	case "$status" in
	ok)
		echo -n "    Status:  "
		if operator_in_catalog "$operator_name" "$(render_file_path "$ocp_version" "$catalog_type")"; then
			echo -e "${GREEN}✅ Found${NC}"
			CHECK_RESULT="present"
		else
			echo -e "${RED}❌ MISSING${NC}"
			CHECK_RESULT="missing"
		fi
		;;
	empty)
		echo -e "    Status:  ${YELLOW}⚠ catalog rendered empty${NC}"
		CHECK_RESULT="render-failed"
		;;
	*)
		echo -e "    Status:  ${YELLOW}⚠ catalog failed to render${NC}"
		CHECK_RESULT="render-failed"
		;;
	esac

	echo -e "    ${YELLOW}↳ $description${NC}"
}

# Distinct catalog types referenced by any version's operator list, in a stable
# order. Catalog images are shared across versions, so the union is checked.
required_catalog_types() {
	printf '%s\n%s\n%s\n' "$COMMON_OPERATORS" "$PRE420_OPERATORS" "$POST420_OPERATORS" |
		while IFS= read -r operator_entry; do
			[ -z "$operator_entry" ] && continue
			printf '%s\n' "${operator_entry%%:*}"
		done |
		sort -u
}

# List all operators in a catalog (for debugging)
list_operators() {
	local ocp_version="$1"
	local catalog_type="$2"
	local search_term="${3:-}"

	echo -e "${BLUE}Listing operators in $catalog_type for OCP $ocp_version...${NC}"

	render_catalog "$ocp_version" "$catalog_type"

	local render_file status
	render_file=$(render_file_path "$ocp_version" "$catalog_type")
	status=$(cat "$(render_status_path "$ocp_version" "$catalog_type")" 2>/dev/null || echo failed)

	if [ "$status" != "ok" ]; then
		echo -e "${RED}Error: catalog $catalog_type for OCP $ocp_version could not be rendered (status: $status)${NC}"

		return 2
	fi

	if [ -n "$search_term" ]; then
		echo "Filtering for: $search_term"
		jq -r 'select(.schema == "olm.package") | .name' "$render_file" |
			grep -i "$search_term" |
			sort -u
	else
		jq -r 'select(.schema == "olm.package") | .name' "$render_file" |
			sort -u
	fi
}

# Get available channels and versions for an operator
get_operator_info() {
	local ocp_version="$1"
	local catalog_type="$2"
	local operator_name="$3"

	echo -e "${BLUE}Getting info for $operator_name in $catalog_type (OCP $ocp_version)...${NC}"

	render_catalog "$ocp_version" "$catalog_type"

	local render_file status
	render_file=$(render_file_path "$ocp_version" "$catalog_type")
	status=$(cat "$(render_status_path "$ocp_version" "$catalog_type")" 2>/dev/null || echo failed)

	if [ "$status" != "ok" ]; then
		echo -e "${RED}Error: catalog $catalog_type for OCP $ocp_version could not be rendered (status: $status)${NC}"

		return 2
	fi

	# Get channels
	echo "Channels:"
	jq -r "select(.schema == \"olm.channel\" and .package == \"$operator_name\") | \"  - \" + .name" "$render_file" |
		sort -u

	# Get bundle versions
	echo "Versions:"
	jq -r "select(.schema == \"olm.bundle\" and .package == \"$operator_name\") | \"  - \" + .name" "$render_file" |
		sort -V |
		tail -10
}

# Main check function. Returns 0 when every operator is present and every catalog
# rendered, 1 when at least one operator is missing, 2 when at least one catalog
# could not be rendered (the check is incomplete).
main_check() {
	local versions="$*"

	if [ -z "$versions" ]; then
		versions="$DEFAULT_VERSIONS"
	fi

	echo -e "${BLUE}========================================${NC}"
	echo -e "${BLUE}  Operator Catalog Check${NC}"
	echo -e "${BLUE}========================================${NC}"
	echo ""
	echo "Checking OCP versions: $versions"
	echo "Date: $(date)"
	if [ "$FRESH_PULL" = "true" ]; then
		echo -e "Mode: ${GREEN}Fresh pull enabled${NC} (re-pulling all catalog images)"
	else
		echo -e "Mode: Using cached images (use -f/--fresh to force re-pull)"
	fi
	echo ""

	local total_count=0
	local present_count=0
	local missing_count=0
	local render_failure_count=0

	for ocp_version in $versions; do
		echo ""
		echo -e "${YELLOW}=== OCP $ocp_version ===${NC}"

		# One render per distinct catalog, shared by every operator in it.
		while IFS= read -r catalog_type; do
			[ -z "$catalog_type" ] && continue
			render_catalog "$ocp_version" "$catalog_type"
		done < <(required_catalog_types)

		echo ""

		while IFS= read -r operator_entry; do
			[ -z "$operator_entry" ] && continue

			local catalog_type operator_name description
			catalog_type="${operator_entry%%:*}"
			operator_name="${operator_entry#*:}"
			operator_name="${operator_name%%:*}"
			description="${operator_entry#*:}"
			description="${description#*:}"

			total_count=$((total_count + 1))
			check_operator "$ocp_version" "$catalog_type" "$operator_name" "$description"

			case "$CHECK_RESULT" in
			present) present_count=$((present_count + 1)) ;;
			missing) missing_count=$((missing_count + 1)) ;;
			*) render_failure_count=$((render_failure_count + 1)) ;;
			esac
		done < <(required_operators_for_version "$ocp_version")
	done

	echo ""
	echo -e "${BLUE}========================================${NC}"
	echo -e "${BLUE}  Summary${NC}"
	echo -e "${BLUE}========================================${NC}"
	echo ""
	echo "OCP versions checked:   $versions"
	echo "Operator entries checked: $total_count"
	echo -e "  ${GREEN}present:${NC}          $present_count"
	echo -e "  ${RED}missing:${NC}          $missing_count"
	echo -e "  ${YELLOW}render failures:${NC}  $render_failure_count"
	echo ""

	if [ "$render_failure_count" -gt 0 ]; then
		echo -e "${YELLOW}⚠ $render_failure_count check(s) could not complete: a catalog failed to render.${NC}"
		echo "  Render errors indicate infrastructure problems (registry auth, opm bugs),"
		echo "  not missing operators."
		echo ""
	fi

	if [ "$missing_count" -gt 0 ]; then
		echo -e "${RED}❌ $missing_count of $total_count operator entries are missing from their catalog.${NC}"

		# A confirmed missing operator is actionable on its own, so it wins over
		# the "incomplete check" exit code.
		return 1
	fi

	if [ "$render_failure_count" -gt 0 ]; then
		return 2
	fi

	echo -e "${GREEN}✅ All $total_count operator entries are present.${NC}"

	return 0
}

# Show usage
usage() {
	echo "Usage: $0 [options] [ocp_versions...]"
	echo ""
	echo "Options:"
	echo "  -f, --fresh                   Force re-pull catalog images (don't use cache)"
	echo "  -v, --verbose                 Show detailed opm commands being run"
	echo "  -l, --list CATALOG [SEARCH]   List all operators in a catalog"
	echo "  -i, --info CATALOG OPERATOR   Get info about a specific operator"
	echo "  -h, --help                    Show this help message"
	echo ""
	echo "Exit codes:"
	echo "  0  every required operator is present and every catalog rendered"
	echo "  1  at least one required operator is missing from its catalog"
	echo "  2  at least one catalog could not be rendered, so the check is incomplete"
	echo "     (takes effect only when nothing is confirmed missing)"
	echo "     -l/--list and -i/--info use 0 on success and 2 when the catalog cannot be rendered"
	echo ""
	echo "Examples:"
	echo "  $0                           # Check all default OCP versions"
	echo "  $0 4.20                      # Check only OCP 4.20"
	echo "  $0 -f 4.20                   # Check OCP 4.20 with fresh image pulls"
	echo "  $0 -v 4.20                   # Check OCP 4.20 with verbose output"
	echo "  $0 4.19 4.20                 # Check OCP 4.19 and 4.20"
	echo "  $0 -l certified-operators cockroach  # List certified operators matching 'cockroach'"
	echo "  $0 -i certified-operators cockroachdb-certified 4.19  # Get info for cockroachdb-certified"
	echo ""
	echo "Catalogs: certified-operators, community-operators, redhat-operators"
	echo ""
	echo "How it works:"
	echo "  This script uses 'opm render' to pull OLM catalog images from registry.redhat.io"
	echo "  and searches for specific operator package names in the rendered JSON output."
	echo "  Each catalog image is rendered once per OCP version and shared by every"
	echo "  operator check for that version."
	echo ""
	echo "Note: By default, cached catalog images may be used. Use -f/--fresh to ensure"
	echo "      you're checking against the latest catalog content."
}

# Allow tests (check-operator-catalogs-test.sh) to source the functions above
# without running the check.
if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then
	return 0
fi

# Parse arguments
POSITIONAL_ARGS=()
while [[ $# -gt 0 ]]; do
	case "$1" in
	-h | --help)
		usage
		exit 0
		;;
	-f | --fresh)
		FRESH_PULL="true"
		shift
		;;
	-v | --verbose)
		VERBOSE="true"
		shift
		;;
	-l | --list)
		check_prerequisites
		EXIT_CODE=0
		list_operators "${4:-4.20}" "$2" "${3:-}" || EXIT_CODE=$?
		exit "$EXIT_CODE"
		;;
	-i | --info)
		check_prerequisites
		EXIT_CODE=0
		get_operator_info "${4:-4.20}" "$2" "$3" || EXIT_CODE=$?
		exit "$EXIT_CODE"
		;;
	*)
		POSITIONAL_ARGS+=("$1")
		shift
		;;
	esac
done

# Run main check with remaining positional arguments
check_prerequisites
EXIT_CODE=0
main_check "${POSITIONAL_ARGS[@]}" || EXIT_CODE=$?
exit "$EXIT_CODE"
