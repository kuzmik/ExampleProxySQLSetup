#!/bin/bash
set -e

mbin="/usr/bin/mysql"
lcon="-h127.0.0.1 -P6032 -uadmin -padmin"
opts="-NB"

hg1_avail=$($mbin $lcon $opts -e"select count(*) from runtime_mysql_servers where status = 'ONLINE'")

if [[ $hg1_avail -gt 1 ]];
then
  echo "hg1 Availability Success"
  exit 0
else
  echo "hg1 Availability Failure - MySQL backends found: $hg1_avail"
  exit 1
fi

# FIXME: THIS IS DOGSHIT AND NEEDS FIXING
