package route

import (
	"fmt"
	"github.com/benderpan/emitter-relay/server/config"
	emitter "github.com/emitter-io/go/v2"
)

type RelayDispatcher struct {
	Cfg        *config.Config
	CurrentHop *emitter.Client
	NextHop    *emitter.Client
}

func (disp *RelayDispatcher) Start() error {
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

func (disp *RelayDispatcher) initCurrentHop() error {
	var err error
	disp.CurrentHop, err = emitter.Connect(disp.Cfg.CurrentHop.Host, func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B CurrentHop] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})
	if err != nil {
		return err
	}
	// Subscribe to the presence demo channel
	return disp.CurrentHop.Subscribe(disp.Cfg.CurrentHop.KeyIn, disp.Cfg.CurrentHop.ChannelIn, disp.currentHopRecvMsgHandler)
}

func (disp *RelayDispatcher) initNextHop() error {
	var err error
	disp.NextHop, err = emitter.Connect(disp.Cfg.NextHop.Host, func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B NextHop] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})
	if err != nil {
		return err
	}
	// Subscribe to the presence demo channel
	return disp.NextHop.Subscribe(disp.Cfg.NextHop.KeyIn, disp.Cfg.NextHop.ChannelIn, disp.nextHopRecvMsgHandler)
}

func (disp *RelayDispatcher) currentHopRecvMsgHandler(client *emitter.Client, msg emitter.Message) {
	fmt.Printf("CurrentHop Received message: %s, Topic=%s\n", msg.Payload(), msg.Topic())
	err := disp.NextHopSendMsg(msg.Payload())
	if err != nil {
		fmt.Printf("转发消息到下一跳失败:%v", err)
	}
}

func (disp *RelayDispatcher) nextHopRecvMsgHandler(client *emitter.Client, msg emitter.Message) {
	fmt.Printf("NextHop Received message: %s, Topic=%s\n", msg.Payload(), msg.Topic())
	err := disp.CurrentHopSendMsg(msg.Payload())
	if err != nil {
		fmt.Printf("回送消息到上一跳失败:%v", err)
	}
}
func (disp *RelayDispatcher) CurrentHopSendMsg(msg interface{}) error {
	return disp.CurrentHop.Publish(disp.Cfg.CurrentHop.KeyOut, disp.Cfg.CurrentHop.ChannelOut, msg)
}

func (disp *RelayDispatcher) NextHopSendMsg(msg interface{}) error {
	return disp.NextHop.Publish(disp.Cfg.NextHop.KeyOut, disp.Cfg.NextHop.ChannelOut, msg)
}

func (disp *RelayDispatcher) Close() {
	_ = disp.CurrentHop.Unsubscribe(disp.Cfg.CurrentHop.KeyIn, disp.Cfg.CurrentHop.ChannelIn)
	_ = disp.NextHop.Unsubscribe(disp.Cfg.NextHop.KeyIn, disp.Cfg.NextHop.ChannelIn)
	disp.CurrentHop.Disconnect(0)
	disp.NextHop.Disconnect(0)
}

func NewDispatcher(cfg *config.Config) *RelayDispatcher {
	return &RelayDispatcher{Cfg: cfg}
}
