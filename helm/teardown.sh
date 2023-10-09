#!/bin/bash
set -eou pipefail

helm uninstall -n proxysql proxysql-cluster
helm uninstall -n proxysql proxysql-cluster-controller

helm uninstall -n mysql mysql-us1
helm uninstall -n mysql mysql-us2

# Probably all we _really_ need to do here is delete the namespaces, but then helm might get confused
kubectl delete ns proxysql
kubectl delete ns mysql
