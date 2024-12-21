package kafka

import (
	"context"
	"encoding/json"
	"gateway_service/internal/models"
	"gateway_service/pkg/logger"
	"gateway_service/pkg/syncmap"
)

type BrokerRepo struct {
	consumer map[string]*Consumer
	producer map[string]*Producer
	requests *syncmap.SyncMap
	ctx      context.Context
	stop     context.CancelFunc
}

func NewBrokerRepo(c context.Context, address string) *BrokerRepo {
	ctx, cancel := context.WithCancel(c)
	lg := logger.GetLogger(ctx)
	repo := &BrokerRepo{requests: syncmap.NewSyncMap(), ctx: ctx, stop: cancel}
	var producers = make(map[string]*Producer)
	var err error
	producers["Notify"], err = NewProducer(address, "notify_request")
	if err != nil {
		lg.Error(ctx, err.Error())
	}
	producers["Ads"], err = NewProducer(address, "ads_request")
	if err != nil {
		lg.Error(ctx, err.Error())
	}
	producers["Auth"], err = NewProducer(address, "auth_request")
	if err != nil {
		lg.Error(ctx, err.Error())
	}
	repo.producer = producers
	var consumers = make(map[string]*Consumer)
	consumers["Notify"] = NewConsumer(address, "notify_response", repo.requests) //make constants instead "Notify" ...etc
	consumers["Ads"] = NewConsumer(address, "ads_response", repo.requests)
	consumers["Auth"] = NewConsumer(address, "auth_response", repo.requests)
	repo.consumer = consumers
	for _, i := range repo.consumer {
		go func() {
			if err = i.Consume(repo.ctx); err != nil {
				lg.Error(ctx, err.Error())
			}
		}()
	}
	return repo
}

func (b *BrokerRepo) SearchAds(uuid string, name string, typ string, location models.Location) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.SearchRequest{
		Name:          name,
		TypeOfFinding: typ,
		Location:      location,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "Ads",
		Action:    "get",
		Data:      string(data),
	}
	data, err = json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Ads"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) MakeAds(login, uuid, name, description, typ string, geo models.Location) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.MakeAdsKafkaRequest{
		Login:         login,
		Name:          name,
		Description:   description,
		TypeOfFinding: typ,
		Location:      geo,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "Ads",
		Action:    "add",
		Data:      string(data),
	}
	data, err = json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Ads"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) ApplyAds(uuid string, findUuid string) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.ApplyKafkaRequest{
		Uuid: findUuid,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "Ads",
		Action:    "respond",
		Data:      string(data),
	}
	data, err = json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Ads"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) Register(uuid string, login, password, email string) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.RegisterRequest{
		Login:    login,
		Password: password,
		Email:    email,
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "auth",
		Action:    "create_user",
		Data:      req,
	}
	data, err := json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Auth"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) Login(uuid string, login, password string) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.LoginRequest{
		Login:    login,
		Password: password,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "auth",
		Action:    "login_user",
		Data:      string(data),
	}
	data, err = json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Auth"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) ValidateToken(uuid string, token string) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.ValidateTokenRequest{
		Token: token,
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "auth",
		Action:    "validate_token",
		Data:      req,
	}
	data, err := json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Auth"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) NotifyUser(uuid string, email, subject, body string) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.NotifyRequest{
		Email:   email,
		Subject: subject,
		Body:    body,
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "Notify",
		Action:    "send_email",
		Data:      req,
	}
	data, err := json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Notify"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) GetLogin(uuid, token string) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.GetLoginRequest{
		Token: token,
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "auth",
		Action:    "get_login_by_token",
		Data:      req,
	}
	data, err := json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Auth"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (b *BrokerRepo) GetEmail(uuid, login string) (chan []byte, error) {
	ch := make(chan []byte)
	b.requests.Write(uuid, ch)
	req := models.GetEmailRequest{
		Login: login,
	}
	kafkaReq := models.KafkaRequest{
		RequestId: uuid,
		Service:   "auth",
		Action:    "get_user_data",
		Data:      req,
	}
	data, err := json.Marshal(kafkaReq)
	if err != nil {
		return nil, err
	}
	err = b.producer["Auth"].SendMessage(b.ctx, uuid, string(data))
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (b *BrokerRepo) Stop() {
	b.stop()
}
