/*
 * @Author: tinson.liu
 * @Date: 2021-04-02 10:02:37
 * @LastEditors: tinson.liu
 * @LastEditTime: 2021-04-12 18:18:30
 * @Description: In User Settings Edit
 * @FilePath: /ts-go-server/pkg/cloud/azure.go
 */
package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"ts-go-server/dto/request"
	"ts-go-server/dto/service"
	"ts-go-server/models"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type AzApiClients struct {
	nicClient      network.InterfacesClient
	vmClient       compute.VirtualMachinesClient
	publicIpClient network.PublicIPAddressesClient
	diskClient     compute.DisksClient
	vmSizeClient   compute.VirtualMachineSizesClient
}

type DiskProp struct {
	DiskId interface{} `json:"DiskId"`
	Type   interface{} `json:"Type"`
	Size   interface{} `json:"Size"`
}

type InstanceSize struct {
	// SizeName string `json:"SizeName"`
	Cpu    int `json:"Cpu"`
	Memory int `json:"Memory"`
}

func NewAzureClientCredentials(subscription string, clientID string, secret string, tenent string) (s auth.EnvironmentSettings) {
	s = auth.EnvironmentSettings{
		Values: map[string]string{},
	}
	s.Values[auth.SubscriptionID] = subscription
	s.Values[auth.ClientID] = clientID
	s.Values[auth.ClientSecret] = secret
	s.Values[auth.TenantID] = tenent
	if v := s.Values[auth.EnvironmentName]; v == "" {
		s.Environment = azure.PublicCloud
	} else {
		s.Environment, _ = azure.EnvironmentFromName(v)
	}
	if s.Values[auth.Resource] == "" {
		s.Values[auth.Resource] = s.Environment.ResourceManagerEndpoint
	}
	return
}

func NewAzApiClient(subscription string) AzApiClients {
	return AzApiClients{
		// subscription : subscription,
		vmClient:       compute.NewVirtualMachinesClient(subscription),
		diskClient:     compute.NewDisksClient(subscription),
		vmSizeClient:   compute.NewVirtualMachineSizesClient(subscription),
		nicClient:      network.NewInterfacesClient(subscription),
		publicIpClient: network.NewPublicIPAddressesClient(subscription),
	}
}

func (clients *AzApiClients) Auth(s auth.EnvironmentSettings) (autorest.Authorizer, error) {
	authorizer, err := s.GetAuthorizer()
	if err != nil {
		return nil, err
	}
	clients.vmClient.Authorizer = authorizer
	clients.nicClient.Authorizer = authorizer
	clients.publicIpClient.Authorizer = authorizer
	clients.diskClient.Authorizer = authorizer
	clients.vmSizeClient.Authorizer = authorizer
	return authorizer, err
}

func (clients *AzApiClients) getDiskInfo(req *request.CreateHostReq, iter *compute.VirtualMachineListResultIterator) {
	if iter.Value().StorageProfile.DataDisks != nil {
		var diskList []DiskProp
		osDisk, _ := clients.diskClient.Get(context.Background(), req.ResourceGroup, *iter.Value().StorageProfile.OsDisk.Name)
		req.OsType = string(osDisk.OsType)
		req.BuyDate = models.LocalTime{Time: osDisk.DiskProperties.TimeCreated.Time}
		diskList = append(diskList, DiskProp{
			DiskId: *iter.Value().StorageProfile.OsDisk.Name,
			Type:   "System",
			Size:   *osDisk.DiskProperties.DiskSizeGB,
		})
		for _, disk := range *iter.Value().StorageProfile.DataDisks {
			dataDisk, _ := clients.diskClient.Get(context.Background(), req.ResourceGroup, *disk.Name)
			diskList = append(diskList, DiskProp{
				DiskId: *iter.Value().StorageProfile.OsDisk.Name,
				Type:   "Data",
				Size:   *dataDisk.DiskProperties.DiskSizeGB,
			})
		}
		jsonDiskList, _ := json.Marshal(diskList)
		req.Disk = string(jsonDiskList)
	}
}

func (clients *AzApiClients) getNetworkInfo(req *request.CreateHostReq, iter *compute.VirtualMachineListResultIterator) {
	for _, networkInterface := range *iter.Value().NetworkProfile.NetworkInterfaces {
		nicName := strings.Split(*networkInterface.ID, "/")[8]
		nic, _ := clients.nicClient.Get(context.Background(), req.ResourceGroup, nicName, "")
		req.MacAddress = *nic.InterfacePropertiesFormat.MacAddress
		if nic.InterfacePropertiesFormat.NetworkSecurityGroup != nil {
			req.SecurityGroupName = strings.Split(*nic.InterfacePropertiesFormat.NetworkSecurityGroup.ID, "/")[8]
		}
		for _, ipconfiguration := range *nic.InterfacePropertiesFormat.IPConfigurations {
			fmt.Println("PrivateIPAddress is: ", *ipconfiguration.PrivateIPAddress)
			req.VirtualNetwork = strings.Split(*ipconfiguration.Subnet.ID, "/")[8]
			req.Subnet = strings.Split(*ipconfiguration.Subnet.ID, "/")[10]
			req.PrivateIp = *ipconfiguration.PrivateIPAddress
			if ipconfiguration.PublicIPAddress != nil {
				publicIpAddressName := strings.Split(*ipconfiguration.PublicIPAddress.ID, "/")[8]
				publicIPAddress, _ := clients.publicIpClient.Get(context.Background(), req.ResourceGroup, publicIpAddressName, "")
				req.PublicIp = *publicIPAddress.PublicIPAddressPropertiesFormat.IPAddress
			}
		}
	}
}

func (clients *AzApiClients) listVmSizes(req *request.CreateHostReq, vmSizeMap map[string]map[string]InstanceSize) {
	if _, ok := vmSizeMap[req.Zone]; !ok {
		vmSizeList, _ := clients.vmSizeClient.List(context.Background(), req.Zone)
		subMap := make(map[string]InstanceSize)
		for _, vmSize := range *vmSizeList.Value {
			subMap[string(*vmSize.Name)] = InstanceSize{Cpu: int(*vmSize.NumberOfCores), Memory: int(*vmSize.MemoryInMB)}
		}
		vmSizeMap[req.Zone] = subMap
	}
}

func (clients *AzApiClients) ScanVm() {
	var req request.CreateHostReq
	vmSizeMap := map[string]map[string]InstanceSize{}
	s := service.New()
	for iter, err := clients.vmClient.ListAllComplete(context.Background(), "false"); iter.NotDone(); err = iter.Next() {
		if err != nil {
			fmt.Println(err)
		}
		req.InstanceId = *iter.Value().VMID
		vmName := *iter.Value().Name
		req.HostName = *iter.Value().Name
		req.ResourceGroup = strings.Split(*iter.Value().ID, "/")[4]
		vmInstanceView, _ := clients.vmClient.InstanceView(context.Background(), req.ResourceGroup, vmName)
		computerName := func() string {
			if vmInstanceView.ComputerName == nil {
				return *iter.Value().Name
			}
			return *vmInstanceView.ComputerName
		}()
		req.SysName = computerName
		req.User = *iter.Value().OsProfile.AdminUsername
		for _, status := range *vmInstanceView.Statuses {
			stateList := strings.Split(*status.Code, "/")
			if stateList[0] == "PowerState" {
				req.Status = stateList[1]
			}
		}
		clients.getNetworkInfo(&req, &iter)
		if iter.Value().StorageProfile.ImageReference.Offer != nil {
			req.OsVersion = *iter.Value().StorageProfile.ImageReference.Offer
		}
		if iter.Value().StorageProfile.ImageReference.Sku != nil {
			req.OsVersion = req.OsVersion + " " + *iter.Value().StorageProfile.ImageReference.Sku
		}
		req.Zone = *iter.Value().Location
		req.InstanceSize = string(iter.Value().HardwareProfile.VMSize)
		clients.listVmSizes(&req, vmSizeMap)
		req.Cpu = vmSizeMap[req.Zone][string(req.InstanceSize)].Cpu
		req.Memory = int(vmSizeMap[req.Zone][string(req.InstanceSize)].Memory / 1024)
		clients.getDiskInfo(&req, &iter)
		clients.getNetworkInfo(&req, &iter)
		_ = s.UpdateHostByField("InstanceId", &req)
	}
}
