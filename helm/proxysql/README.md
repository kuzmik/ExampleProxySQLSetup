# ProxySQL on Kubernetes

These charts are adapted from [the offical k8s repo](https://github.com/ProxySQL/kubernetes), which I have been told is not technically canonical, and is in fact experimental. Because of that, I have felt empowered to delete a bunch of things that we don't need and make some changes to the rest.

* We're using the [latest proxysql docker image](https://hub.docker.com/r/proxysql/proxysql) rather than some random one they created.
* We're using the cluster controller and the cluster follower charts
