package mq

import (
	"encoding/json"
	"eshop_cart/log"
	"eshop_cart/service"
)

func ConsumeSaveCartMessage() {
	ch, err := Conn.Channel()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	defer ch.Close()
	// 消费消息

	msgs, err := ch.Consume(
		"save_cart_queue",          // 队列名称
		"save_cart_queue_consumer", // 消费者标签
		false,                      // 是否自动确认消息
		false,                      // 是否独占消费者（仅限于本连接）
		false,                      // 是否阻塞等待服务器确认
		false,                      // 是否使用内部排他队列
		nil,                        // 其他参数
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	log.Infof("save_cart_queue_consumer消费者启动")
	for msg := range msgs {
		log.Infof("save_cart_queue_consumer消费者收到消息: %s", string(msg.Body))
		// 更新库
		message := make(map[string]interface{})
		err = json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Errorf("err: %v", err)
			return
		}
		uid := message["uid"].(string)
		service.UpdateCart(uid)
		if err != nil {
			log.Errorf("err: %v", err)
			return
		}
		// 手动确认消息已被消费
		msg.Ack(false)
	}
}
