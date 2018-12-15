package player

import (
	"os"
	"errors"
	//"log"
	"sync"
	"time"
	"regexp"
	"strconv"
	"fmt"
	"io/ioutil"
)

const FileCachTime = 5 * 60

type FileCache struct{
	fileName			string
	lastVisitTime 		int64
	len                 int64
	content             []byte
}

type FileManager struct{
	files		map[string]*FileCache
	mutex       *sync.Mutex
	cachePool   *sync.Pool
}

func NewFileManager() *FileManager{
	m := &FileManager{files:make(map[string]*FileCache)}
	m.mutex = new(sync.Mutex)
	m.cachePool = new(sync.Pool)
	return m
}

func (m *FileManager) newCache() *FileCache{
	m.mutex.Lock()
	defer m.mutex.Unlock()
	cache := m.cachePool.Get()
	if cache != nil {
		return cache.(*FileCache)
	}
	newCache := new(FileCache)
	return newCache
}

func (m *FileManager) getFile(fileName string) *FileCache {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if file , ok := m.files[fileName] ; ok {
		file.lastVisitTime = time.Now().Unix()
		return file
	}
	return nil
}

func (m *FileManager) setFile(fileName string, cache* FileCache) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	cache.lastVisitTime = time.Now().Unix()
	m.files[fileName] = cache
}

func (m *FileManager) clearExpiredFile() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	expireTimePoint := time.Now().Unix() - FileCachTime
	for fileName , cache := range m.files {
		if cache.lastVisitTime < expireTimePoint {
			delete(m.files , fileName)
			m.cachePool.Put(cache)
		}
	}
}

var fileManager *FileManager = NewFileManager()

func init() {
	go func(){
		for true {
			fileManager.clearExpiredFile()
			time.Sleep(FileCachTime * time.Second)
		}
	}()
}


type FileReader struct {
	FileName	string
	StartPos	int64
	EndPos		int64
	FileLen     int64
}



func NewFileReader(fileName, frange string) *FileReader {
	reader := &FileReader{FileName:fileName,StartPos:0,EndPos:-1,FileLen:0}
	re := regexp.MustCompile("bytes=([0-9]+)-([0-9]*)")
	matches := re.FindStringSubmatch(frange)
	if matches != nil && len(matches) == 3 {
		startPos , err := strconv.ParseInt(string(matches[1]) , 10 , 64)
		if err == nil {
			reader.StartPos = startPos
		}
		endPos , err := strconv.ParseInt(string(matches[2]) , 10 ,64)
		if err == nil{
			reader.EndPos = endPos
		}
	}
	return reader
}

func (r *FileReader) readFromCache(cache *FileCache)([]byte , error) {
	if r.StartPos > cache.len - 1 {
		return nil , errors.New("start pos > file len")
	}
	if r.EndPos == -1 {
		r.EndPos = cache.len - 1
	}
	if r.StartPos > r.EndPos {
		return nil , errors.New("start pos > end pos")
	}
	//go slice的区间是函数是左闭右开
	return cache.content[r.StartPos:r.EndPos+1] , nil
}

func (r *FileReader) Read() ([]byte,error) {
	cache := fileManager.getFile(r.FileName)
	if cache != nil {
		return r.readFromCache(cache)
	}
	file , err := os.Open(r.FileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("open file:%s fail" , r.FileName))
	}

	newCache := fileManager.newCache()
	newCache.fileName = r.FileName
	newCache.content , err  = ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("read file:%s fail",r.FileName))
	}
	fileManager.setFile(r.FileName , newCache)
	return r.readFromCache(newCache)
}
