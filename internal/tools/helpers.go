package tools

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

// SetValue устанавливает значение заголовка "ключ-значение"
func setValue(msg *kafka.Message, key string, val []byte) {
	//c[key] = val
	msg.Headers = append(msg.Headers, kafka.Header{Key: key, Value: val})
}

// NewMessageCarrier оборачивает сообщение кафки в структуру - необходимо для передачи контекста
func newMessageCarrier(msg *kafka.Message) *MessageCarrier {
	return &MessageCarrier{msg: msg}
}

// Get - переопределение метода Get для поддержки интерфейса TextMapCarrier
func (c MessageCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

// Set - переопределение метода Get для поддержки интерфейса TextMapCarrier
func (c MessageCarrier) Set(key string, value string) {
	// Ensure uniqueness of keys
	for i := 0; i < len(c.msg.Headers); i++ {
		if c.msg.Headers[i].Key == key {
			c.msg.Headers = append(c.msg.Headers[:i], c.msg.Headers[i+1:]...)
			i--
		}
	}
	c.msg.Headers = append(c.msg.Headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

// Keys - переопределение метода Get для поддержки интерфейса TextMapCarrier
func (c MessageCarrier) Keys() []string {
	out := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		out[i] = h.Key
	}
	return out
}
