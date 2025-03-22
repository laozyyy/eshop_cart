package crons

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"big_market/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func scanTask() {
	taskEntities, err := database.QueryNoSendMessageTaskList(database.DB)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	if taskEntities == nil {
		log.Infof("无待发送mq")
		return
	}
	log.Infof("有待发送mq %d 条", len(taskEntities))
	for _, task := range taskEntities {
		log.Infof("当前发送mq:%+v", task)
		task := task
		go func() {
			publish := amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(*task.Message),
			}
			err = mq.SendMessage(constant.NormalExchangeName, task.Topic, publish)
			if err != nil {
				log.Errorf("err: %v", err)
				return
			}
			err = database.UpdateTaskCompleted(database.DB, *task)
			if err != nil {
				log.Errorf("err: %v", err)
				return
			}
			log.Infof("更新taskID: %v 为completed", task.ID)
		}()
	}
}
