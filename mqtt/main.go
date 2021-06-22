package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"example.com/m/model"
	"example.com/m/tool"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var c chan os.Signal

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

/*
func Producer() {
LOOP:
	for {
		select {
		case s := <-c:
			fmt.Println()
			fmt.Println("Producer | get", s)
			break LOOP
		default:
		}
		time.Sleep(500 * time.Millisecond)
	}

}
*/

func main() {

	fmt.Println(os.Args)

	//influxdb2
	const token = "zwvS0JXTQU2LUiEnWCmLjr6mq_E1UPJagrpePLalFO-SvsmVxKFoC-f1oDZDTU_PTuIGKiVuseFQIn2OR9YFvw=="
	//const bucket = "devops"
	//const org = "devops"
	influxdb = influxdb2.NewClient("http://192.168.2.60:8086", token)
	defer influxdb.Close()

	var broker = "yang5896336.tpddns.cn"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("yang")
	opts.SetPassword("yang@5896336")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	sub(client)
	for {
		if quit {
			client.Disconnect(250)
			os.Exit(3)
		}
	}
	//publish(client)
	//Producer()
	client.Disconnect(250)
}

/*
func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message yaoyao %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}

}
*/
func sub(client mqtt.Client) {
	topic := "test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
