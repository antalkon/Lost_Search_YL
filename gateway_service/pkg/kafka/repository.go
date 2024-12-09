package kafka

import (
	"gateway_service/pkg/syncmap"
	"sync"
)

type BrokerRepo struct {
	consumer *Consumer
	producer *Producer
	mutex    *sync.Mutex
	requests *syncmap.SyncMap
}

func NewBrokerRepo(consumer *Consumer, producer *Producer) *BrokerRepo {
	return &BrokerRepo{consumer: nil, producer: nil, requests: syncmap.NewSyncMap()}
}

func (b *BrokerRepo) SearchAds(uuid string, data string) (chan string, error) {
	ch := make(chan string)
	b.requests.Write(uuid, ch)
	err := b.producer.SearchAds(data)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) MakeAds(uuid string, data string) (chan string, error) {
	return nil, nil
}

func (b *BrokerRepo) ApplyAds(uuid string, data string) (chan string, error) {
	return nil, nil
}

func (b *BrokerRepo) Login(uuid string) error {

}

func (b *BrokerRepo) Register(uuid string) error {}

func (b *BrokerRepo) DeleteChan(uuid string) {
	b.requests.Delete(uuid)
}
