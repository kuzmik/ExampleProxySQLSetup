# ProxySQL on Kubernetes

These charts are adapted from [the offical k8s repo](https://github.com/ProxySQL/kubernetes), which I have been told is not technically canonical, and is in fact experimental. Because of that, I have felt empowered to delete a bunch of things that we don't need and make some changes to the rest.

* We're using the [latest proxysql docker image](https://hub.docker.com/r/proxysql/proxysql) rather than some random one they created.
* Well ACTUALLY, I've built a [Dockerfile](Dockerfile) that automatically installs a bunch of tooling so that I don't have to keep doing it over and over everytime I reinstall.


## Controller Cluster

Charts for this deployment (statefulset technically) are in the [cluster-controller](cluster-controller) directory.

Install the controllers via:

```shell
# create the namespace if it doesn't exist
kubectl create ns proxysql

helm install proxysql-cluster-controller -n proxysql ./helm/proxysql/cluster-controller
```

The controller cluster is a small (3 nodes currently, but can be scaled down as needed) cluster of proxysql servers that communicate with each other. These instances will not serve proxysql traffic, and exist only to manage the configuration of the real proxysql cluster.

Configuration changes on any of these pods will be propagated to the other pods in the controller cluster, and to any follower cluster that is online.

Resources created by the charts:

```
23:19:31 <nick@marais:ExampleProxySQLSetup(kuzmik/k8s-cluster)(✘!?) $ > k get all
NAME                                READY   STATUS    RESTARTS   AGE
pod/proxysql-cluster-controller-0   1/1     Running   0          28s
pod/proxysql-cluster-controller-1   1/1     Running   0          26s
pod/proxysql-cluster-controller-2   1/1     Running   0          15s

NAME                                  TYPE        CLUSTER-IP        EXTERNAL-IP   PORT(S)    AGE
service/proxysql-cluster-controller   ClusterIP   192.168.194.163   <none>        6032/TCP   28s

NAME                                           READY   AGE
statefulset.apps/proxysql-cluster-controller   3/3     28s
```

## Follower cluster

Charts for this deployment are in the [cluster-follower](cluster-follower) directory.

Install the followers via:

```shell
helm install proxysql-cluster -n proxysql ./helm/proxysql/cluster-follower
```

This is the actual proxysql cluster that will serve traffic. On boot, each pod will connect to the `proxysql-cluster-controller` service (see above), which will distribute the configuration to the pod. This will allow scaling in and out to be easier and any new pods will automatically join the "cluster" so to speak.

Configuration changes on these pods will NOT propagate up to the contoller cluster, and therefore will not make it to any other proxysql-follower pod.
