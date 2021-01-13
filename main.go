package main

import (
	"fmt"
	"log_collect/tools"
	"strings"
	"sync"
	"time"

	"log_collect/conf"
	"log_collect/etcd"
	"log_collect/kafka"
	"log_collect/taillog"

	"gopkg.in/ini.v1"
)

var config = new(conf.Config)

func main() {
	// 0. 加载配置文件
	err := ini.MapTo(config, "./conf/config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		return
	}
	fmt.Println(config)

	// 1. 初始化 kafka 连接
	err = kafka.Init(strings.Split(config.Kafka.Address, ";"), config.Kafka.ChanMaxSize)
	if err != nil {
		fmt.Printf("Init kafka failed, err: %v\n", err)
		return
	}
	fmt.Println("Init kafka successful!")

	// 2. 初始化 etcd
	err = etcd.Init(config.Etcd.Address, time.Duration(config.Etcd.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("Init etcd failed, err: %v\n", err)
		return
	}
	fmt.Printf("Init etcd success\n")

	// 实现每个 logagent 都拉取自己独有的配置，所以要以自己的ip地址实现热加载
	ip, err := tools.GetOurboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(config.Etcd.Key, ip)

	// 2.1 从 etcd 中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(config.Etcd.Key)
	if err != nil {
		fmt.Printf("etcd.GetConf failed, err: %v\n", err)
		return
	}
	fmt.Printf("Get conf from etcd success, %v\n", logEntryConf)

	for index, value := range logEntryConf {
		fmt.Printf("index: %v  value: %v\n", index, value)
	}

	// 3. 收集日志发往 kafka
	taillog.Init(logEntryConf)

	var wg sync.WaitGroup
	wg.Add(1)
	// 哨兵发现最新的配置信息会通知上面的通道
	go etcd.WatchConf(etcdConfKey, taillog.NewConfChan())
	wg.Wait()

	// err = taillog.Init(config.Log.FileName)
	// if err != nil {
	// 	fmt.Printf("Init tailLog failed, err: %v\n", err)
	// 	return
	// }
	// fmt.Println("Init taillog successful!")

	// run()
}

// func run() {
// 	for {
// 		select {
// 		case line := <-taillog.ReadChan():
// 			kafka.SendToKafka(config.Kafka.Topic, line.Text)
// 		default:
// 			time.Sleep(time.Second)
// 		}
// 	}
// }
