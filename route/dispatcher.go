package route

import (
	"fmt"
	"github.com/benderpan/emitter-relay/config"
	emitter "github.com/emitter-io/go/v2"
)

type Dispatcher struct {
	Cfg        *config.Config
	CurrentHop *emitter.Client
	NextHop    *emitter.Client
}

func (disp *Dispatcher) Start() error {
	err := disp.initCurrentHop()
	if err != nil {
		return err
	}
	err = disp.initNextHop()
	if err != nil {
		return err
	}
	return nil
}

func (disp *Dispatcher) initCurrentHop() error {
	var err error
	disp.CurrentHop, err = emitter.Connect(disp.Cfg.CurrentHop.Host, func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B CurrentHop] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})
	if err != nil {
		return err
	}
	// Subscribe to the presence demo channel
	return disp.CurrentHop.Subscribe(disp.Cfg.CurrentHop.Key, disp.Cfg.CurrentHop.ChannelIn(), disp.currentHopRecvMsgHandler)
}

func (disp *Dispatcher) initNextHop() error {
	var err error
	disp.NextHop, err = emitter.Connect(disp.Cfg.NextHop.Host, func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B NextHop] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})
	if err != nil {
		return err
	}
	// Subscribe to the presence demo channel
	return disp.NextHop.Subscribe(disp.Cfg.NextHop.Key, disp.Cfg.NextHop.ChannelOut(), disp.nextHopRecvMsgHandler)
}

func (disp *Dispatcher) currentHopRecvMsgHandler(client *emitter.Client, msg emitter.Message) {
	fmt.Printf("CurrentHop Received message: %s\n", msg.Payload())
}

func (disp *Dispatcher) nextHopRecvMsgHandler(client *emitter.Client, msg emitter.Message) {
	fmt.Printf("NextHop Received message: %s\n", msg.Payload())
}
func (disp *Dispatcher) CurrentHopSendMsg(msg string) error {
	return disp.currentHopPublish(disp.Cfg.CurrentHop.Key, disp.Cfg.CurrentHop.ChannelOut(), msg)
}

func (disp *Dispatcher) NextHopSendMsg(msg string) error {
	return disp.nextHopPublish(disp.Cfg.NextHop.Key, disp.Cfg.NextHop.ChannelIn(), msg)
}

func (disp *Dispatcher) currentHopPublish(key string, channel string, msg string) error {
	return disp.CurrentHop.Publish(key, channel, msg)
}

func (disp *Dispatcher) nextHopPublish(key string, channel string, msg string) error {
	return disp.NextHop.Publish(key, channel, msg)
}

func (disp *Dispatcher) Close() {
	_ = disp.CurrentHop.Unsubscribe(disp.Cfg.CurrentHop.Key, disp.Cfg.CurrentHop.ChannelIn())
	_ = disp.NextHop.Unsubscribe(disp.Cfg.NextHop.Key, disp.Cfg.NextHop.ChannelOut())
	disp.CurrentHop.Disconnect(0)
	disp.NextHop.Disconnect(0)
}

func NewDispatcher(cfg *config.Config) *Dispatcher {
	return &Dispatcher{Cfg: cfg}
}
