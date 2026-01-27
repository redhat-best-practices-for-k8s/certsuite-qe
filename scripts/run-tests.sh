#!/usr/bin/env bash

GOPATH="${GOPATH:-~/go}"
export PATH=$PATH:$GOPATH/bin
EXCLUDED_FOLDERS=''
GINKGO_SEED_FLAG=''
ALL_TESTS_FOLDERS="$(ls -d ./tests/*/)"

if [[ "${GINKGO_SEED_NUMBER}" != "" ]]; then
	echo "Using ginkgo seed number: ${GINKGO_SEED_NUMBER}"
	GINKGO_SEED_FLAG="--seed=${GINKGO_SEED_NUMBER}"
fi

# Check if the user has specified to run in parallel
PFLAG=""
if [[ ${ENABLE_PARALLEL} == "true" ]]; then
	echo "Running tests in parallel"
	PFLAG="--procs=16"
fi

# Allow for flake retries
FFLAG=""
if [[ ${ENABLE_FLAKY_RETRY} == "true" ]]; then
	echo "Retrying flaky tests"
	FFLAG="--flake-attempts=2"
fi

function run_tests {
	case $1 in
	all)
		echo "#### Run all tests ####"
		all_default_suites=""
		for folder in ${ALL_TESTS_FOLDERS}; do
			for excluded_folder in ${EXCLUDED_FOLDERS}; do
				if [[ $folder == *"${excluded_folder}"* ]]; then
					folder=''
				fi
			done
			if [ -n "$folder" ]; then
				all_default_suites+=" $folder"
			fi
		done
		# strip any spaces from the all_default_suites variable
		all_default_suites=$(echo "$all_default_suites" | xargs)
		# check if the all_default_suites variable is empty
		if [ -z "$all_default_suites" ]; then
			echo "No tests found"
			exit 1
		fi

		echo "Running tests in the following folders: ${all_default_suites}"
		# shellcheck disable=SC2086
		ginkgo -timeout=24h -v ${PFLAG} ${FFLAG} --keep-going "${GINKGO_SEED_FLAG}" --show-node-events --require-suite -r $all_default_suites
		;;
	features)
		if [ -z "$FEATURES" ]; then
			echo "FEATURES env var is empty. Please export FEATURES"
			exit 1
		fi
		echo "#### Run feature tests: ${FEATURES} ####"

		LABEL_FILTER=""
		OCP_FILTER=""
		for feature in ${FEATURES}; do
			# Check for -ocp or -k8s suffix to filter by cluster type
			# Examples:
			#   platformalteration1-ocp -> only run ocp-required tests from platformalteration1
			#   accesscontrol1-k8s -> only run non-ocp-required tests from accesscontrol1
			#   accesscontrol1 -> run all tests from accesscontrol1 (no filter)
			if [[ $feature =~ ^(.+)-ocp$ ]]; then
				feature="${BASH_REMATCH[1]}"
				OCP_FILTER=" && ocp-required"
				echo "Filtering to OCP-required tests only"
			elif [[ $feature =~ ^(.+)-k8s$ ]]; then
				feature="${BASH_REMATCH[1]}"
				OCP_FILTER=" && !ocp-required"
				echo "Filtering to K8S-compatible tests only (excluding OCP-required)"
			fi

			# Check if feature ends with a number (e.g., accesscontrol1)
			if [[ $feature =~ ^([a-z]+)([0-9]+)$ ]]; then
				base_feature="${BASH_REMATCH[1]}"
				# Find the base directory (e.g., accesscontrol for accesscontrol1)
				for dir in tests/*; do
					if [[ $dir != *"util"* ]] && [[ $dir == *"${base_feature}"* ]]; then
						command+=" "$dir
						LABEL_FILTER="${feature}"
					fi
				done
			else
				# Original behavior for non-numbered features
				for dir in tests/*; do
					if [[ $dir != *"util"* ]] && [[ $dir == *"${feature}"* ]]; then
						command+=" "$dir
					fi
				done
			fi
		done

		# strip any spaces from the command variable
		command=$(echo "$command" | xargs)
		# check if the command variable is empty
		if [ -z "$command" ]; then
			echo "No tests found for feature: ${FEATURES}"
			exit 1
		fi

		# Add label filter flag if set, combined with OCP filter
		LFLAG=""
		if [[ -n "$LABEL_FILTER" ]]; then
			LFLAG="--label-filter=${LABEL_FILTER}${OCP_FILTER}"
		elif [[ -n "$OCP_FILTER" ]]; then
			# If no label filter but OCP filter is set, use "true" as base
			LFLAG="--label-filter=true${OCP_FILTER}"
		fi

		# shellcheck disable=SC2086
		ginkgo -v ${PFLAG} ${FFLAG} ${LFLAG} --keep-going ${GINKGO_SEED_FLAG} --output-interceptor-mode=none --timeout=24h --show-node-events --require-suite $command
		;;
	*)
		echo "Unknown case"
		exit 1
		;;
	esac
}

run_tests "${1}"
