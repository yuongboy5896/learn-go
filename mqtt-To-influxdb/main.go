package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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

	fmt.Println(os.Args)

	infurl := flag.String("infurl", "http://192.168.2.60:8086", "http url")
	fmt.Println(infurl)
	mqttip := flag.String("mqttip", "192.168.48.100", "mqttip")
	mqttport := flag.Int("mqttport", 1883, "mqttport")
	username := flag.String("username", "thpower", "username")
	password := flag.String("password", "Thp@IOT12345678", "password")
	clientid := flag.String("clientid", "123345455", "username")
	var stoken string

	flag.StringVar(&stoken, "token", "zwvS0JXTQU2LUiEnWCmLjr6mq_E1UPJagrpePLalFO-SvsmVxKFoC-f1oDZDTU_PTuIGKiVuseFQIn2OR9YFvw==", "token")
	flag.Parse()
	//influxdb2
	//const token = stoken
	//const bucket = "devops"
	//const org = "devops"
	influxdb = influxdb2.NewClient(*infurl, stoken)
	defer influxdb.Close()

	var broker = mqttip
	var port = mqttport
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", *broker, *port))
	opts.SetClientID(*clientid)
	opts.SetUsername(*username)
	opts.SetPassword(*password)
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetDefaultPublishHandler(IOTPubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	for i := 0; i < len(deviceTops); i++ {

	}
	sub(client)

	for {
		if quit {
			client.Disconnect(250)
			os.Exit(3)
		}
	}

	//client.Disconnect(250)
}

func sub(client mqtt.Client) {
	topic := "test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
