package crons

import (
	"eshop_cart/log"
	"eshop_cart/service"
	"github.com/robfig/cron/v3"
)

func AddCron() {
	c := cron.New(cron.WithSeconds())

	// 每五秒更新一次数据库
	_, err := c.AddFunc("*/5 * * * * *", service.UpdateCart)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	go c.Start()
	defer c.Stop()
}
