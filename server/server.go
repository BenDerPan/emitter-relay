package main

import (
	"flag"
	"fmt"
	"github.com/benderpan/emitter-relay/server/config"
	"github.com/benderpan/emitter-relay/server/route"
	"os"
	"os/signal"
	"time"
)

var c = flag.String("c", "emitter-relay.yaml", "运行参数")

func main() {
	flag.Parse()
	cfg := &config.Config{}
	err := cfg.Load(*c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n%#v\n", cfg.CurrentHop, cfg.NextHop)

	disp := route.NewDispatcher(cfg)
	err = disp.Start()
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	time.Sleep(time.Duration(3) * time.Second)
	disp.Close()

}
