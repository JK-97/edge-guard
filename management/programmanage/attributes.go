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
#numprocs_start=not support
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
#stopasgroup=true
#killasgroup=true
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
#numprocs_start=not support
autostart=true
startsecs=3
startretries=3
autorestart=true
exitcodes=0,2
stopsignal=TERM
stopwaitsecs=10
#stopasgroup=true
#killasgroup=true
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
command=/edge/mnt/db/bin/sync-db-arm64 serve --repo mongodb://172.17.0.1:27017 --src mongodb://172.17.0.1:27017
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
#stopasgroup=true
#killasgroup=true
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
var ProgramCfgMap = map[string]string{
	"Filelistener": filelistener,
	"Telegraf":     telegraf,
	"Db":Db,
	//"gateway": gateway,
}

var ProgramSetting = ``
