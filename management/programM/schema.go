package programM

type ProgramM struct {
    MyProgram []Program
}
type Program struct {
    Name          string
    Command       string
    LogPath       string
    startretries  int32
    startsecs     int64
    isAutoStart   bool
    isAutoRestart bool
}
