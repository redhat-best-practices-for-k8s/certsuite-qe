#!/bin/bash

# A script that will pre-emptively delete namespaces that are used in QE testing to make
# sure that the namespaces are not left behind after the test run.

# Check if oc is installed, if not exit gracefully
if ! command -v oc &>/dev/null; then
	echo >&2 "oc is not installed. Skipping deletion of namespaces"
	exit 0
fi

NAMESPACE_STRINGS_TO_GREP=(accesscontrol ac-test ac-rq-test my-ns affiliated lifecycle-tests manageability networking net-tests observability operator-ns performance platform-alteration cnfsuite)

for NS in "${NAMESPACE_STRINGS_TO_GREP[@]}"; do
	for NAMESPACE in $(oc get namespaces | grep "$NS" | awk '{print $1}'); do
		echo "Deleting namespace $NAMESPACE"
		oc delete namespace "$NAMESPACE" --ignore-not-found=true
	done
done
