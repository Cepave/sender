package main

import (
	"flag"
	"fmt"
	"github.com/Cepave/sender/cron"
	"github.com/Cepave/sender/g"
	"github.com/Cepave/sender/http"
	"github.com/Cepave/sender/redis"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	cron.InitWorker()
	redis.InitConnPool()

	go http.Start()
	go cron.ConsumeSms()
	go cron.ConsumeMail()
	go cron.ConsumeQQ()
	go cron.ConsumeServerchan()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		redis.ConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
