package programM

const (
    DependOnBase string= "depends_on=gateway"
)

var ProgramMconfig = `[supervisord]
logfile=%(here)s/logfile/jxcore.log
logfile_maxbytes=50MB
logfile_backups=10
loglevel=info
pidfile=/tmp/supervisord.pid
`

var BaseDepend = `[program:gateway]
command=/home/marshen/gateway -c /home/marshen/gateway.cfg
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
stdout_logfile=%(here)s/logfile/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=%(here)s/logfile/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
restart_when_binary_changed=true
`

var gateway = `[program:gateway]
command=/edge/gateway/bin/gateway -c /edge/gateway/bin/gateway.cfg
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
stdout_logfile=%(here)s/logfile/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=%(here)s/logfile/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
restart_when_binary_changed=true
`
var filelistener = `[program:filelistener]
depends_on=gateway
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
stopasgroup=true
killasgroup=true
user=root
redirect_stderr=false
stdout_logfile=%(here)s/logfile/%(program_name)s_stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
stdout_capture_maxbytes=0
stdout_events_enabled=true
stderr_logfile=%(here)s/logfile/%(program_name)s_stderr.log
stderr_logfile_maxbytes=50MB
stderr_logfile_backups=10
stderr_capture_maxbytes=0
stderr_events_enabled=false
`

var ProgramCfgMap = map[string]string{
    "Filelistener":   filelistener,
    //"gateway": gateway,
}

var ProgramSetting = ``
