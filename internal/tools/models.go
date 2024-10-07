package tools

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type MessageCarrier struct {
	msg *kafka.Message
}
