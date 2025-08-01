#!/bin/bash

# Set the namespace
NAMESPACE="openshift-marketplace"

# Get the list of pods in the namespace
PODS=$(kubectl get pods -n "$NAMESPACE" -o name)

# Delete each pod
for pod in $PODS; do
	kubectl delete "$pod" -n "$NAMESPACE" --ignore-not-found
done
