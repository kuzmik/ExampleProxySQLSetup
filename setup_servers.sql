INSERT INTO mysql_servers(hostgroup_id,hostname,port) VALUES (1,'mysql1',3307);
INSERT INTO mysql_servers(hostgroup_id,hostname,port) VALUES (1,'mysql2',3308);
UPDATE global_variables SET variable_value='monitoruser' WHERE variable_name='mysql-monitor_username';
UPDATE global_variables SET variable_value='monitorpass' WHERE variable_name='mysql-monitor_password';
UPDATE global_variables SET variable_value='2000' WHERE variable_name IN ('mysql-monitor_connect_interval','mysql-monitor_ping_interval','mysql-monitor_read_only_interval');
LOAD MYSQL VARIABLES TO RUNTIME;
SAVE MYSQL VARIABLES TO DISK;
