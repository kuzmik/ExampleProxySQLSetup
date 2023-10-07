# ProxySQL CLuster for Replicated Databases

## Setup

This assumes you have a k8s cluster on hand; we're installing some helm charts into it.

### Create the MySQL cluster

We'll create one primary and one secondary mysql instance, and make sure they are replicating. We're using the bitnami mysql charts for this.

```bash
kubectl get namespace | grep -q "^mysql" || kubectl create ns mysql

helm install mysql-us1 -n mysql ./helm/mysql-1

helm install mysql-us2 -n mysql ./helm/mysql-2
```

### Create the ProxySQL cluster

For this step, we're creating a proxysql controller statefulset and a proxysql cluster deployment. The controller is the "leader" and is in charge of distributing the configuration changes to the followers. The followers are configured to automatically connect to the leader.

```bash
# create the proxysql leader and followers in the proxysql namespace
kubectl get namespace | grep -q "^proxysql" || kubectl create ns proxysql

helm install proxysql-cluster-controller -n proxysql ./helm/proxysql/cluster-controller
helm install proxysql-cluster -n proxysql ./helm/proxysql/cluster-follower
```

-----

## Creds

### MySQL - US1 shard

* database: persona-web-us1_local
* username: persona-web--us1
* password: persona-web--us1

### MySQL - US2 shard

* database: persona-web-us2_local
* username: persona-web--us2
* password: persona-web--us2

### To Connect to MySQL Directly

This connects to the database as root, using the root password k8s secret (which is just... `rootpw`). You can connect to any service avaiable in the mysql namespace, as long as they have a ClusterIP (ie: not the headless services).

```bash
US1_IP="$(kubectl get services -n mysql mysql-us1-primary -o jsonpath='{.spec.clusterIP}')"
MYSQL_ROOT_PASSWORD=$(kubectl get secret --namespace mysql mysql-us2 -o jsonpath="{.data.mysql-root-password}" | base64 -d)

mysql -h$US1_IP -uroot -p$MYSQL_ROOT_PASSWORD
```

There is the assumption here that you have the `mysql` command available; if not, you can use the docker image from bitnami:

```bash
kubectl run mysql-us2-client --rm --tty -i --restart='Never' --image  docker.io/bitnami/mysql:8.0.34-debian-11-r56 --namespace mysql --env MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD --command -- bash
```

### Connecting to MySQL Backends via ProxySQL

Connect to us1 via proxysql

```bash
mysql -h$(k get service proxysql-cluster --output jsonpath='{.spec.clusterIP}') -P6033 -upersona-web-us1 -ppersona-web-us1
```

-----

## Teardown

```bash
helm uninstall -n proxysql proxysql-cluster
helm uninstall -n proxysql proxysql-cluster-controller

helm uninstall -n mysql mysql-us1
helm uninstall -n mysql mysql-us2

kubectl delete ns proxysql
kubectl delete ns mysql
```
