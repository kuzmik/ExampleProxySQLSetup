FROM proxysql/proxysql:2.5.5

# Install additional tools since ProxySQL Docker image is very bare-bones
# Install Ruby (for the probes.rb), mariadb-client so we can use the mysql cli, etc
RUN apt-get update && apt-get install -y --no-install-recommends --fix-missing \
  bind9-dnsutils \
  curl \
  jq \
  less \
  lsof \
  mariadb-client \
  procps \
  ruby \
  vim

# BUG: Reset entrypoint to override base image.
ENTRYPOINT []

CMD ["proxysql", "-f", "--idle-threads", "--no-version-check"]
