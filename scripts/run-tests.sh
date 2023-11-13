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
	PFLAG="-procs=4"
fi

# Allow for flake retries
FFLAG=""
if [[ ${ENABLE_FLAKY_RETRY} == "true" ]]; then
	echo "Retrying flaky tests"
	FFLAG="--flake-attempts=3"
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
		# shellcheck disable=SC2086
		ginkgo -timeout=24h -v ${PFLAG} ${FFLAG} --keep-going "${GINKGO_SEED_FLAG}" --require-suite -r $all_default_suites
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

		# shellcheck disable=SC2086
		ginkgo -timeout=24h -v ${PFLAG} ${FFLAG} --keep-going "${GINKGO_SEED_FLAG}" --output-interceptor-mode=none --require-suite $command
		;;
	*)
		echo "Unknown case"
		exit 1
		;;
	esac
}

run_tests "${1}"
