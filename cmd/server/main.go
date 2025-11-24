package main

import (
	"gibhub.com/l-jessie/test-im/internal/api"
)

func main() {
	// Init 路由
	router := api.Router()

	// run server
	if err := router.Run(":8070"); err != nil {
		panic(err)
	}
}
