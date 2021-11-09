#!/usr/bin/env bash

GOPATH="${GOPATH:-~/go}"
export PATH=$PATH:$GOPATH/bin
export ACK_GINKGO_DEPRECATIONS=1.16.5
EXCLUDED_FOLDERS=""
ALL_TESTS_FOLDERS=$(ls -d ./tests/*/)

function run_tests {
    case $1 in
        all)
            echo "#### Run all tests ####"
            all_default_suites=""
            for folder in ${ALL_TESTS_FOLDERS}
            do
              for excluded_folder in ${EXCLUDED_FOLDERS}
                do
                  if [[ $folder == *"${excluded_folder}"* ]]; then
                    folder=''
                  fi
                done
                if ! [ -z "$folder" ]; then
                  all_default_suites+=" $folder"
                fi
            done
            ginkgo -v --keepGoing -requireSuite -r $all_default_suites
            ;;
        features)
            if [ -z "$FEATURES" ]; then {
                echo "FEATURES env var is empty. Please export FEATURES"
                exit 1
            } fi
            echo "#### Run feature tests: ${FEATURES} ####"
            for feature in ${FEATURES}
                do
                    for dir in tests/*/*; do
                    if [[ $dir != *"util"* ]] && [[ $dir == *"${feature}"* ]]; then {
                        command+=" "$dir
                    } fi
                    done
                done
            ginkgo -v --keepGoing -requireSuite $command
            ;;
        *)
        echo "Unknown case"
        exit 1
        ;;
    esac
}

run_tests ${1}
