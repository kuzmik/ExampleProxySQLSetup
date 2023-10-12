#!/bin/bash
set -e

mbin="/usr/bin/mysql"
lcon="-h127.0.0.1 -P6032 -uadmin -padmin"
opts="-NB"

hg1_avail=$($mbin $lcon $opts -e"select count(*) from runtime_mysql_servers where status = 'ONLINE'")

if [[ $hg1_avail -gt 0 ]];
then
  echo "Backends available"
  exit 0
else
  echo "Backends availability failure - MySQL backends found: $hg1_avail"
  exit 1
fi
