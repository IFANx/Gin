package main

import (
	"Gin/common"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	db := common.InitDB()
	defer fmt.Println(db)
	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run()) // listen and serve on 0.0.0.0:8080
}
