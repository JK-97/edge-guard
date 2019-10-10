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
directory=/edge/tools/nodetools/filelistener/bin/
depends_on=gateway
command=./filelistener -c filelistener.cfg
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

var telegraf = `[program:telegraf]
directory=/edge/monitor/telegraf/bin
depends_on=gateway
command=./telegraf --config telegraf.cfg
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
	"Filelistener": filelistener,
	"Telegraf":     telegraf,
	//"gateway": gateway,
}

var ProgramSetting = ``
