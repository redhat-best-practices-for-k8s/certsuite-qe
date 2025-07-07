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
		for feature in ${FEATURES}; do
			for dir in tests/*; do
				if [[ $dir != *"util"* ]] && [[ $dir == *"${feature}"* ]]; then
					command+=" "$dir
				fi
			done
		done

		# strip any spaces from the command variable
		command=$(echo "$command" | xargs)
		# check if the command variable is empty
		if [ -z "$command" ]; then
			echo "No tests found for feature: ${FEATURES}"
			exit 1
		fi

		# shellcheck disable=SC2086
		ginkgo -v ${PFLAG} ${FFLAG} --keep-going ${GINKGO_SEED_FLAG} --output-interceptor-mode=none --timeout=24h --show-node-events --require-suite $command
		;;
	*)
		echo "Unknown case"
		exit 1
		;;
	esac
}

run_tests "${1}"
