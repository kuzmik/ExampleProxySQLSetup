# TODO

- [ ] On deploy, followers log a lot and then lose their minds if the controllers are not reachable
  - If the controllers go away after the followers are running, the followers will log a lot, but not crash
  - If the controllers come back, the followers will still be unahppy because the IPs will be different for the new controllers. Maybe we need a lightweight cluster.rb for the followers too?

# DONE

- [X] Make proxysql run `cluster.rb` via the scheduler on the controllers to maintain cluster state
- [X] End to end connections working (mysql -> proxysql -> mysql backends) :party:
- [X] Basic mysql schema with some test data
- [X] Use only one mysql helm chart, and specify the values on the commandline
- [X] Simple setup/teardown scripts for easier bootstrapping
- [X] Fix the proxysql healthchecks - good enough
