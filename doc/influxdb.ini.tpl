[config]
influxdb_number = 2

[influxdb1]
host = {{.VpnIP}}
port = 8086
database = statsite
username = root
password = root
version = 1.7
timeout = 2

[influxdb2]
host = {{.MASTER_IP}}
port = 8086
database = statsite
username = root
password = root
version = 1.7
timeout = 2
