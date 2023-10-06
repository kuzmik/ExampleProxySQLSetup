# ProxySQL CLuster for Replicated Databases

## Setup

This assumes you have a k8s cluster on hand; we're installing some helm charts into it.

### Create the MySQL cluster

We'll create one primary and one secondary mysql instance, and make sure they are replicating. We're using the bitnami mysql charts for this.

```bash
# create the mysql primary and replica servers in the mysql namespace
kubectl create ns mysql
helm install mysql -n mysql ./helm/mysql
```

### Create the ProxySQL cluster

For this step, we're creating a proxysql controller statefulset and a proxysql cluster deployment. The controller is the "leader" and is in charge of distributing the configuration changes to the followers. The followers are configured to automatically connect to the leader.

```bash
# create the proxysql leader and followers in the proxysql namespace
kubectl create ns proxysql
helm install proxysql-cluster-controller -n proxysql ./helm/proxysql/cluster-controller
helm install proxysql-cluster -n proxysql ./helm/proxysql/cluster-follower
```

## Teardown

```bash
helm uninstall -n proxysql proxysql-cluster
helm uninstall -n proxysql proxysql-cluster-controller

helm uninstall -n mysql mysql

kubectl delete ns proxysql
kubectl delete ns mysql
```
