package programM

const (
    DependOnBase = "depends_on=mongo,influx"
)

var ProgramMconfig = `[supervisord]
logfile=%(here)s/logfile/jxcore.log
logfile_maxbytes=50MB
logfile_backups=10
loglevel=info
pidfile=/tmp/supervisord.pid
`

var BaseDepend = `[program:mongo]
command=mongod /etc/mongod.conf
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
`

var gateway = `[program:gateway]
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
var nginx = `[program:nginx]
command=docker run nginx
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
`

var ProgramCfgMap = map[string]string{
    "nginx":   nginx,
    "gateway": gateway,
}

var ProgramSetting = ``
