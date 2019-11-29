version: '3'

services:
  cadvisor:
    image: registry.jiangxingai.com:5000/cadvisor:arm64v8-cpu-0.1.0
    container_name: cadvisor
    hostname: {{.WORKER_ID}}
    restart: always
    ports: 
     - 3080:8080
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:rw
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
    command: -storage_driver=influxdb -storage_driver_host={{.MASTER_IP}}:8086 -storage_driver_db=telegraf
