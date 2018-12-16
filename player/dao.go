package player

import(
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "os"
    "path/filepath"
    "log"
    "fmt"
    "path"
    "encoding/json"
)

const (
     DB_FILE_NAME = "p.db"
     DB_FILE_PATH = "db"
)


func PathExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    if os.IsNotExist(err) {
        return false
    }
    return false
}



func openDB() *sql.DB {
    dbPath  := filepath.Join(".",DB_FILE_PATH,DB_FILE_NAME)
    if !PathExists(dbPath) {
        //第一次启动初始化DB
        db, err := sql.Open("sqlite3", dbPath)
        if err != nil {
            log.Fatal(err)
        }
        sqlStmtList := []string{
            "create table kv (k text not null primary key, v text);",
            "create table love(id integer PRIMARY KEY autoincrement  , vid text unique);",
            "create table hate(id integer PRIMARY KEY autoincrement  , vid text unique);"}
        for _,sqlStmt := range sqlStmtList {
            _, err = db.Exec(sqlStmt)
            if err != nil {
                log.Fatal(err)
            } else {
                log.Printf("excute %s succ\n" , sqlStmt)
            }
        }
        return db
    } else {
        db, err := sql.Open("sqlite3", dbPath)
        if err != nil {
            log.Fatal(err)
        }
        return db
    }
}

var db *sql.DB = openDB()

func GetKVMap() map[string]string {
    rows, err := db.Query("select k, v from kv")
    if err != nil {
        log.Print(err)
        return nil
    }
    kvMap := make(map[string]string , 5 )
    defer rows.Close()
    for rows.Next() {
        var k , v string
        err = rows.Scan(&k, &v)
        if err != nil {
            log.Print(err)
            return nil
        }
        kvMap[k] = v
    }
    err = rows.Err()
    if err != nil {
        log.Print(err)
    }
    return kvMap
}



func SetKV(key , val string) {
    sqlStmt := "replace into kv(k,v) values(?,?)"
    db.Exec(sqlStmt , key , val)
}

func GetConfig() string {
    kvMap := GetKVMap()
    config:= kvMap["config"]
    if config == "" {
        var dirs VDirsConf
        jBytes,_ := json.Marshal(dirs)
        return string(jBytes)
    }
    return config
}

func SetConfig(val string) {
    SetKV("config" , val)
}

type DirConf  struct {
    Id                int64           `json:"id"`
    Path            string          `json:"path"`
    Info            string          `json:"desc"`
    Files           []string              `json:"-"`
    Vids            []string              `json:"-"`
    Vid2File        map [string]string    `json:"-"`
}

type VDirsConf struct {
    Dirs            []DirConf      `json:"dirs"`
    ConfAccessKey   string         `json:"confAccessKey"`
}

func GetDirFileList(path string) []string {
    files := make([]string , 100)
    filterFile := func(path string, info os.FileInfo, err error) error{
        if !info.IsDir() {
            files = append(files , path)
        }
        return nil
    }
    filepath.Walk(path , filterFile)
    return files
}

func GetDirConfList() []DirConf  {
    var dirs VDirsConf
    err := json.Unmarshal([]byte(GetConfig()) , &dirs)
    if err != nil {
        return []DirConf{}
    }
    for i := 0; i < len(dirs.Dirs) ; i++ {
        dir := &dirs.Dirs[i]
        dir.Files = make([]string,0)
        dir.Vid2File = make(map[string]string)
        dir.Files = GetDirFileList(dir.Path)
        for _ , filePath := range dir.Files {
            fileName := path.Base(filePath)
            if fileName != "." {
                dir.Vid2File[fileName] = filePath
                dir.Vids = append(dir.Vids , fileName)
            }
        }
    }
    return dirs.Dirs
}

func GetPathByVid(vid string) string {
    dirList := GetDirConfList()
    for _ , dir := range dirList {
        if filePath , ok := dir.Vid2File[vid]; ok {
            return filePath
        } 
    }
    return ""
}

func GetDirConfById(id int64) *DirConf {
    dirList := GetDirConfList()
    for _ , dir := range dirList {
        if dir.Id == id {
            return &dir
        }
    }
    return nil
}
func AddVid(table,vid string) {
    sqlStmt := fmt.Sprintf("insert into %s (vid)  values(?)",table)
    _ , err := db.Exec(sqlStmt,vid)
    if err != nil {
        log.Print(err)
    }
}

func DelVid(table,vid string) {
    sqlStmt := fmt.Sprintf("delete from %s where vid = ?" , table)
    _ , err := db.Exec(sqlStmt,vid)
    if err != nil {
        log.Print(err)
    }
}

func GetVidList(table string) []string  {
    log.Print("GetVidList")
    data := make([]string, 0 ,100)
    rows, err := db.Query(fmt.Sprintf("select vid from %s" , table))
    if err != nil {
        log.Print(err)
        return data
    }
    defer rows.Close()
    for rows.Next() {
        var vid string
        err = rows.Scan(&vid)
        if err != nil {
            log.Print(err)
            return data
        }
        data = append(data , vid)
    }
    err = rows.Err()
    if err != nil {
        log.Print(err)
    }
    return data    
}

func GetFavoriteVidList() []string {
    log.Print("GetFavoriteVidList")
    return GetVidList("love")
}

func GetHateVidList() []string {
    return GetVidList("hate")
}