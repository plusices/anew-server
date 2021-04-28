/*
 * @Author: tinson.liu
 * @Date: 2021-03-03 12:00:21
 * @LastEditors: tinson.liu
 * @LastEditTime: 2021-04-14 18:09:24
 * @Description: In User Settings Edit
 * @FilePath: /ts-go-server/dto/response/asset_host.go
 */
package response

import (
	"ts-go-server/models"
)

type HostListResp struct {
	Id                uint             `json:"id"`
	HostName          string           `json:"host_name"`
	SysName           string           `json:"sys_name"`
	ResourceGroup     string           `json:"resource_group"`
	IpAddress         string           `json:"ip_address"`
	MacAddress        string           `json:"mac_address"`
	HostType          string           `json:"host_type"`
	Port              string           `json:"port"`
	AuthType          string           `json:"auth_type"`
	User              string           `json:"user"`
	Password          string           `json:"password"`
	PrivateKey        string           `json:"privatekey"`
	KeyPassphrase     string           `json:"key_passphrase"`
	Creator           string           `json:"creator"`
	InstanceId        string           `json:"instance_id"`
	Cpu               int              `json:"cpu"`
	Memory            int              `json:"memory"`
	Disk              string           `json:"disk"`
	PrivateIp         string           `json:"private_ip"`
	PublicIp          string           `json:"public_ip"`
	BuyDate           models.LocalTime `json:"buy_date"`
	Eip               string           `json:"eip"`
	Owner             string           `json:"owner"`
	InstanceSize      string           `json:"instance_size"`
	SnNumber          string           `json:"sn_number"`
	Subnet            string           `json:"subnet"`
	VirtualNetwork    string           `json:"virtual_network"`
	OsType            string           `json:"os_type"`
	Zone              string           `json:"zone"`
	Status            string           `json:"status"`
	Desc              string           `json:"desc"`
	VpcId             string           `json:"vpc_id"`
	ImageId           string           `json:"image_id"`
	VswitchId         string           `json:"vswitch_id"`
	Provider          string           `json:"provider"`
	WarrantyDate      models.LocalTime `json:"warranty_date"`
	OsVersion         string           `json:"os_version"`
	SecurityGroupName string           `json:"security_group_name"`
}

type ConnectionResp struct {
	Key         string           `json:"key"`
	UserName    string           `json:"user_name"`
	Name        string           `json:"name"`
	HostName    string           `json:"host_name"`
	IpAddress   string           `json:"ip_address"`
	Port        string           `json:"port"`
	ConnectTime models.LocalTime `json:"connect_time"`
}

type FileInfo struct {
	Name   string           `json:"name"`
	Path   string           `json:"path"`
	IsDir  bool             `json:"isDir"`
	Mode   string           `json:"mode"`
	Size   string           `json:"size"`
	Mtime  models.LocalTime `json:"mtime"` // 修改时间
	IsLink bool             `json:"isLink"`
}
