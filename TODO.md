# TODO and/or TOFIGUREOUT

- [ ] time to tackle mysql_query_rules

# DONE

- [X] Make proxysql run `maintain-cluster.rb` via the scheduler on the controllers to maintain cluster state
- [X] End to end connections working (mysql -> proxysql -> mysql backends) :party:
- [X] Basic mysql schema with some test data
- [X] Use only one mysql helm chart, and specify the values on the commandline
- [X] Simple setup/teardown scripts for easier bootstrapping
- [X] Fix the proxysql healthchecks - good enough
- [X] FIGUREOUT: on deploy (or if the leaders get cycled out and change IP addresses), followers log a lot and then lose their minds if the controllers are not reachable
  - proxysql in the cluster still functions normally, but it fills the logs with 1 error per proxysql_server, per second
  - follower self-healing done via `maintain-cluster.rb` in the cluster chart
  - leader self-healing was already done, also via (a different) `maintain-cluster.rb` process
  - both leader and followers run the cluster script every 10s via the proxysql scheduler config
- [X] FIGUREOUT: What happens if a follower config is updated via the admin interface?
  - does it get wiped out, or does it get put into a state that won't pull changes from the leaders?
  - yes, config gets wiped out as soon as the leaders run the maintain-cluster script
  - to test:
    1. i added a mysql rule on 1 follower pod. the rule stuck around...
    2. until i scaled down the leders by 1 pod, which triggered their maintain-cluster.rb
    3. maintain-cluster.rb ran `LOAD MYSQL QUERY RULES TO RUNTIME` which propagated an empty query rules set to the followers
    4. RIP my test rule on the follower pod
  - i think this is perfectly acceptable, we should not expect followers to have their own state, and should fully expect any changes NOT initiated by a leader to get wiped out at any time
