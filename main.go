package main

import (
	"net/http"
	"./player"
	"log"
	"path/filepath"
)

func main(){
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
	err := http.ListenAndServe(":8105", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}