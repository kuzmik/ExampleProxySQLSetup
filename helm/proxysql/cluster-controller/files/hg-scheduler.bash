#!/bin/bash

CUR_CS=$(/usr/bin/md5sum $0 | /usr/bin/awk '{print $1}')
PRE_CS=$(/bin/cat /tmp/cs.file | /usr/bin/awk '{print $1}')

if [[ $CUR_CS != $PRE_CS ]];
then
  /bin/echo "Diffs detected, executing"
  /usr/bin/mysql -uadmin -padmin -h127.0.0.1 -P6032 -e"
    DELETE FROM proxysql_servers;
    INSERT INTO proxysql_servers (hostname, port, weight, comment)
      VALUES ('proxysql-cluster-controller', 6032, 0, 'proxysql-cluster-controller');
    LOAD PROXYSQL SERVERS TO RUNTIME;
    LOAD MYSQL SERVERS TO RUNTIME;
    LOAD MYSQL USERS TO RUNTIME;
    LOAD MYSQL QUERY RULES TO RUNTIME;
  "

  /bin/echo $CUR_CS > /tmp/cs.file
fi

