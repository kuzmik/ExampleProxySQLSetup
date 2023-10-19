#!/bin/bash
set -eou pipefail

# Script to connect to the proxysql-core via mysql, without execing into the pod

mysql_info=$(kubectl get service -n proxysql proxysql-core --output=json | jq -r '.spec.clusterIP, .spec.ports[0].port')

mysql -h$(echo "$mysql_info" | awk 'NR==1') -P$(echo "$mysql_info" | awk 'NR==2') -uradmin -pradmin
