#!/bin/bash
set -eou pipefail

# Create the mysql infra

## Create the mysql namespace, unless it already exists
kubectl get namespace mysql > /dev/null 2>&1 \
  || kubectl create ns mysql

## Create some Configmaps that hold the mysql init scripts, if they don't already exist
kubectl get configmap -n mysql us1-initdb > /dev/null 2>&1 \
  || kubectl create configmap -n mysql us1-initdb --from-file=./helm/data/mysql-us1.sql
kubectl get configmap -n mysql us2-initdb > /dev/null 2>&1 \
  || kubectl create configmap -n mysql us2-initdb --from-file=./helm/data/mysql-us2.sql

## Install the mysql us1 and us2 instances, each of which has 1 replica
helm install mysql-us1 -n mysql ./helm/mysql \
  --set nameOverride="mysql-us1" \
  --set architecture="replication" \
  --set auth.rootPassword="rootpw" \
  --set auth.replicationPassword="replication" \
  --set auth.database="persona-web-us1" \
  --set auth.username="persona-web-us1" \
  --set auth.password="persona-web-us1" \
  --set initdbScriptsConfigMap="us1-initdb"

helm install mysql-us2 -n mysql ./helm/mysql \
  --set nameOverride="mysql-us2" \
  --set architecture="replication" \
  --set auth.rootPassword="rootpw" \
  --set auth.replicationPassword="replication" \
  --set auth.database="persona-web-us2" \
  --set auth.username="persona-web-us2" \
  --set auth.password="persona-web-us2" \
  --set initdbScriptsConfigMap="us2-initdb"

# End MySQL

echo "Sleeping 10s to allow mysql to finish coming up"
sleep 10

# Create the ProxySQL infra

## Create the ProxySQL namespace, unless it already exists
kubectl get namespace proxysql > /dev/null 2>&1 \
  || kubectl create ns proxysql

## ProxySQL controller cluster, which manages the configuration state of the rest of the cluster
helm install proxysql-cluster-controller -n proxysql ./helm/proxysql/cluster-controller

## ProxySQL "followers" cluster, which will be serving the actual proxied sql traffic
helm install proxysql-cluster -n proxysql ./helm/proxysql/cluster-follower

# End ProxySQL infra

## Cleanup the configmaps now that mysql is created... just in case we want to make
# changes to them, there won't be any confusion as to why they aren't taking effect
kubectl delete configmap -n mysql us1-initdb
kubectl delete configmap -n mysql us2-initdb
