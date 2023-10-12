#!/bin/bash
set -eou pipefail

DIR=$(dirname -- "${BASH_SOURCE[0]}")

source "$DIR/.lib/assert.sh"

proxysql_instance=$(kubectl get service proxysql-cluster -n proxysql --output jsonpath='{.spec.clusterIP}')

echo "----- us1 primary -----"

us1_hostname=$(mysql --defaults-extra-file="$DIR/.lib/us1-client.cfg" -h"$proxysql_instance" -P6033 -NB -e 'select @@hostname')
assert_eq "mysql-us1-primary-0" "$us1_hostname" "Not equivalent"
echo "us1_primary hostname: mysql-us1-primary-0 == $us1_hostname"

us1_primary_ro=$(mysql --defaults-extra-file="$DIR/.lib/us1-client.cfg" -h"$proxysql_instance" -P6033 -NB -e 'select @@global.read_only')
assert_eq "0" "$us1_primary_ro" "Not equivalent"
echo "us1_primary readonly flag: 0 == $us1_primary_ro"

us1_user=$(mysql --defaults-extra-file="$DIR/.lib/us1-client.cfg"  -h"$proxysql_instance" -P6033  persona-web-us1 -NB -e 'select email from users where id = 1')
assert_eq "rick@persona-us1.com" "$us1_user" "Not equivalent"
echo "us1_primary user: rick@persona-us1.com == $us1_user"

echo ''
echo "----- us1 secondary -----"

us1_ro_hostname=$(mysql --defaults-extra-file="$DIR/.lib/us1-ro-client.cfg" -h"$proxysql_instance" -P6033 -NB -e 'select @@hostname')
assert_eq "mysql-us1-secondary-0" "$us1_ro_hostname" "Not equivalent"
echo "us1_secondary hostname: mysql-us1-primary-0 == $us1_ro_hostname"

us1_secondary_ro=$(mysql --defaults-extra-file="$DIR/.lib/us1-ro-client.cfg" -h"$proxysql_instance" -P6033 -NB -e 'select @@global.read_only')
assert_eq "1" "$us1_secondary_ro" "Not equivalent"
echo "us1_secondary readonly flag: 1 == $us1_secondary_ro"

us1_secondary_user=$(mysql --defaults-extra-file="$DIR/.lib/us1-client.cfg"  -h"$proxysql_instance" -P6033  persona-web-us1 -NB -e 'select email from users where id = 1')
assert_eq "rick@persona-us1.com" "$us1_user" "Not equivalent"
echo "us1_secondary user: rick@persona-us1.com == $us1_user"

echo ''
echo "----- us2 primary -----"

us2_hostname=$(mysql --defaults-extra-file="$DIR/.lib/us2-client.cfg" -h"$proxysql_instance" -P6033 -NB -e 'select @@hostname')
assert_eq "mysql-us2-primary-0" "$us2_hostname" "Not equivalent"
echo "us2_primary hostname: mysql-us2-primary-0 == $us2_hostname"

us2_primary_ro=$(mysql --defaults-extra-file="$DIR/.lib/us2-client.cfg" -h"$proxysql_instance" -P6033 -NB -e 'select @@global.read_only')
assert_eq "0" "$us2_primary_ro" "Not equivalent"
echo "us2_primary readonly flag: 0 == $us2_primary_ro"

us2_user=$(mysql --defaults-extra-file="$DIR/.lib/us2-client.cfg"  -h"$proxysql_instance" -P6033  persona-web-us2 -NB -e 'select email from users where id = 1')
assert_eq "charles@persona-us2.com" "$us2_user" "Not equivalent"
echo "us2_primary user: charles@persona-us2.com == $us2_user"

echo ''
echo "----- us2 secondary -----"

us2_ro_hostname=$(mysql --defaults-extra-file="$DIR/.lib/us2-ro-client.cfg" -h"$proxysql_instance" -P6033 -NB -e 'select @@hostname')
assert_eq "mysql-us2-secondary-0" "$us2_ro_hostname" "Not equivalent"
echo "us2_secondary hostname: mysql-us2-primary-0 == $us2_ro_hostname"

us2_secondary_ro=$(mysql --defaults-extra-file="$DIR/.lib/us2-ro-client.cfg" -h"$proxysql_instance" -P6033 -NB -e 'select @@global.read_only')
assert_eq "1" "$us2_secondary_ro" "Not equivalent"
echo "us2_secondary readonly flag: 1 == $us2_secondary_ro"

us2_secondary_user=$(mysql --defaults-extra-file="$DIR/.lib/us2-client.cfg"  -h"$proxysql_instance" -P6033  persona-web-us2 -NB -e 'select email from users where id = 1')
assert_eq "charles@persona-us2.com" "$us2_user" "Not equivalent"
echo "us2_secondary user: charles@persona-us2.com == $us2_user"

echo ''
echo '-*-*-*-*- all tests passed, insert emojis here -*-*-*-*-'
echo ''
