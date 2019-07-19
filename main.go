package main

import (
	"fmt"
	"github.com/benderpan/emitter-relay/config"
	"github.com/benderpan/emitter-relay/route"
	"time"
)

func main() {
	cfg := &config.Config{}
	err := cfg.Load("emitter-relay.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n%#v\n", cfg.CurrentHop, cfg.NextHop)

	disp := route.NewDispatcher(cfg)
	disp.Start()
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		disp.CurrentHopSendMsg("你好，我是当前节点")
		disp.NextHopSendMsg("你好，我是下一跳节点")
	}
	time.Sleep(time.Duration(3) * time.Second)
	disp.Close()

}
