package manager

import (
	"io"
	"io/ioutil"
	"jxcore/oplog/logs"
	"jxcore/oplog/types"
	"os"
	"strings"

	"sync"
)

type Manager struct {
	logFile *os.File
	mux     sync.Mutex
}

func NewManager(logPath string) *Manager {

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	return &Manager{
		logFile: logFile,
		mux:     sync.Mutex{},
	}
}

func (m *Manager) Insert(o types.Oplog) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	_, err := m.logFile.Write(o.Marshal())
	return err
}

func (m *Manager) Remove(o types.Oplog) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	return nil
}

func (m *Manager) ListAll() ([]types.Oplog, error) {
	return m.listByFiliter()
}
func (m *Manager) listByFiliter(f ...types.FilterFunc) ([]types.Oplog, error) {
	m.logFile.Seek(0, io.SeekStart)
	data, err := ioutil.ReadAll(m.logFile)
	if err != nil {
		return nil, err
	}
	rawData := strings.Split(strings.TrimSpace(string(data)), "\n")

	res := make([]types.Oplog, 0)
	for _, line := range rawData {
		logMessage := &logs.LogMessage{}
		logMessage.UnMarshal([]byte(line))
		if len(f) > 0 {
			for _, fn := range f {
				if !fn.Filter(logMessage) {
					continue
				}
			}
		}
		res = append(res, logMessage)

	}
	return res, nil
}

func (m *Manager) FindMany(f ...types.FilterFunc) ([]types.Oplog, error) {
	return m.listByFiliter(f...)
}

func (m *Manager) GetLogFileName() string {
	return m.logFile.Name()
}
