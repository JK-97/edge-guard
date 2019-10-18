package programmanage

const (
    DependOnBase string = ""
)

var ProgramMconfig = `[supervisord]
logfile=/edge/logs/jxcore.log
logfile_maxbytes=50MB
logfile_backups=10
loglevel=info
pidfile=/tmp/supervisord.pid
`

var BaseDepend = ``

var filelistener = `[program:filelistener]
#directory=/edge/tools/nodetools/filelistener/bin/
restart_when_binary_changed=true
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

var telegraf = `[program:telegraf]
#directory=/edge/monitor/telegraf/bin
restart_when_binary_changed=true
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
restart_when_binary_changed=true
`

var Db = `[program:Db]
#directory=/edge/mnt/db/bin
restart_when_binary_changed=true
command=/edge/mnt/db/bin/db serve --repo mongodb://172.17.0.1:27017 --src mongodb://172.17.0.1:27017
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
restart_when_binary_changed=true
`

var Powermanagement = `[program:Powermanagement]
restart_when_binary_changed=true
command=/edge/tools/mcutools/powermanagement/bin/powermanagement
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
restart_when_binary_changed=true
`

var Watchdog = `[program:WatchDog]
restart_when_binary_changed=true
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false

`

var Mcuserver = `[program:McuServer]
restart_when_binary_changed=true
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
`
var Mq = `[program:Mq]
restart_when_binary_changed=true
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
`
var Jxserving = `[program:Jxserving]
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

var Cleaner = `[program:Cleaner]
restart_when_binary_changed=true
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
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=/edge/logs/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
`
var ProgramCfgMap = map[string]string{
    "Filelistener":    filelistener,
    "Telegraf":        telegraf,
    "Db":              Db,
    "Mcuserver":       Mcuserver,
    "Powermanagement": Powermanagement,
    "Watchdog":        Watchdog,
    "Mq":              Mq,
    "Jxserving":       Jxserving,
    "Cleaner":         Cleaner,
}

var ProgramSetting = ``
