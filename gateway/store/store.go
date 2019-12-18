package store

import (
	// "log"
	"os"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"

	log "jxcore/lowapi/logger"
)

// Iterator 迭代器
type Iterator interface {
	// First moves the iterator to the first key/value pair. If the iterator
	// only contains one key/value pair then First and Last would moves
	// to the same key/value pair.
	// It returns whether such pair exist.
	First() bool

	// Last moves the iterator to the last key/value pair. If the iterator
	// only contains one key/value pair then First and Last would moves
	// to the same key/value pair.
	// It returns whether such pair exist.
	Last() bool

	// Seek moves the iterator to the first key/value pair whose key is greater
	// than or equal to the given key.
	// It returns whether such pair exist.
	//
	// It is safe to modify the contents of the argument after Seek returns.
	Seek(key []byte) bool

	// Next moves the iterator to the next key/value pair.
	// It returns false if the iterator is exhausted.
	Next() bool

	// Prev moves the iterator to the previous key/value pair.
	// It returns false if the iterator is exhausted.
	Prev() bool

	// Key returns the key of the current key/value pair, or nil if done.
	// The caller should not modify the contents of the returned slice, and
	// its contents may change on the next call to any 'seeks method'.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done.
	// The caller should not modify the contents of the returned slice, and
	// its contents may change on the next call to any 'seeks method'.
	Value() []byte
	Release()
}

// Store 存储配置
type Store interface {
	Open() error
	Get(key []byte) (value []byte, err error)
	Put(key, value []byte) error
	Delete(key []byte) error
	NewIterator(prefix string) Iterator
	Close() error
}

// LevelDBStore 使用 LevelDB 保存配置
type LevelDBStore struct {
	*leveldb.DB
	File string
	mu   sync.Locker
}

// NewLevelDBStore 打开指定路径的 LevelDB
func NewLevelDBStore(file string) *LevelDBStore {
	return &LevelDBStore{File: file, mu: new(sync.Mutex)}
}

// Open 打开DB文件
func (s *LevelDBStore) Open() error {
	db, err := leveldb.OpenFile(s.File, nil)
	if err != nil {
		if errors.IsCorrupted(err) {
			if err = os.RemoveAll(s.File); err != nil {
				return err
			}
			db, err = leveldb.OpenFile(s.File, nil)
			s.DB = db
		} else {

		}
		return err
	}
	s.DB = db

	return err
}

func (s *LevelDBStore) db() (db *leveldb.DB, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.DB == nil {
		err = s.Open()
		if err != nil {
			return
		}
	}
	return s.DB, nil
}

// Get 读取配置
func (s *LevelDBStore) Get(key []byte) ([]byte, error) {
	db, err := s.db()
	if err != nil {
		return nil, err
	}
	return db.Get(key, nil)
}

// Put 修改指定的 Key
func (s *LevelDBStore) Put(key, value []byte) error {
	db, err := s.db()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return db.Put(key, value, nil)
}

// Delete 删除指定的 Key
func (s *LevelDBStore) Delete(key []byte) error {
	db, err := s.db()
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return db.Delete(key, nil)
}

// NewIterator 指定前缀迭代查找 Key
func (s *LevelDBStore) NewIterator(prefix string) Iterator {
	db, err := s.db()
	if err != nil {
		log.Panic(err)
	}
	return db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
}

// Close 关闭数据库
func (s *LevelDBStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.DB == nil {
		return nil
	}
	err := s.DB.Close()
	s.DB = nil
	return err
}
