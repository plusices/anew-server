package asset

import (
	getuser "anew-server/api/v1/system"
	"anew-server/dto/request"
	"anew-server/dto/response"
	"anew-server/dto/service"
	"anew-server/models/system"
	"anew-server/pkg/common"
	"anew-server/pkg/utils"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	// "github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	// "github.com/Azure-Samples/azure-sdk-for-go-samples/internal/util"
	"github.com/Azure/go-autorest/autorest/azure/auth"

	// "github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
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

func ScanAzureHost(c *gin.Context) {
	os.Setenv("AZURE_CLIENT_ID", "0e9afd24-b96c-4a5e-975a-a4722ffcba8b")
	os.Setenv("AZURE_TENANT_ID", "9d8b91e7-d3bc-4146-8c5b-311dd8a23b51")
	os.Setenv("AZURE_CLIENT_SECRET", "Kv_1Z.T5=BMG1utYbS/WAhyU5B9ES4Zi")
	authorizer, err := auth.NewAuthorizerFromEnvironment()

	vmClient := compute.NewVirtualMachinesClient("2f14f163-de21-4afa-b4e8-78496304745e")
	if err == nil {
		vmClient.Authorizer = authorizer
	}
	vmClient.Authorizer = authorizer

	for iter, err := vmClient.ListAllComplete(context.Background(), "false"); iter.NotDone(); err = iter.Next() {
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("found VM with name %s", *iter.Value().Name)
		fmt.Println("ComputerName: %s", *iter.Value().OsProfile.ComputerName)
		fmt.Println("Secrets: %s", *iter.Value().OsProfile.Secrets)
		fmt.Println("AdminUsername: %s", *iter.Value().OsProfile.AdminUsername)
		// fmt.Printf("%+v\n", *iter.Value().StorageProfile.ImageReference)
		// fmt.Println(reflect.TypeOf(*iter))
		// if *iter.Value().StorageProfile != nil {
		if iter.Value().StorageProfile.ImageReference.Offer != nil {
			fmt.Println("ImageReference Offer: %s", *iter.Value().StorageProfile.ImageReference.Offer)
		}
		if iter.Value().StorageProfile.ImageReference.Sku != nil {
			fmt.Println("ImageReference Sku: %s", *iter.Value().StorageProfile.ImageReference.Sku)
		}
		// }
		// fmt.Println("found VM with name %s", *iter.Value().HardwareProfile.VMSize)
		fmt.Println("Location: %s", *iter.Value().Location)
		// if *iter.Value().StorageProfile.DataDisks != nil {
		// 	for disk := range *iter.Value().StorageProfile.DataDisks {
		// 		fmt.Println("found VM with name %s", disk.Name)
		// 		fmt.Println("found VM with name %s", disk.DiskSizeGB)
		// 	}
		// }

	}
	// all_vm_list := &vmlist.vmlr.Value
	// for i, v := range vmlist.vmlr.Value {
	// 	fmt.Println(v.OsProfile)
	// }
}
