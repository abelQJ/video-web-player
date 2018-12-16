package main

import (
    "net/http"
    "./player"
    "log"
    "path/filepath"
    "flag"
	"expvar"
	"fmt"
)

func init(){
	confAccessKey := expvar.NewString("pi")
	expvar.Publish("confAccessKey" , confAccessKey)
}

func main(){
    var confAccessKey,port,https string
	flag.StringVar(&confAccessKey , "confAccessKey" , "pi" ,"-confAccessKey key")
	flag.StringVar(&port , "port" , "8080" ,"-port num")
	flag.StringVar(&https , "https" , "false" ,"-https true|false")
    flag.Parse()
    expvar.Get("confAccessKey").(*expvar.String).Set(confAccessKey)
    http.Handle("/", player.NotFoundPage{})
    http.Handle("/page/list", player.NewListPage())
    http.Handle("/page/config", player.NewConfigPage())
    http.Handle("/page/home" , player.NewHomePage())
    http.Handle("/page/play" , player.NewPlayPage())
    http.Handle("/cgi/setconfig" , player.SetConfigCgi{})
    http.Handle("/cgi/getconfig" , player.GetConfigCgi{})
    http.Handle("/cgi/getVideo" , player.GetVideoCgi{})
    http.Handle("/cgi/loveVideo" , player.LoveVideoCgi{})
    http.Handle("/cgi/cancelLoveVideo" , player.CancelLoveVideoCgi{})
    http.Handle("/cgi/hateVideo" , player.HateVideoCgi{})
    http.Handle("/cgi/cancelHateVideo" , player.CancelHateVideoCgi{})
	http.Handle("/resources/" , http.FileServer(http.Dir(filepath.Join(".","static"))))
	listenPort := fmt.Sprintf(":%s",port)
	var err error
	if https == "true" {
	  serverCrt := filepath.Join(".","ca","server.crt")
	  serverKey := filepath.Join(".","ca","server.key")
      err = http.ListenAndServeTLS(listenPort, serverCrt , serverKey, nil)
	} else {
		err = http.ListenAndServe(listenPort, nil)
	}
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
 }
}