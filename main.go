package main

import (
	"eshop_cart/database"
	"eshop_cart/kitex_gen/eshop/cart/cartservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {
	database.Init()
	svr := cartservice.NewServer(
		new(CartServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "CartService"}),
	)
	err := svr.Run()
	if err != nil {
		panic(err)
	}
}
