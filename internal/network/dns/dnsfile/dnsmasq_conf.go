package dnsfile

const DNSMasqConf = `
# This file is managed by Jxcore, please don't modify.

resolv-file=/etc/dnsmasq.resolv.conf
interface=
listen-address=172.18.1.1
bind-interfaces
conf-dir=/etc/dnsmasq.d
addn-hosts=/etc/dnsmasq.hosts
`
