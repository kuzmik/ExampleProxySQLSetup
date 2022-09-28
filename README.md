## Starting up

```
docker-compose rm --force
docker-compose up

# Wait
./setup_replication.sh

```

## Checking the ProxySQL

```
mysql -u radmin -pradmin --protocol=tcp -h 127.0.0.1 -P16032 --prompt='RAdmin> '
```

## Connecting to the two MySQL DBs

```
mysql --protocol tcp  -P 3307 -u root -pmysql1
mysql --protocol tcp  -P 3308 -u root -pmysql2
```

## Connectiong through the proxy

```
mysql --protocol tcp -P 16033 -u monitoruser -pmonitorpass
```
