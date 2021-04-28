package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	getuser "ts-go-server/api/v1/system"
	"ts-go-server/dto/request"
	"ts-go-server/dto/response"
	"ts-go-server/dto/service"
	"ts-go-server/models/system"
	"ts-go-server/pkg/cloud"
	"ts-go-server/pkg/common"
	"ts-go-server/pkg/utils"

	"github.com/gin-gonic/gin"

	// "github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	// "github.com/Azure-Samples/azure-sdk-for-go-samples/internal/util"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// 获取列表
func GetHosts(c *gin.Context) {
	// 绑定参数
	var req request.HostReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	hosts, err := s.GetHosts(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.HostListResp
	utils.Struct2StructByJson(hosts, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	response.SuccessWithData(resp)
}

// 创建
func CreateHost(c *gin.Context) {
	user := getuser.GetCurrentUserFromCache(c)
	// 绑定参数
	var req request.CreateHostReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.(system.SysUser).Name
	// 创建服务
	s := service.New()
	err = s.CreateHost(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 获取当前主机信息
func GetHostInfo(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	hostId := utils.Str2Uint(c.Param("hostId"))
	if hostId == 0 {
		response.FailWithMsg("接口编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	host, err := s.GetHostById(hostId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var connStruct response.HostListResp
	utils.Struct2StructByJson(host, &connStruct)
	response.SuccessWithData(connStruct)
}

// 更新
func UpdateHostById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	hostId := utils.Str2Uint(c.Param("hostId"))
	if hostId == 0 {
		response.FailWithMsg("接口编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateHostById(hostId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批删除
func BatchDeleteHostByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteHostByIds(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

type DiskProp struct {
	DiskId interface{} `json:"DiskId"`
	Type   interface{} `json:"Type"`
	Size   interface{} `json:"Size"`
}

func ScanAzureHostNew(c *gin.Context) {
	subscription := "2f14f163-de21-4afa-b4e8-78496304745e"
	clientId := "0e9afd24-b96c-4a5e-975a-a4722ffcba8b"
	secret := "Kv_1Z.T5=BMG1utYbS/WAhyU5B9ES4Zi"
	tenant := "9d8b91e7-d3bc-4146-8c5b-311dd8a23b51"
	azureClientCredential := cloud.NewAzureClientCredentials(subscription, clientId, secret, tenant)
	azureApiClients := cloud.NewAzApiClient(subscription)
	_, err := azureApiClients.Auth(azureClientCredential)
	if err != nil {
		panic(err)
	}
	azureApiClients.ScanVm()
}

func ScanAzureHost(c *gin.Context) {
	os.Setenv("AZURE_CLIENT_ID", "0e9afd24-b96c-4a5e-975a-a4722ffcba8b")
	os.Setenv("AZURE_TENANT_ID", "9d8b91e7-d3bc-4146-8c5b-311dd8a23b51")
	os.Setenv("AZURE_CLIENT_SECRET", "Kv_1Z.T5=BMG1utYbS/WAhyU5B9ES4Zi")
	// config := NewClientCredentialsConfig(clientID, secret, tenantID)

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	nicClient := network.NewInterfacesClient("2f14f163-de21-4afa-b4e8-78496304745e")
	vmClient := compute.NewVirtualMachinesClient("2f14f163-de21-4afa-b4e8-78496304745e")
	publicIpClient := network.NewPublicIPAddressesClient("2f14f163-de21-4afa-b4e8-78496304745e")
	diskClient := compute.NewDisksClient("2f14f163-de21-4afa-b4e8-78496304745e")
	vmSizeClient := compute.NewVirtualMachineSizesClient("2f14f163-de21-4afa-b4e8-78496304745e")
	if err == nil {
		vmClient.Authorizer = authorizer
		nicClient.Authorizer = authorizer
		publicIpClient.Authorizer = authorizer
		diskClient.Authorizer = authorizer
		vmSizeClient.Authorizer = authorizer
	}

	for iter, err := vmClient.ListAllComplete(context.Background(), "false"); iter.NotDone(); err = iter.Next() {
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("found VM with name %s", *iter.Value().Name)
		fmt.Println("found VM with id %s", *iter.Value().ID)
		vmName := *iter.Value().Name
		resourceGroup := strings.Split(*iter.Value().ID, "/")[4]
		vmInstanceView, _ := vmClient.InstanceView(context.Background(), resourceGroup, vmName)
		computerName := vmInstanceView.ComputerName
		fmt.Println("ComputerName: %s", computerName)
		// fmt.Printf("instance view :%v", reflect.ValueOf(vmInstanceView))
		// fmt.Println("ComputerName: %s", *iter.Value().OsProfile.ComputerName)
		fmt.Println("Secrets: %s", *iter.Value().OsProfile.Secrets)
		fmt.Println("AdminUsername: %s", *iter.Value().OsProfile.AdminUsername)
		for _, status := range *vmInstanceView.Statuses {
			// fmt.Printf("status :%v", reflect.ValueOf(vmInstanceView))
			stateList := strings.Split(*status.Code, "/")
			if stateList[0] == "PowerState" {
				fmt.Println("PowerState is: %s", stateList[1])
			}
		}
		for _, networkInterface := range *iter.Value().NetworkProfile.NetworkInterfaces {
			fmt.Println("NetworkInterface's id is: ", *networkInterface.ID)
			nicName := strings.Split(*networkInterface.ID, "/")[8]
			nic, _ := nicClient.Get(context.Background(), resourceGroup, nicName, "")
			fmt.Println("MacAddress is: ", *nic.InterfacePropertiesFormat.MacAddress)
			if nic.InterfacePropertiesFormat.NetworkSecurityGroup != nil {
				fmt.Println("NetWorkSecurityGroup is: ", strings.Split(*nic.InterfacePropertiesFormat.NetworkSecurityGroup.ID, "/")[8])
			}
			for _, ipconfiguration := range *nic.InterfacePropertiesFormat.IPConfigurations {
				fmt.Println("PrivateIPAddress is: ", *ipconfiguration.PrivateIPAddress)
				if ipconfiguration.PublicIPAddress != nil {
					publicIpAddressName := strings.Split(*ipconfiguration.PublicIPAddress.ID, "/")[8]
					fmt.Println("network id : ", publicIpAddressName)
					publicIPAddress, _ := publicIpClient.Get(context.Background(), resourceGroup, publicIpAddressName, "")
					fmt.Println("publicIPAddress : ", *publicIPAddress.PublicIPAddressPropertiesFormat.IPAddress)
				}
			}
			fmt.Println("PrivateIPAdress is ", nic)
		}
		if iter.Value().StorageProfile.ImageReference.Offer != nil {
			fmt.Println("ImageReference Offer: %s", *iter.Value().StorageProfile.ImageReference.Offer)
		}
		if iter.Value().StorageProfile.ImageReference.Sku != nil {
			fmt.Println("ImageReference Sku: %s", *iter.Value().StorageProfile.ImageReference.Sku)
		}
		// }
		fmt.Println("VmSize : ", iter.Value().HardwareProfile.VMSize)
		vmSizeList, _ := vmSizeClient.List(context.Background(), *iter.Value().Location)
		fmt.Printf("vmsize %+v", vmSizeList.Value)
		fmt.Println("Location : ", *iter.Value().Location)

		if iter.Value().StorageProfile.DataDisks != nil {
			var diskList []DiskProp
			fmt.Println("OsDisk Name is : ", *iter.Value().StorageProfile.OsDisk.Name)
			osDisk, _ := diskClient.Get(context.Background(), resourceGroup, *iter.Value().StorageProfile.OsDisk.Name)
			fmt.Println("osdisk created at: ", osDisk.DiskProperties.TimeCreated)
			diskList = append(diskList, DiskProp{
				DiskId: *iter.Value().StorageProfile.OsDisk.Name,
				Type:   "System",
				Size:   *osDisk.DiskProperties.DiskSizeGB,
			})
			// fmt.Printf("Datadisk :%v\n", reflect.ValueOf(*iter.Value().StorageProfile.DataDisks))
			for _, disk := range *iter.Value().StorageProfile.DataDisks {
				dataDisk, _ := diskClient.Get(context.Background(), resourceGroup, *disk.Name)
				diskList = append(diskList, DiskProp{
					DiskId: *iter.Value().StorageProfile.OsDisk.Name,
					Type:   "Data",
					Size:   *dataDisk.DiskProperties.DiskSizeGB,
				})
				// fmt.Println("Datadisk Name %s", *disk.Name)
				// fmt.Println("Datadisk Size %sGB", *disk.DiskSizeGB)
				// fmt.Println("Datadisk Size %s", disk.DiskSizeGB)
			}
			// fmt.Printf("DiskList is : %+v \n", diskList)
			jsonDiskList, _ := json.Marshal(diskList)
			fmt.Println("DiskList is : ", string(jsonDiskList))
		}

	}
	// all_vm_list := &vmlist.vmlr.Value
	// for i, v := range vmlist.vmlr.Value {
	// 	fmt.Println(v.OsProfile)
	// }
}
