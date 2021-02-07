package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var dir = `D:\workspace\golang\project\go-fftool\video_split_temp`

func main() {
	router := gin.Default()

	router.Use(cors.Default())
	router.Static("/", dir)

	if err := router.Run("0.0.0.0:8081"); err != nil {
		panic(err)
	}

}
