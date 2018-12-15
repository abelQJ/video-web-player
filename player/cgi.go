package player

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"errors"
	"strconv"
	"log"
)

type GetConfigCgi struct {

}

type SetConfigCgi struct {

}

type GetVideoCgi struct {

}

type LoveVideoCgi struct {

}

type CancelLoveVideoCgi struct {

}

type HateVideoCgi struct {

}

type CancelHateVideoCgi struct {

}

func AddDelVid(w http.ResponseWriter, r *http.Request , table string, add bool) {
	LogAccess(w,r)
	defer HandlePanic(w,r)
	w.Header().Set("Content-Type" , "text/json")
	vid := r.FormValue("vid")
	if vid == "" {
		panic(errors.New("invalid param"))
	}
	if add {
		AddVid(table,vid)
	} else {
		DelVid(table,vid)
	}
	var data CgiCommonRet
	data.Ret = 0
	data.Msg = "succ"
	w.Write([]byte(data.ToJsonBytes()))
}

func (c LoveVideoCgi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	AddDelVid(w,r,"love" , true)
}

func (c CancelLoveVideoCgi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	AddDelVid(w,r,"love" , false)
}


func (c HateVideoCgi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	AddDelVid(w,r,"hate" , true)
}

func (c CancelHateVideoCgi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	AddDelVid(w,r,"hate" , false)
}



func (c GetConfigCgi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LogAccess(w,r)
	defer HandlePanic(w,r)
	w.Header().Set("Content-Type" , "text/json")
	w.Write([]byte(GetConfig()))
}

func (c SetConfigCgi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LogAccess(w,r)
	defer HandlePanic(w,r)
	w.Header().Set("Content-Type" , "text/json")
	var data CgiCommonRet
	newConfVal , err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(errors.New("param read error"))
	}
	var vConf VDirsConf
	err = json.Unmarshal(newConfVal , &vConf)
	if err != nil {
		panic(errors.New("param json decode error"))
	}
	SetConfig(string(newConfVal))
	data.Ret = 0
	data.Msg = "set confit succ"
	w.Write(data.ToJsonBytes())
}

func (c GetVideoCgi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LogAccess(w,r)
	defer HandlePanic(w,r)
	
	reqRange := r.Header.Get("Range")
	if reqRange == "" {
		reqRange = "bytes=0-"
	}
	log.Printf("http req range:%s\n" , reqRange)
	vid := r.FormValue("vid")
	filePath := GetPathByVid(vid)
	if vid == "" || filePath == "" {
		panic(errors.New("param error"))
	}
	fr := NewFileReader(filePath , reqRange)
	content , err := fr.Read()
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type" , "video/mp4")
	w.Header().Set("Content-Length" , strconv.Itoa(len(content)))
	w.Header().Set("Accept-Ranges" , "bytes")
	w.Header().Set("Connection",  "keep-alive")
	contentRange := fmt.Sprintf("bytes %d-%d/%d",fr.StartPos,fr.EndPos,fr.FileLen)
	w.Header().Set("Content-Range" , contentRange )
	log.Printf("content range:%s\n" , contentRange)
	w.WriteHeader(206)
	w.Write(content)
}