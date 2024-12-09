package kafka

import "sync"

type BrokerRepo struct {
	consumer *Consumer
	producer *Producer
	mutex    *sync.Mutex
	requests map[string]chan string
}

func NewBrokerRepo(consumer *Consumer, producer *Producer) *BrokerRepo {
	return &BrokerRepo{consumer: nil, producer: nil,
		mutex: &sync.Mutex{}, requests: make(map[string]chan string)}
}

func (b *BrokerRepo) SearchAds(uuid string, data string) (chan string, error) {
	ch := make(chan string)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.requests[uuid] = ch
	err := b.producer.SearchAds(data)
	if err != nil {
		return nil, err
	}
	return ch, nil
}
