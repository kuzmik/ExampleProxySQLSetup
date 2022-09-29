### Add some user data

```
mysql --protocol tcp  -P 3307 -u root -pmysql1 < users.sql
```

### Login to the proxy 

```
mysql --protocol tcp -P 16033 -u monitoruser -pmonitorpass
mysql> use application;
Database changed
mysql> SELECT * FROM Users;
+--------+----------+-----------+--------------------+
| UserID | LastName | FirstName | Email              |
+--------+----------+-----------+--------------------+
|      1 | Lee      | Geddy     | geddy@rush.com     |
|      2 | Lifeson  | Alex      | alex@rush.com      |
+--------+----------+-----------+--------------------+

```

### Add a User

```
mysql> INSERT into Users values (3, "Peart", "Neil", "professor@rush.com");
mysql> SELECT * FROM Users;
+--------+----------+-----------+--------------------+
| UserID | LastName | FirstName | Email              |
+--------+----------+-----------+--------------------+
|      1 | Lee      | Geddy     | geddy@rush.com     |
|      2 | Lifeson  | Alex      | alex@rush.com      |
|      3 | Peart    | Neil      | professor@rush.com |
+--------+----------+-----------+--------------------+
```

### In a separate shell disable replication on the slave

```
mysql --protocol tcp  -P 3308 -u root -pmysql2
mysql> STOP SLAVE;
```

### In the proxy shell insert another user 

```
mysql> INSERT into Users values (4, "Clapton", "Eric", "eric@layla.com");
```

### Now since the SELECT is routed to the slave you won't see this 

```
mysql> SELECT * FROM Users;
+--------+----------+-----------+--------------------+
| UserID | LastName | FirstName | Email              |
+--------+----------+-----------+--------------------+
|      1 | Lee      | Geddy     | geddy@rush.com     |
|      2 | Lifeson  | Alex      | alex@rush.com      |
|      3 | Peart    | Neil      | professor@rush.com |
+--------+----------+-----------+--------------------+
```

### Restart the slave

```
mysql --protocol tcp  -P 3308 -u root -pmysql2
mysql> START SLAVE;
```

### Run the query again on the proxy and you'll see all the users

```
mysql> SELECT * FROM Users;
+--------+----------+-----------+--------------------+
| UserID | LastName | FirstName | Email              |
+--------+----------+-----------+--------------------+
|      1 | Lee      | Geddy     | geddy@rush.com     |
|      2 | Lifeson  | Alex      | alex@rush.com      |
|      3 | Peart    | Neil      | professor@rush.com |
|      4 | Clapton  | Eric      | eric@layla.com     |
+--------+----------+-----------+--------------------+
```


