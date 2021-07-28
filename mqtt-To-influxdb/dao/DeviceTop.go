package dao

import (
	"example.com/m/model"
	"example.com/m/tool"
)

type DeviceIoT struct {
	*tool.Orm
}

//实例化Dao对象
func NewDeviceIoTDao() *DeviceIoT {
	return &DeviceIoT{tool.DbEngine}
}

//从数据库中查询所有topic
func (DevIoT *DeviceIoT) QuerydeviceIots() ([]model.IOT_DeviceTopic, error) {
	var deviceIots []model.IOT_DeviceTopic
	if err := DevIoT.Engine.Find(&deviceIots); err != nil {
		return nil, err
	}
	return deviceIots, nil
}

//按类型从数据库中查询所有topic
func (DevIoT *DeviceIoT) QuerydeviceIotsByType() ([]model.IOT_DeviceTopic, error) {
	var deviceIots []model.IOT_DeviceTopic
	if err := DevIoT.Engine.Where("TopicType = 1 and FunctionType = 7 and OptAuthType =0").Find(&deviceIots); err != nil {
		return nil, err
	}
	return deviceIots, nil
}
