package main

import (
	"github.com/namsral/flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	crontabPath string
)

func init() {
	flag.StringVar(&crontabPath, "file", "crontab", "crontab file path")
}

func main() {
	flag.Parse()

	file, err := os.Open(crontabPath)
	if err != nil {
		log.Fatalf("crontab path:%v err:%v", crontabPath, err)
	}

	parser, err := NewParser(file)
	if err != nil {
		log.Fatalf("Parser read err:%v", err)
	}

	runner, err := parser.Parse()
	if err != nil {
		log.Fatalf("Parser parse err:%v", err)
	}

	file.Close()

	var wg sync.WaitGroup
	shutdown(runner, &wg)

	runner.Start()
	wg.Add(1)

	wg.Wait()
	log.Println("End cron")
}

func shutdown(runner *Runner, wg *sync.WaitGroup) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Println("Got signal: ", s)
		runner.Stop()
		wg.Done()
	}()
}
