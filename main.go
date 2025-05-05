package main

import (
	"eshop_cart/database"
	"eshop_cart/kitex_gen/eshop/cart/cartservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	zkregistry "github.com/kitex-contrib/registry-zookeeper/registry"
	"time"
)

func main() {
	database.Init()
	r, err := zkregistry.NewZookeeperRegistry([]string{"117.72.72.114:2181"}, 40*time.Second)
	if err != nil {
		panic(err)
	}
	svr := cartservice.NewServer(new(CartServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "eshop_cart"}))
	//svr := cartservice.NewServer(
	//	new(CartServiceImpl),
	//	server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "CartService"}),
	//)
	err = svr.Run()
	if err != nil {
		panic(err)
	}
}
