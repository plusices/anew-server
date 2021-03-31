/*
 * @Author: tinson.liu
 * @Date: 2021-03-03 12:00:21
 * @LastEditors: tinson.liu
 * @LastEditTime: 2021-03-29 15:12:22
 * @Description: In User Settings Edit
 * @FilePath: /anew-server/models/asset/asset_host.go
 */
package asset

import (
	"anew-server/models"
	"time"
)

// 主机表
type AssetHost struct {
	models.Model
	HostName          string       `gorm:"comment:'主机名';size:128" json:"host_name"`
	InstanceId        string       `gorm:"comment:'实例id';size:128" json:"instance_id"`
	Cpu               int          `gorm:"comment:'cpu';" json:"cpu"`
	Memory            int          `gorm:"comment:'内存(单位M)'" json:"memory"`
	Disk              int          `gorm:"comment:'磁盘容量'" json:"disk"`
	IpAddress         string       `gorm:"comment:'主机地址';size:128" json:"ip_address"`
	PrivateIp         string       `gorm:"comment:'私有ip';size:64" json:"private_ip"`
	PublicIp          string       `gorm:"comment:'公有ip';size:64" json:""public_ip`
	BuyDate           time.Time    `gorm:"comment:购买日期;" json:"buy_date"`
	Eip               string       `gorm:"comment:'弹性ip';size:64" json:"eip"`
	Owner             string       `gorm:"comment:'负责人';size:64" json:"owner"`
	InstanceSize      string       `gorm:"comment:'实例类型;size:128'" json:"instance_size"`
	SnNumber          string       `gorm:"comment:'SN序列号';size:128" json:"sn_number"`
	Subnet            string       `gorm:"comment:'Azure subnet';size:128" json:"subnet"`
	VirtualNetwork    string       `gorm:"comment:'Azure virtual_network';size:128" json:"virtual_network"`
	Port              string       `gorm:"comment:'SSh端口';size:64" json:"port"`
	OsType            string       `gorm:"comment:'系统类型';size:64" json:"os_type"`
	Zone              string       `gorm:"comment:'区域';size:64" json:"zone"`
	Status            string       `gorm:"comment:'状态';size:64" json:"status"`
	Desc              string       `gorm:"comment:'描述信息';size:256" json:"desc"`
	VpcId             string       `gorm:"comment:'VPC网络id';size:64" json:"vpc_id"`
	ImageId           string       `gorm:"comment:'镜像id';size:128" json:"image_id"`
	VswitchId         string       `gorm:"comment:'虚拟交换机id';size:64" json:"vswitch_id"`
	Provider          string       `gorm:"comment:'服务商';size:64" json:"provider"`
	WarrantyDate      time.Time    `gorm:"comment:'到保日期'" json:"warranty_date"`
	OsVersion         string       `gorm:"comment:'系统版本';size:128" json:"os_version"`
	HostType          string       `gorm:"comment:'主机类型';size:64" json:"host_type"`
	AuthType          string       `gorm:"comment:'认证类型'" json:"auth_type"`
	User              string       `gorm:"comment:'认证用户';size:64" json:"user"`
	Password          string       `gorm:"comment:'认证密码';size:64" json:"password"`
	PrivateKey        string       `gorm:"comment:'秘钥';size:128" json:"privatekey"`
	KeyPassphrase     string       `gorm:"comment:'秘钥';size:64" json:"key_passphrase"`
	Creator           string       `gorm:"comment:'创建人';size:64" json:"creator"`
	SecurityGroupName string       `gorm:"comment:'安全组';size:128" json:"security_group_name"`
	Groups            []AssetGroup `gorm:"many2many:relation_group_host;" json:"groups"`
}

func (m AssetHost) TableName() string {
	return m.Model.TableName("asset_host")
}
