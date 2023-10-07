# Testing the cluster

## Connect to a Database Through ProxySQL

```bash
# us1
$ > mysql -h$(k get service proxysql-cluster --output jsonpath='{.spec.clusterIP}') -P6033 -upersona-web-us1 -ppersona-web-us1 -e 'select @@hostname'
+---------------------+
| @@hostname          |
+---------------------+
| mysql-us1-primary-0 |
+---------------------+

# us2
$ > mysql -h$(k get service proxysql-cluster --output jsonpath='{.spec.clusterIP}') -P6033 -upersona-web-us2 -ppersona-web-us2 -e 'select @@hostname'
+---------------------+
| @@hostname          |
+---------------------+
| mysql-us2-primary-0 |
+---------------------+
```

## Other stuff

```sql
SET mysql-eventslog_default_log = 1;
SET mysql-eventslog_format = 2; # json format
SET mysql-eventslog_filename = '/tmp/mysql_events.log';

# match all queries and do not log them to the event log
INSERT INTO mysql_query_rules (active, match_digest, log, apply) VALUES (1, '.', 0, 0);

# match select count queries and DO log them, but don't redirect them
insert into mysql_query_rules (active, match_pattern, log, apply) values (1, 'select count*', 1, 1);
```

