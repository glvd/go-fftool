package main

import (
	"net/http"
)

var dir = `D:\workspace\golang\project\go-fftool\video_split_temp`

func main() {
	err := http.ListenAndServe(":8081", http.FileServer(http.Dir(dir)))
	if err != nil {
		panic(err)
	}
}
