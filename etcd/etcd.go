package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var cli *clientv3.Client

// LogEntry 日志配置信息的结构
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// Init 初始化 etcd 配置
func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(addr, ";"),
		DialTimeout: timeout,
	})
	return
}

// GetConf 从 etcd 中获取配置
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("Get from etcd failed, err: %v\n", err)
		return
	}

	for _, ev := range resp.Kvs {
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("Unmarshal etcd value failed, err: %v\n", err)
			return
		}
	}
	return
}

// WatchConf 监视配置是否发生变化
func WatchConf(key string, newConfChan chan<- []*LogEntry) {
	rch := cli.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s  Key: %s  Value: %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			var newConf []*LogEntry
			if ev.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("Unmarshal failed, err: %v\n", err)
					continue
				}
			}

			fmt.Printf("Get new conf: %v\n", newConf)
			newConfChan <- newConf
		}
	}
}
