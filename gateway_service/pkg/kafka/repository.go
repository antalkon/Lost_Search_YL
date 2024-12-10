package kafka

import (
	"context"
	"gateway_service/pkg/syncmap"
)

type BrokerRepo struct {
	consumer map[string]*Consumer
	producer map[string]*Producer
	requests *syncmap.SyncMap
	ctx      context.Context
	stop     context.CancelFunc
}

func NewBrokerRepo(address string) *BrokerRepo {
	ctx, cancel := context.WithCancel(context.Background())
	repo := &BrokerRepo{requests: syncmap.NewSyncMap(), ctx: ctx, stop: cancel}
	var producers = make(map[string]*Producer)
	var err error
	producers["Notify"], err = NewProducer(address, "NotifyRequest")
	if err != nil {
		panic(err)
	}
	producers["Ads"], err = NewProducer(address, "AdsRequest")
	if err != nil {
		panic(err)
	}
	producers["Auth"], err = NewProducer(address, "AuthRequest")
	if err != nil {
		panic(err)
	}
	repo.producer = producers
	var consumers = make(map[string]*Consumer)
	consumers["Notify"] = NewConsumer(address, "Notify", "NotifyResponse", repo.requests) //make constants instead "Notify" ...etc
	consumers["Ads"] = NewConsumer(address, "Ads", "AdsResponse", repo.requests)
	consumers["Auth"] = NewConsumer(address, "Auth", "AuthResponse", repo.requests)
	repo.consumer = consumers
	for _, i := range repo.consumer {
		go func() {
			if err = i.Consume(repo.ctx); err != nil {
				panic(err)
			}
		}()
	}
	return repo
}

func (b *BrokerRepo) SearchAds(uuid string, data string) (chan string, error) {
	ch := make(chan string)
	b.requests.Write(uuid, ch)
	err := b.producer["Ads"].SendMessage(b.ctx, uuid, data)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) MakeAds(uuid string, data string) (chan string, error) {
	ch := make(chan string)
	b.requests.Write(uuid, ch)
	err := b.producer["Ads"].SendMessage(b.ctx, uuid, data)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) ApplyAds(uuid string, data string) (chan string, error) {
	ch := make(chan string)
	b.requests.Write(uuid, ch)
	err := b.producer["Ads"].SendMessage(b.ctx, uuid, data)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) Login(uuid string, data string) (chan string, error) {
	ch := make(chan string)
	b.requests.Write(uuid, ch)
	err := b.producer["Auth"].SendMessage(b.ctx, uuid, data)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) Register(uuid string, data string) (chan string, error) {
	ch := make(chan string)
	b.requests.Write(uuid, ch)
	err := b.producer["Ads"].SendMessage(b.ctx, uuid, data)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) NotifyUser(uuid string, data string) (chan string, error) {
	ch := make(chan string)
	b.requests.Write(uuid, ch)
	err := b.producer["Auth"].SendMessage(b.ctx, uuid, data)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) DeleteChan(uuid string) {
	b.requests.Delete(uuid)
}

func (b *BrokerRepo) Stop() {
	b.stop()
}
