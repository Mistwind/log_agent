package main

import (
	"fmt"
	"strings"
	"time"

	"./conf"
	"./kafka"
	"./taillog"

	"gopkg.in/ini.v1"
)

var config = new(conf.Config)

func main() {
	err := ini.MapTo(config, "./conf/config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		return
	}
	fmt.Println(config)

	err = kafka.Init(strings.Split(config.Kafka.Address, ";"))
	if err != nil{
		fmt.Printf("Init kafka failed, err: %v\n", err)
		return
	}
	fmt.Println("Init kafka successful!")

	err = taillog.Init(config.Log.FileName)
	if err != nil {
		fmt.Printf("Init tailLog failed, err: %v\n", err)
		return
	}
	fmt.Println("Init taillog successful!")

	run()
}

func run() {
	for {
		select {
		case line := <-taillog.ReadChan():
			kafka.SendToKafka(config.Kafka.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}
