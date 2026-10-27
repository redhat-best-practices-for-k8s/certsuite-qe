#!/bin/bash
#
# Tests for the version-aware operator selection in check-operator-catalogs.sh.
# Runs without registry access: the catalogs are fixtures in a temp directory.
#
# Usage:
#   ./check-operator-catalogs-test.sh
#
# Exit codes:
#   0  every assertion passed
#   1  at least one assertion failed

set -uo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

# shellcheck source=./check-operator-catalogs.sh
# shellcheck disable=SC1091
source "$SCRIPT_DIR/check-operator-catalogs.sh"

# The sourced script enables errexit; assertions must not abort on a failure.
set +e

FAILURES=0
CHECK_RESULT=""

pass() {
	printf 'PASS: %s\n' "$1"
}

fail() {
	printf 'FAIL: %s\n' "$1" >&2
	FAILURES=$((FAILURES + 1))
}

assert_line() {
	local output="$1"
	local expected="$2"
	local description="$3"

	if printf '%s\n' "$output" | grep -Fxq -- "$expected"; then
		pass "$description"
	else
		fail "$description (missing line: $expected)"
	fi
}

refute_line() {
	local output="$1"
	local unexpected="$2"
	local description="$3"

	if printf '%s\n' "$output" | grep -Fxq -- "$unexpected"; then
		fail "$description (unexpected line: $unexpected)"
	else
		pass "$description"
	fi
}

COCKROACH="certified-operators:cockroachdb-certified:Used for affiliated certification tests (OCP 4.14-4.19)"
MONGODB="certified-operators:mongodb-enterprise:Used for affiliated certification tests (OCP 4.20+)"
CLUSTER_LOGGING="redhat-operators:cluster-logging:Used for cluster-wide operator tests"

# Version-aware selection, matching tests/utils/operatorversions.
operators_419=$(required_operators_for_version "4.19")
assert_line "$operators_419" "$COCKROACH" "4.19 selects cockroachdb-certified"
refute_line "$operators_419" "$MONGODB" "4.19 excludes mongodb-enterprise"
assert_line "$operators_419" "$CLUSTER_LOGGING" "4.19 includes the common operators"

operators_420=$(required_operators_for_version "4.20")
assert_line "$operators_420" "$MONGODB" "4.20 selects mongodb-enterprise"
refute_line "$operators_420" "$COCKROACH" "4.20 excludes cockroachdb-certified"
assert_line "$operators_420" "$CLUSTER_LOGGING" "4.20 includes the common operators"

operators_421=$(required_operators_for_version "4.21")
assert_line "$operators_421" "$MONGODB" "4.21 selects mongodb-enterprise"

operators_nightly=$(required_operators_for_version "4.20.0-0.nightly-2025-01-01")
assert_line "$operators_nightly" "$MONGODB" "a full 4.20 version string selects mongodb-enterprise"

operators_50=$(required_operators_for_version "5.0")
assert_line "$operators_50" "$MONGODB" "a major version above 4 selects mongodb-enterprise"

operators_partial=$(required_operators_for_version "4.x")
assert_line "$operators_partial" "$COCKROACH" "a non-numeric minor version falls back to the pre-4.20 set"

operators_malformed=$(required_operators_for_version "not-a-version")
assert_line "$operators_malformed" "$COCKROACH" "malformed versions fall back to the pre-4.20 set"
refute_line "$operators_malformed" "$MONGODB" "malformed versions do not select mongodb-enterprise"

# Catalog images are version-independent: the union must be deduplicated.
expected_types="certified-operators
community-operators
redhat-operators"
catalog_types=$(required_catalog_types)
if [ "$catalog_types" = "$expected_types" ]; then
	pass "required_catalog_types returns the deduplicated union"
else
	fail "required_catalog_types returns the deduplicated union (got: $catalog_types)"
fi

# FOUND/MISSING detection against a fixture catalog.
RENDER_DIR=$(mktemp -d)
trap 'rm -rf "$RENDER_DIR"' EXIT

cat >"$RENDER_DIR/certified-operators-v4.19.json" <<'EOF'
{"schema":"olm.package","name":"cockroachdb-certified"}
{"schema":"olm.package","name":"redis-operator"}
EOF
printf 'ok\n' >"$RENDER_DIR/certified-operators-v4.19.status"

check_operator "4.19" "certified-operators" "cockroachdb-certified" "fixture" >/dev/null
if [ "$CHECK_RESULT" = "present" ]; then
	pass "a package in the rendered catalog reports present"
else
	fail "a package in the rendered catalog reports present (got: $CHECK_RESULT)"
fi

check_operator "4.19" "certified-operators" "mongodb-enterprise" "fixture" >/dev/null
if [ "$CHECK_RESULT" = "missing" ]; then
	pass "a package absent from the rendered catalog reports missing"
else
	fail "a package absent from the rendered catalog reports missing (got: $CHECK_RESULT)"
fi

printf 'failed\n' >"$RENDER_DIR/certified-operators-v4.19.status"
check_operator "4.19" "certified-operators" "cockroachdb-certified" "fixture" >/dev/null
if [ "$CHECK_RESULT" = "render-failed" ]; then
	pass "a failed render reports render-failed instead of missing"
else
	fail "a failed render reports render-failed instead of missing (got: $CHECK_RESULT)"
fi

# Sourcing must not break direct execution.
if "$SCRIPT_DIR/check-operator-catalogs.sh" --help >/dev/null 2>&1; then
	pass "direct execution still parses arguments"
else
	fail "direct execution still parses arguments"
fi

echo ""

if [ "$FAILURES" -gt 0 ]; then
	printf '%d assertion(s) failed\n' "$FAILURES" >&2

	exit 1
fi

echo "All assertions passed."
