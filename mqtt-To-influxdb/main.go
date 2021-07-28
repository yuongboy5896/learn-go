package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"example.com/m/dao"
	"example.com/m/model"
	"example.com/m/tool"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

//var c chan os.Signal

//
var quit = false
var influxdb influxdb2.Client
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	r := ioutil.NopCloser(bytes.NewReader([]byte(msg.Payload())))
	var mqttMsg model.MqttMsg

	err := tool.Decode(r, &mqttMsg)
	if err != nil {
		print("参数解析失败 %s \n", err)
		return
	}
	fmt.Printf("RH %.2f,TMP: %.2f \n", mqttMsg.RH, mqttMsg.TMP)
	const bucket = "test"
	const org = "devops"

	filemap := make(map[string]interface{})
	filemap["rh"] = mqttMsg.RH
	filemap["tmp"] = mqttMsg.TMP
	writeAPI := influxdb.WriteAPI(org, bucket)
	p := influxdb2.NewPoint("home",
		map[string]string{"unit": "temperature"},
		filemap,
		time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	writeAPI.Flush()

}

//

var IOTPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	Payloadstr := string(msg.Payload())
	topic := strings.Split(msg.Topic(), "/")
	var deviceName string
	if len(topic) > 3 {
		deviceName = topic[3]
	} else {
		println("没有获取设备名称")
		return
	}
	fmt.Printf("设备名称 %s\n", deviceName)
	//string_slice := strings.Split(Payloadstr, "===========")
	string_slice := strings.Split(Payloadstr, "==========")
	var data map[string]interface{}
	if len(string_slice) > 1 {
		relust := strings.ReplaceAll(string_slice[1], "=", "")
		if err := json.Unmarshal([]byte(relust), &data); err != nil {
			print("参数解析失败 %s \n", err)
			return
		}
	}

	const bucket = "IOT"
	const org = "devops"
	tagmap := make(map[string]string)
	tagmap["IP"] = string_slice[0]

	params := data["params"].(map[string]interface{})
	if params == nil {
		println("转行失败")
		return
	}
	writeAPI := influxdb.WriteAPI(org, bucket)
	p := influxdb2.NewPoint(deviceName,
		tagmap,
		params,
		time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	writeAPI.Flush()

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	//fmt.Printf("Connect lost: %v", err)
	if err != nil {
		quit = true
		fmt.Printf("Connect lost: %v", err)
	}
}

func main() {

	// 初始化配置
	cfg, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		println(err)
	}
	//
	_, err = tool.OrmEngine(cfg)
	if err != nil {
		println(err)
		return
	}
	deviceIoTDao := dao.NewDeviceIoTDao()
	deviceTops, err := deviceIoTDao.QuerydeviceIotsByType()
	if err != nil {
		println(err)
		return
	}
	fmt.Println(len(deviceTops))
	//flag.Parse()
	//fmt.Println(os.Args)

	//infurl := flag.String("infurl", "http://192.168.2.60:8086", "http url")
	//fmt.Println(infurl)
	//mqttip := flag.String("mqttip", "192.168.48.100", "mqttip")
	//mqttport := flag.Int("mqttport", 1883, "mqttport")
	//username := flag.String("username", "thpower", "username")
	//password := flag.String("password", "Thp@IOT12345678", "password")
	//clientid := flag.String("clientid", "123345455", "clientid")

	//var stoken string
	//fmt.Printf("mqttip %s \n", *mqttip)
	//flag.StringVar(&stoken, "token", "zwvS0JXTQU2LUiEnWCmLjr6mq_E1UPJagrpePLalFO-SvsmVxKFoC-f1oDZDTU_PTuIGKiVuseFQIn2OR9YFvw==", "token")

	//influxdb2
	//const token = stoken
	//const bucket = "devops"
	//const org = "devops"
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
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	filters := make(map[string]byte)
	for i := 0; i < len(deviceTops); i++ {
		filters[deviceTops[i].TopicName] = 1
		//subIot(client, deviceTops[i].TopicName)
	}
	//sub(client)
	subIotMultiple(client, filters)
	for {
		if quit {
			client.Disconnect(250)
			os.Exit(3)
		}
	}

	//client.Disconnect(250)
}

/*
func sub(client mqtt.Client) {
	topic := "test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
*/
func subIotMultiple(client mqtt.Client, filters map[string]byte) {
	token := client.SubscribeMultiple(filters, nil)
	token.Wait()
}

//func subIot(client mqtt.Client, topic string) {
//	token := client.SubscribeMultiple(topic, 1, nil)
//token.Wait()
//fmt.Printf("Subscribed to topic: %s \n", topic)
//}
