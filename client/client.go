package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/benderpan/emitter-relay/client/config"
	emitter "github.com/emitter-io/go/v2"
	"os"
)

type RelayClient struct {
	EmitterClient *emitter.Client
	Cfg           *config.Config
}

func (c *RelayClient) Start() error {
	var err error
	c.EmitterClient, err = emitter.Connect(c.Cfg.ServerNode.Host, func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B CurrentHop] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})
	if err != nil {
		return err
	}
	// Subscribe to the presence demo channel
	return c.EmitterClient.Subscribe(c.Cfg.ServerNode.KeyIn, c.Cfg.ServerNode.ChannelIn, func(client *emitter.Client, msg emitter.Message) {
		fmt.Printf("Client收到任务: '%s'， Channel: '%s'\n", msg.Payload(), msg.Topic())
		//err:=c.EmitterClient.Publish(c.Cfg.ServerNode.KeyOut,c.Cfg.ServerNode.ChannelOut,fmt.Sprintf("应答:%s",msg.Payload()))
		//fmt.Printf("应答发送失败:%v\n",err)
	})

}

func (c *RelayClient) SendMsg(msg interface{}) error {
	return c.EmitterClient.Publish(c.Cfg.ServerNode.KeyOut, c.Cfg.ServerNode.ChannelOut, msg)
}

func (c *RelayClient) Close() {
	_ = c.EmitterClient.Unsubscribe(c.Cfg.ServerNode.KeyIn, c.Cfg.ServerNode.ChannelIn)
	c.EmitterClient.Disconnect(3)
}

func NewRelayClient(cfg *config.Config) *RelayClient {
	client := &RelayClient{}
	client.Cfg = cfg
	return client
}

var c = flag.String("c", "emitter-relay-client.yaml", "运行参数")

func main() {
	flag.Parse()
	cfg := &config.Config{}
	err := cfg.Load(*c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", cfg.ServerNode)
	client := NewRelayClient(cfg)
	err = client.Start()
	if err != nil {
		panic(err)
	}
	input := bufio.NewScanner(os.Stdin)
	fmt.Printf(">")
	for input.Scan() { //扫描输入内容
		line := input.Text() //把输入内容转换为字符串
		if line == "exit" {
			break
		}
		_ = client.SendMsg(line)
		fmt.Printf(">")

	}
	client.Close()

}
