package cmd

import (
	"fmt"
	"zoho-center/core/config"
	"zoho-center/core/database"
	"zoho-center/core/event"
	"zoho-center/core/log"
	"zoho-center/core/router"
	"zoho-center/job/v1/purchaseorder"
	"zoho-center/job/v1/salesorder"
)

func Run(args []string) {
	config.LoadConfig(args[1])
	log.ConfigLogger()
	// cache.ConfigCache()
	event.Subscribe(salesorder.Subscribe)
	database.ConfigMysql()
	runType := config.ReadConfig("application.type")
	fmt.Println(runType)
	if runType == "api" {
		r := router.InitRouter()
		// router.InitPublicRouter(r, auth.Routers)
		// router.InitAuthRouter(r, organization.Routers, project.Routers, event.Routers, component.Routers, auth.AuthRouter, client.Routers, position.Routers)
		router.RunServer(r)
	} else if runType == "job" {
		purchaseorder.GetPurchaseorderList()
	} else {
		fmt.Println("type error")
	}

}
