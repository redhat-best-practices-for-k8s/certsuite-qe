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
# TODO: Change this to the latest OCP versions once they are available.
DEFAULT_VERSIONS="4.14 4.17 4.18 4.19 4.20"

# Operators we depend on for QE tests
# Format: catalog_type:operator_package_name:description
#
# Note: We have version-specific operators due to catalog availability:
#   - cockroachdb-certified: Used for OCP 4.14-4.19 (certified operator)
#   - mongodb-enterprise: Used for OCP 4.20+ (alternative certified operator)
#   - postgresql: Used for OCP 4.14-4.19 (lightweight operator)
#   - prometheus-exporter-operator: Used for OCP 4.20+ (alternative lightweight operator)
# See tests/utils/operatorversions/ for the mapping logic.
REQUIRED_OPERATORS="
certified-operators:cockroachdb-certified:Used for affiliated certification tests (OCP 4.14-4.19)
certified-operators:mongodb-enterprise:Used for affiliated certification tests (OCP 4.20+)
community-operators:grafana-operator:Used for operator tests
community-operators:postgresql:Used as lightweight operator (OCP 4.14-4.19)
community-operators:prometheus-exporter-operator:Used as lightweight operator (OCP 4.20+)
redhat-operators:cluster-logging:Used for cluster-wide operator tests
"

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

# Check if logged into registry
check_registry_login() {
	echo -e "${BLUE}Checking registry authentication...${NC}"

	# Try a simple catalog query to check auth
	if ! opm render "registry.redhat.io/redhat/certified-operator-index:v4.14" 2>&1 | head -1 >/dev/null; then
		echo -e "${YELLOW}Warning: You may need to authenticate to registry.redhat.io${NC}"
		echo ""
		echo "Run: podman login registry.redhat.io"
		echo "  or: docker login registry.redhat.io"
		echo ""
		echo "You'll need a Red Hat account. Get one at: https://access.redhat.com"
		echo ""
	fi
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

# Check a single operator in a catalog
check_operator() {
	local ocp_version="$1"
	local catalog_type="$2"
	local operator_name="$3"
	local description="$4"

	local catalog_base
	catalog_base=$(get_catalog_image "$catalog_type")
	local catalog_image="${catalog_base}:v${ocp_version}"

	echo -e "  ${BLUE}Checking:${NC} $operator_name"
	echo -e "    Catalog: $catalog_type"
	echo -e "    Image:   $catalog_image"

	# Force pull if requested
	if [ "$FRESH_PULL" = "true" ]; then
		force_pull_image "$catalog_image"
	fi

	if [ "$VERBOSE" = "true" ]; then
		echo -e "    ${YELLOW}Running: opm render $catalog_image | grep \"name\":\"$operator_name\"${NC}"
	fi

	# Try to find the operator in the catalog
	# opm render pulls the catalog image and outputs all packages/channels/bundles as JSON
	# We grep for the operator package name (handles both "name": "X" and "name":"X" formats)
	# Note: Use [ ]* instead of \s* for portable regex (works on both GNU and BSD grep)
	echo -n "    Status:  "
	# Note: 'timeout' may not exist on macOS, so we run without it
	if opm render "$catalog_image" 2>/dev/null | grep -qE "\"(name|package)\":[ ]*\"$operator_name\""; then
		echo -e "${GREEN}✅ Found${NC}"
		return 0
	else
		echo -e "${RED}❌ MISSING${NC}"
		return 1
	fi
}

# List all operators in a catalog (for debugging)
list_operators() {
	local ocp_version="$1"
	local catalog_type="$2"
	local search_term="${3:-}"

	local catalog_base
	catalog_base=$(get_catalog_image "$catalog_type")
	local catalog_image="${catalog_base}:v${ocp_version}"

	echo -e "${BLUE}Listing operators in $catalog_type for OCP $ocp_version...${NC}"

	if [ -n "$search_term" ]; then
		echo "Filtering for: $search_term"
		opm render "$catalog_image" 2>/dev/null |
			jq -r 'select(.schema == "olm.package") | .name' |
			grep -i "$search_term" |
			sort -u
	else
		opm render "$catalog_image" 2>/dev/null |
			jq -r 'select(.schema == "olm.package") | .name' |
			sort -u
	fi
}

# Get available channels and versions for an operator
get_operator_info() {
	local ocp_version="$1"
	local catalog_type="$2"
	local operator_name="$3"

	local catalog_base
	catalog_base=$(get_catalog_image "$catalog_type")
	local catalog_image="${catalog_base}:v${ocp_version}"

	echo -e "${BLUE}Getting info for $operator_name in $catalog_type (OCP $ocp_version)...${NC}"

	# Get channels
	echo "Channels:"
	opm render "$catalog_image" 2>/dev/null |
		jq -r "select(.schema == \"olm.channel\" and .package == \"$operator_name\") | \"  - \" + .name" |
		sort -u

	# Get bundle versions
	echo "Versions:"
	opm render "$catalog_image" 2>/dev/null |
		jq -r "select(.schema == \"olm.bundle\" and .package == \"$operator_name\") | \"  - \" + .name" |
		sort -V |
		tail -10
}

# Main check function
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

	local missing_count=0
	local total_count=0

	for ocp_version in $versions; do
		echo ""
		echo -e "${YELLOW}=== OCP $ocp_version ===${NC}"

		echo "$REQUIRED_OPERATORS" | while IFS= read -r operator_entry; do
			# Skip empty lines
			[ -z "$operator_entry" ] && continue

			catalog_type=$(echo "$operator_entry" | cut -d: -f1)
			operator_name=$(echo "$operator_entry" | cut -d: -f2)
			description=$(echo "$operator_entry" | cut -d: -f3)

			total_count=$((total_count + 1))

			if ! check_operator "$ocp_version" "$catalog_type" "$operator_name" "$description"; then
				missing_count=$((missing_count + 1))
				echo -e "    ${YELLOW}↳ $description${NC}"
			fi
		done
	done

	echo ""
	echo -e "${BLUE}========================================${NC}"
	echo -e "${BLUE}  Summary${NC}"
	echo -e "${BLUE}========================================${NC}"
	echo ""

	# Re-count since subshell doesn't persist variables
	missing_count=0
	total_count=0
	for ocp_version in $versions; do
		echo "$REQUIRED_OPERATORS" | while IFS= read -r operator_entry; do
			[ -z "$operator_entry" ] && continue
			total_count=$((total_count + 1))
		done
	done

	echo "Completed catalog check for OCP versions: $versions"
	echo -e "${GREEN}Check complete. Review output above for any missing operators.${NC}"
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
	echo "  The catalog images contain all operators available for a specific OCP version."
	echo ""
	echo "Note: By default, cached catalog images may be used. Use -f/--fresh to ensure"
	echo "      you're checking against the latest catalog content."
}

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
		list_operators "${4:-4.20}" "$2" "${3:-}"
		exit 0
		;;
	-i | --info)
		check_prerequisites
		get_operator_info "${4:-4.20}" "$2" "$3"
		exit 0
		;;
	*)
		POSITIONAL_ARGS+=("$1")
		shift
		;;
	esac
done

# Run main check with remaining positional arguments
check_prerequisites
check_registry_login
main_check "${POSITIONAL_ARGS[@]}"
