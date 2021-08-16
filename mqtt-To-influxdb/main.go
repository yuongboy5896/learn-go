package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"example.com/m/tool"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

//var c chan os.Signal

//
var quit = false
var influxdb influxdb2.Client
var logs = true
var IOTPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if logs {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

	Payloadstr := string(msg.Payload())
	Topicstr := string(msg.Topic())
	insertInfluxdb(Payloadstr, Topicstr)
}

func insertInfluxdb(Payloadstr string, Topicstr string) {
	//Payloadstr := string(Payload)

	if !strings.Contains(Payloadstr, "post") {
		//fmt.Println("数据不正确")
		return
	}
	topic := strings.Split(Topicstr, "/")
	var deviceName string
	if len(topic) > 6 {
		deviceName = topic[3]
	} else {
		println("没有获取设备名称")
		//return
	}
	//fmt.Printf("设备名称 %s\n", deviceName)
	string_slice := strings.Split(Payloadstr, "==========")
	var data map[string]interface{}
	if len(string_slice) > 1 {
		relust := strings.ReplaceAll(string_slice[1], "=", "")
		if err := json.Unmarshal([]byte(relust), &data); err != nil {
			print("参数解析失败 %s \n", err)
			//return
		}
	}

	const bucket = "IOT"
	const org = "devops"
	tagmap := make(map[string]string)
	tagmap["IP"] = string_slice[0]

	params := data["params"].(map[string]interface{})
	if params != nil {
		writeAPI := influxdb.WriteAPI(org, bucket)
		p := influxdb2.NewPoint(deviceName,
			tagmap,
			params,
			time.Now())
		// write point asynchronously
		writeAPI.WritePoint(p)
		writeAPI.Flush()
	}

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	if err != nil {
		quit = true
		fmt.Printf("Connect lost: %v", err)
	}
}

func main() {
	timeout, _ := time.ParseDuration("10s")
	// 初始化配置
	cfg, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		println(err)
	}
	logs = cfg.Mqtt.Bugger
	influxdb = influxdb2.NewClient(cfg.Infludb.Infurl, cfg.Infludb.Token)
	defer influxdb.Close()

	var broker = cfg.Mqtt.Mqttip
	var port = cfg.Mqtt.Mqttport
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(cfg.Mqtt.Clientid)
	opts.SetUsername(cfg.Mqtt.MqttUname)
	opts.SetPassword(cfg.Mqtt.MqttPwd)
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetDefaultPublishHandler(IOTPubHandler)
	opts.OnConnect = connectHandler
	opts.ConnectTimeout = 3 * time.Minute
	opts.OnConnectionLost = connectLostHandler
	opts.SetAutoReconnect(true).SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectRetry(true)
	opts.SetConnectTimeout(timeout)
	opts.SetPingTimeout(6 * timeout)
	opts.SetKeepAlive(3 * timeout)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	filters := make(map[string]byte)
	//filters["/sys/#"] = 1
	filters["/sys/+/+/thing/event/property/post_reply"] = 1

	subIotMultiple(client, filters)
	for {
		if quit {
			client.Disconnect(250)
			fmt.Println("程序退出")
			os.Exit(3)
		}
		time.Sleep(2 * time.Second)
	}
}

func subIotMultiple(client mqtt.Client, filters map[string]byte) {
	if token := client.SubscribeMultiple(filters, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
