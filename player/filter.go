package player

import (
	"net/http"
	"log"
	"encoding/json"
	"runtime/debug"
	"fmt"
)

type CgiCommonRet struct {
	Ret         int64      `json:"ret"`
	Msg         string     `json:"msg"`
}

func (ret CgiCommonRet) ToJsonBytes() []byte {
	jBytes , _ := json.Marshal(ret)
	return jBytes
}

func HandlePanic(w http.ResponseWriter, r *http.Request){
	if err := recover(); err != nil {
		log.Printf("error happen:%v \n" , string(debug.Stack()))
		var data CgiCommonRet
		data.Ret = -99
		data.Msg = fmt.Sprintf("unkown error happen:%v" , err)
		w.Write(data.ToJsonBytes())
	}
}

func LogAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("access|Method:%s,Host:%s,RemoteAddr:%s,RequestURI:%s\n" ,
	         r.Method , r.Host, r.RemoteAddr , r.RequestURI)
}