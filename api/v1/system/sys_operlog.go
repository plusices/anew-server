package system

import (
	"anew-server/dto/cacheService"
	"anew-server/dto/request"
	"anew-server/dto/response"
	"anew-server/dto/service"
	"anew-server/models/system"
	"anew-server/pkg/redis"
	"anew-server/pkg/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// 获取操作日志列表
func GetOperLogs(c *gin.Context) {
	// 绑定参数
	var req request.OperLogListReq
	reqErr := c.Bind(&req)
	if reqErr != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	var operationLogs []system.SysOperLog
	var err error
	// 创建缓存对象
	cache := cacheService.New(redis.NewStringOperation(), time.Second*15, cacheService.SERILIZER_JSON)
	key := "operationLogs"
	cache.DBGetter = func() interface{} {
		// 创建服务
		s := service.New()
		operationLogs, err = s.GetOperLogs(&req)
		return operationLogs
	}
	// 获取缓存
	cache.GetCacheForObject(key, &operationLogs)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.OperationLogListResp
	utils.Struct2StructByJson(operationLogs, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	response.SuccessWithData(resp)
}

// 批量删除操作日志
func BatchDeleteOperLogByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteOperationLogByIds(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
