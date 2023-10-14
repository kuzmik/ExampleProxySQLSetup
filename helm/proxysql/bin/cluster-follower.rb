#!/usr/bin/env ruby
# frozen_string_literal: true

# count the number of proxysql servers that:
#   - are not named "proxysql-leader" because that is the default server that is loaded from the config,
#     either on boot or by this script
#   - haven't been seen in over 30s (last_check_ms)
#   - have an uptime > 0
#     - when a follower first joins the cluster, and before the leaders propagate config, last_check_ms
#       will continue to grow but uptime will remain 0
missing_command = "SELECT count(hostname) FROM stats_proxysql_servers_metrics WHERE last_check_ms > 30000 and hostname != 'proxysql-leader' and Uptime_s > 0"
all_command = 'SELECT count(hostname) FROM stats_proxysql_servers_metrics'

missing_count = `mysql -h127.0.0.1 -P6032 -uadmin -padmin -NB -e"#{missing_command}"`.to_i
all_count = `mysql -h127.0.0.1 -P6032 -uadmin -padmin -NB -e"#{all_command}"`.to_i

# if there are any servers that are returned by the command, the entire leader cluster is probably down; it was either
# removed, or has been deployed recently and got different IPs. because we aren't using static IPs, we want to bootstrap
# the proxysql servers again. joining the cluster will pull down data from the leaders automatically, though it tends to
# take ~20 seconds or so once the leaders are back
if missing_count.positive?
  puts "#{missing_count}/#{all_count} proxysql leaders haven't been seen in over 30s, resetting leader state"

  commands = [
    'DELETE FROM proxysql_servers',
    'LOAD PROXYSQL SERVERS FROM CONFIG',
    'LOAD PROXYSQL SERVERS TO RUNTIME'
  ].join('; ')

  `mysql -h127.0.0.1 -P6032 -uadmin -padmin -NB -e"#{commands}"`
end
