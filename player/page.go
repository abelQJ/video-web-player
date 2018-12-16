package player

import (
    "net/http"
    "html/template"
    "path/filepath"
    "strconv"
    "log"
)


var drivers *template.Template

func init(){
    pattern := filepath.Join("." , "tpl" , "*.html")
    drivers = template.Must(template.ParseGlob(pattern))
}


type PageCommon struct {
    Title    string
}

type ConfigPage struct {
    PageCommon
}

type HomePage struct {
    PageCommon
    DirList []DirConf
}

type ListPage struct {
    PageCommon
    VidList []string
}

type PlayPage struct {
    PageCommon
    Vid   string
}

func NewConfigPage() *ConfigPage{
    return &ConfigPage{PageCommon{Title:""}}
}

func NewHomePage() *HomePage{
    return &HomePage{PageCommon:PageCommon{Title:""},DirList:make([]DirConf,10)}
}

func NewListPage() *ListPage{
    return &ListPage{PageCommon:PageCommon{Title:""},VidList:make([]string,10)}
}

func NewPlayPage() *PlayPage{
    return &PlayPage{PageCommon:PageCommon{Title:""},Vid:""}
}

func (p ConfigPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    LogAccess(w,r)
    defer HandlePanic(w,r)
    pageData := NewConfigPage()
    pageData.Title = "配置页"
    drivers.ExecuteTemplate(w, "config",pageData )
}



func (p HomePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    LogAccess(w,r)
    defer HandlePanic(w,r)
    pageData := NewHomePage()
    pageData.Title = "主页"
    pageData.DirList = GetDirConfList()
    drivers.ExecuteTemplate(w, "home", pageData)
}



func (p ListPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    LogAccess(w,r)
    defer HandlePanic(w,r)
    pageData := NewListPage()
    target := r.FormValue("target")
    log.Printf("target:%s \n" , target)
    switch target {
    case "favorite":
        log.Println("step favorite")
        pageData.Title = "收藏页"
        pageData.VidList = GetFavoriteVidList()
        break
    case "hate":
        log.Println("step hate")
        pageData.Title = "黑名单"
        pageData.VidList = GetHateVidList()
        break
    default:
        log.Println("step default")
        pageData.Title = "unkown"
        dirId , _ := strconv.ParseInt(target , 10 , 64)
        dirInfo := GetDirConfById(dirId)
        if dirInfo != nil {
            pageData.Title = dirInfo.Info
            pageData.VidList = dirInfo.Vids
        }
        break
    }
    drivers.ExecuteTemplate(w, "list", pageData)
}



func (p PlayPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    LogAccess(w,r)
    defer HandlePanic(w,r)
    pageData := NewPlayPage()
    pageData.Title = "播放页"
    pageData.Vid = r.FormValue("vid")
    if pageData.Vid == "" {
        panic("valid param")
    }
    drivers.ExecuteTemplate(w, "play", pageData)
}