#!/bin/bash
set -eou pipefail

# if we aren't in one of the orbstack/docker-desktop contexts, bail out. basically i want to prevent accidentally un-deploying
# this stuff to staging (or god help us, prod).
context=$(kubectl config current-context)
if [[ "$context" != "orbstack" ]] && [[ "$context" != "docker-desktop" ]]; then
  echo "You are not in the right kube context, current context is: $context. We want 'orbstack' or 'docker-desktop'"
  exit 1
fi

helm uninstall -n proxysql proxysql-leader
helm uninstall -n proxysql proxysql-cluster

helm uninstall -n mysql mysql-us1
helm uninstall -n mysql mysql-us2

# Probably all we _really_ need to do here is delete the namespaces, but then helm might get confused
kubectl delete ns proxysql
kubectl delete ns mysql
