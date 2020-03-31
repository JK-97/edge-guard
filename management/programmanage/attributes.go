package programmanage

const (
	DependOnBase string = ""
)

var ProgramMconfig = `[supervisord]
logfile=/edge/logs/edge-guard.log
logfile_maxbytes=50MB
logfile_backups=10
loglevel=info
pidfile=/tmp/supervisord.pid
[inet_http_server]
port = :9001
`

var BaseDepend = ""

var filelistener = `[program:filelistener]
#directory=/edge/tools/nodetools/filelistener/bin/
command=/edge/tools/nodetools/filelistener/bin/filelistener -c /edge/tools/nodetools/filelistener/bin/filelistener.cfg
process_name=%(program_name)s
numprocs=1
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

var telegraf = `[program:telegraf]
#directory=/edge/monitor/telegraf/bin
command=/edge/monitor/telegraf/bin/telegraf --config /edge/monitor/telegraf/bin/telegraf.cfg
process_name=%(program_name)s
numprocs=1
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

var Db = `[program:db]
#directory=/edge/mnt/db/bin
command=/edge/mnt/db/bin/db serve --config /edge/mnt/db/bin/db.toml
process_name=%(program_name)s
numprocs=1
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

// Watchdog 与 mcu 通信，防止系统被误杀
var Watchdog = `[program:watchDog]
command=/edge/tools/mcutools/watchdog/bin/watchdog
process_name=%(program_name)s
numprocs=1
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false

`

// Mcuserver 与 mcu 通信，获取当前的开机模式
var Mcuserver = `[program:mcuServer]
command=/edge/tools/mcutools/mcuserver/bin/mcuserver
process_name=%(program_name)s
numprocs=1
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

// Mq 消息队列同步
var Mq = `[program:mq]
command=/edge/mnt/mq/bin/mq -c /edge/mnt/mq/bin/mq.cfg
process_name=%(program_name)s
numprocs=1
#numprocs_start=not support
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
stopasgroup=true
killasgroup=true
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

// Fs 文件系统同步
var Fs = `[program:fs]
command=/edge/mnt/fs/bin/fs
process_name=%(program_name)s
numprocs=1
#numprocs_start=not support
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
stopasgroup=true
killasgroup=true
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`
var Jxserving = `[program:jxserving]
command=python3 /jxserving/run.py
process_name=%(program_name)s
numprocs=1
#numprocs_start=not support
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
stopasgroup=true
killasgroup=true
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

var Cleaner = `[program:cleaner]
command=/edge/tools/nodetools/cleaner/bin/cleaner -c /edge/tools/nodetools/cleaner/bin/cleaner.cfg
process_name=%(program_name)s
numprocs=1
#numprocs_start=not support
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
stopasgroup=true
killasgroup=true
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

var Fsindex = `[program:fsindex]
directory=/edge/fsindex/bin/
command=/edge/fsindex/bin/fsindexer 
process_name=%(program_name)s
numprocs=1
#numprocs_start=not support
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
stopasgroup=true
killasgroup=true
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

var Speaker = `[program:speaker]
directory=/edge/tools/speaker/
command=/edge/tools/speaker/device-speaker
process_name=%(program_name)s
numprocs=1
#numprocs_start=not support
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
stopasgroup=true
killasgroup=true
user=root
redirect_stderr=false
stdout_logfile=/edge/logs/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=1
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=1
stderr_capture_maxbytes=0
stderr_events_enabled=false
`
var Cadvisor = `[program:cadvisor]
command=docker-compose -f /jxbootstrap/worker/docker-compose.d/cadvisor/docker-compose.yaml up 
process_name=%(program_name)s
numprocs=1
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
stopasgroup=true
killasgroup=true
user=root
`

var ProgramCfgMap = map[string]string{
	"filelistener": filelistener,
	"telegraf":     telegraf,
	"db":           Db,
	"fs":           Fs,
	"mcuserver":    Mcuserver,
	"watchdog":     Watchdog,
	"mq":           Mq,
	"jxserving":    Jxserving,
	"cleaner":      Cleaner,
	"cadvisor":     Cadvisor,
	"fsindex":      Fsindex,
	"speaker":      Speaker,
}

var ProgramSetting = ``
