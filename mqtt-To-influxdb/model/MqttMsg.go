package model

type MqttMsg struct {
	RH  float64 `json:"RH"`
	TMP float64 `json:"TMP"`
}

type Iot_devicetopic1 struct {
	Id             int32  `xorm:"pk autoincr" json:"ID32"`
	Devicename     string `xorm:"varchar(32)" json:"DeviceName"`
	OptAuthtype    string `xorm:"char(1)" json:"OptAuthType"`
	Topicname      string `xorm:" varchar(256)" json:"TopicName"`
	Customtopic    string `xorm:" varchar(256)" json:"CustomTopic"`
	Descr          string `xorm:" varchar(1000)" json:"Descr"`
	Topictype      int    `xorm:" int(1)" json:"TopicType"`
	Functiontype   int    `xorm:" int(1)" json:"FunctionType"`
	Producttopicid int    `xorm:" int(11)" json:"ProductTopicId"`
	Datastatesmark string `xorm:" varchar(4)" json:"DataStatesMark"`
	Syncflag       int    `xorm:" int(10)" json:"SyncFlag"`
	Iphyindex      int    `xorm:" int(10)" json:"iPhyIndex"`
}
type IOT_DeviceTopic1 struct {
	ID32           int32  `xorm:"pk autoincr" json:"ID32"`
	DeviceName     string `xorm:"varchar(32)" json:"DeviceName"`
	OptAuthType    string `xorm:"char(1)" json:"OptAuthType"`
	TopicName      string `xorm:" varchar(256)" json:"TopicName"`
	CustomTopic    string `xorm:" varchar(256)" json:"CustomTopic"`
	Descr          string `xorm:" varchar(1000)" json:"Descr"`
	TopicType      int    `xorm:" int(1)" json:"TopicType"`
	FunctionType   int    `xorm:" int(1)" json:"FunctionType"`
	ProductTopicId int    `xorm:" int(11)" json:"ProductTopicId"`
	DataStatesMark string `xorm:" varchar(4)" json:"DataStatesMark"`
	SyncFlag       int    `xorm:" int(10)" json:"SyncFlag"`
	IPhyIndex      int    `xorm:" int(10)" json:"iPhyIndex"`
}
type IOT_DeviceTopic struct {
	ID32           int32  `xorm:"pk autoincr" json:"ID32"`
	DeviceName     string `xorm:"varchar(32)" json:"DeviceName"`
	OptAuthType    string `xorm:"char(1)" json:"OptAuthType"`
	TopicName      string `xorm:" varchar(256)" json:"TopicName"`
	CustomTopic    string `xorm:" varchar(256)" json:"CustomTopic"`
	MessageCount   string `xorm:" varchar(256)" json:"MessageCount"`
	Descr          string `xorm:" varchar(1000)" json:"Descr"`
	TopicType      int    `xorm:" int(1)" json:"TopicType"`
	FunctionType   int    `xorm:" int(1)" json:"FunctionType"`
	ProductTopicId int    `xorm:" int(11)" json:"ProductTopicId"`
	DataStatesMark string `xorm:" varchar(4)" json:"DataStatesMark"`
	SyncFlag       int    `xorm:" int(10)" json:"SyncFlag"`
	IPhyIndex      int    `xorm:" int(10)" json:"iPhyIndex"`
}
