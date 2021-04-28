/*
 * @Author: tinson.liu
 * @Date: 2021-03-03 12:00:21
 * @LastEditors: tinson.liu
 * @LastEditTime: 2021-04-16 09:19:51
 * @Description: In User Settings Edit
 * @FilePath: /ts-go-server/dto/request/asset_host.go
 */
package request

import (
	"ts-go-server/dto/response"
	"ts-go-server/models"
)

// 获取接口列表结构体
type HostReq struct {
	HostName          string `json:"host_name" form:"host_name"`
	PublicIp          string `json:"public_ip" form:"public_ip"`
	PrivateIp         string `json:"private_ip" form:"private_ip"`
	OSVersion         string `json:"os_version" form:"os_version"`
	HostType          string `json:"host_type" form:"host_type"`
	AuthType          string `json:"auth_type" form:"auth_type"`
	Creator           string `json:"creator" form:"creator"`
	GroupID           string `json:"group_id" form:"group_id"`
	response.PageInfo        // 分页参数
}

// 创建接口结构体
type CreateHostReq struct {
	HostName          string           `json:"host_name" form:"host_name"`
	SysName           string           `json:"sys_name" form:"sys_name"`
	ResourceGroup     string           `json:"resource_group" form:"resource_group"`
	MacAddress        string           `json:"mac_address" form:"mac_address"`
	HostType          string           `json:"host_type" form:"host_type"`
	Port              string           `json:"port" form:"port"`
	AuthType          string           `json:"auth_type" form:"auth_type" validate:"required"`
	User              string           `json:"user" form:"user"`
	Password          string           `json:"password" form:"password"`
	PrivateKey        string           `json:"privatekey" form:"privatekey"`
	KeyPassphrase     string           `json:"key_passphrase"`
	Creator           string           `json:"creator" form:"creator"`
	InstanceId        string           `json:"instance_id" form:"instance_id"`
	Cpu               int              `json:"cpu" form:"cpu"`
	Memory            int              `json:"memory" form:"memory"`
	Disk              string           `json:"disk" form:"disk"`
	PrivateIp         string           `json:"private_ip" form:"private_ip"`
	PublicIp          string           `json:"public_ip" form:"public_ip"`
	BuyDate           models.LocalTime `json:"buy_date" form:"buy_date"`
	Eip               string           `json:"eip" form:"eip"`
	Owner             string           `json:"owner" form:"owner"`
	InstanceSize      string           `json:"instance_size" form:"instance_size"`
	SnNumber          string           `json:"sn_number" form:"sn_number"`
	Subnet            string           `json:"subnet" form:"subnet"`
	VirtualNetwork    string           `json:"virtual_network" form:"virtual_network"`
	OsType            string           `json:"os_type" form:"os_type"`
	Zone              string           `json:"zone" form:"zone"`
	Status            string           `json:"status" form:"status"`
	Desc              string           `json:"desc" form:"desc"`
	VpcId             string           `json:"vpc_id" form:"vpc_id"`
	ImageId           string           `json:"image_id" form:"image_id"`
	VswitchId         string           `json:"vswitch_id" form:"vswitch_id"`
	Provider          string           `json:"provider" form:"provider"`
	WarrantyDate      models.LocalTime `json:"warranty_date" form:"warranty_date"`
	OsVersion         string           `json:"os_version" form:"os_version"`
	SecurityGroupName string           `json:"security_group_name" form:"security_group_name"`
}

// SSh结构体
type SShTunnelReq struct {
	HostId uint   `json:"hostId" form:"host_id"`
	Width  int    `json:"width" form:"width"`
	Hight  int    `json:"hight" form:"hight"`
	Token  string `json:"token" form:"token"`
}

// 文件管理req
type FileReq struct {
	HostId uint   `json:"host_id" form:"host_id"` // hostId
	Path   string `json:"path" form:"path"`
	Key    string `json:"key" form:"key"`
}

// 翻译需要校验的字段名称
func (s CreateHostReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["AuthType"] = "认证类型"
	return m
}
