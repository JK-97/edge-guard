package yaml

var defaultConfig = `
#cri: true
monitor:
  telegraf: true
#devicemanagement:
#  camera: true
#  rs485: true
#  aiserving: true
tools:
  mcutools:
    watchdog: true
    mcuserver: true
  nodetools:
    #cleaner: true
    #usblistener: true
    filelistener: true
jxserving: true
synctools:
  db: true
  tsdb: true
  mq: true
  fs: true

fixedresolver: ""
fsindex: true

iface:
  priority:
    - "eth0"
    - "eth1"
    - "usb0"
    - "usb1"
  backup: "usb0"
  switch_interval: 5s

mount_cfg:
    "/dev/mmcblk1p1": "/media/card"
`
