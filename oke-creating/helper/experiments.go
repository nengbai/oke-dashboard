package helper

import (
	"context"
	"fmt"
	"oke-creating/lib"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/containerengine"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

var (
	compartmentId               = "ocid1.compartment.oc1..aaaaaaaarrh2qh5ehj226yge7k75o4uk6bo6g7eud4hszl7fw4fhbtv2v77q"
	clusterName                 = "OKE-Cluster"
	vcnDisplayName              = "OKE-Cluster-VCN"
	subnetDisplayName1          = "OCI-AZ-k8sSubnet"
	subnetDisplayName2          = "OCI-AZ-nodeSubnet"
	subnetDisplayName3          = "OCI-AZ-svcSubnet"
	nodePoolName                = "pool1"
	kubeVersion                 = "v1.24.1"
	VMShape                     = "VM.Standard.E4.Flex"
	OSImage                     = "Oracle-Linux-7.9-2022.10.04-0"
	cidrBlock                   = "192.168.0.0/16"
	subnetDisplayName1cidrBlock = "192.168.0.0/28"
	subnetDisplayName2cidrBlock = "192.168.10.0/24"
	subnetDisplayName3cidrBlock = "192.168.20.0/24"
)

func CreateOKE() {
	ctx := context.Background()
	c, clerr := containerengine.NewContainerEngineClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	compute, err := core.NewComputeClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	identityClient, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)
	req := identity.ListAvailabilityDomainsRequest{}
	req.CompartmentId = common.String(compartmentId)
	ads, err := identityClient.ListAvailabilityDomains(ctx, req)
	helpers.FatalIfError(err)
	
	vcn := lib.CreateVcn(vcnDisplayName, cidrBlock, compartmentId, clusterName)

	fmt.Println("WAITING for CreateVcn")
	time.Sleep(1 * 30 * time.Second)

	internetGateway := lib.CreateInternetGateway(*vcn.Id, compartmentId)
	fmt.Println("WAITING for CreateInternetGateway")
	time.Sleep(1 * 30 * time.Second)

	natGateway := lib.CreateNatGateway(*vcn.Id, compartmentId)
	fmt.Println("WAITING for CreateNatGateway")
	time.Sleep(1 * 30 * time.Second)

	serviceGateway := lib.CreateServiceGateway(*vcn.Id, compartmentId)
	fmt.Println("WAITING for CreateServiceGateway")
	time.Sleep(1 * 30 * time.Second)

	//fmt.Println("Update for ServiceGateway")
	//err = lib.UpdateServiceGateway(*serviceGateway.Id)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	publicRouteTable := lib.CreatePublicRouteTable(*vcn.Id, compartmentId, *internetGateway.Id)
	fmt.Println("WAITING for CreatePublicRouteTable")
	time.Sleep(1 * 30 * time.Second)

	privateRouteTable := lib.CreatePrivateRouteTable(*vcn.Id, compartmentId, *natGateway.Id, *serviceGateway.Id)
	fmt.Println("WAITING for CreatePrivateRouteTable")
	time.Sleep(1 * 30 * time.Second)

	svcSubnet := lib.CreateSubnet(common.String(subnetDisplayName3), compartmentId, common.String(subnetDisplayName3cidrBlock), common.String("svcSubnetDns"), nil, vcn, publicRouteTable, nil, false) // svc subnet
	fmt.Println("WAITING for CreateSubnet")
	time.Sleep(1 * 60 * time.Second)

	k8sSecurityList := lib.CreateK8sSecurityList(*vcn.Id, compartmentId)
	k8sSubnet := lib.CreateSubnet(common.String(subnetDisplayName1), compartmentId, common.String(subnetDisplayName1cidrBlock), common.String("k8sSubnetDns"), nil, vcn, publicRouteTable, &k8sSecurityList, false)
	fmt.Println("WAITING for CreateSubnet  k8sSubnet ...")
	time.Sleep(1 * 30 * time.Second)

	nodeSecurityList := lib.CreateNodeSecurityList(*vcn.Id, compartmentId)
	nodeSubnet := lib.CreateSubnet(common.String(subnetDisplayName2), compartmentId, common.String(subnetDisplayName2cidrBlock), common.String("nodeSubnetDns"), nil, vcn, privateRouteTable, &nodeSecurityList, true)
	fmt.Println("WAITING for CreateSubnet  nodeSubnet ...")
	time.Sleep(1 * 30 * time.Second)

	//createClusterResponse := lib.CreateCluster(*vcn.Id, compartmentId, *svcSubnet.Id, *k8sSubnet.Id, kubeVersion)
	createClusterResponse := lib.CreateCluster(*vcn.Id, compartmentId, *svcSubnet.Id, *k8sSubnet.Id, clusterName, kubeVersion)

	// wait until work request complete
	workReqResp := lib.WaitUntilWorkRequestComplete(c, createClusterResponse.OpcWorkRequestId)
	fmt.Println("cluster created")
	fmt.Println("WAITING for cluster created ...")
	time.Sleep(5 * 60 * time.Second)

	clusterID := getResourceID(workReqResp.Resources, containerengine.WorkRequestResourceActionTypeCreated, "CLUSTER")
	id := k8sSubnet.Id
	fmt.Println("ID:", *clusterID, *id)
	migrateClusterResponse := lib.MigrateToVcnNativeCluster(*clusterID, *id)
	migreateReqResp := lib.WaitUntilWorkRequestComplete(c, migrateClusterResponse.OpcWorkRequestId)
	fmt.Println("cluster migrated")

	// wait until migrate complete
	getResourceID(migreateReqResp.Resources, containerengine.WorkRequestResourceActionTypeCreated, "CLUSTER")

	// // get Image Id
	image := getImageID(ctx, compute)

	fmt.Println(image)
	lib.CreateNodePool(VMShape, nodePoolName, kubeVersion, *clusterID, *image.Id, compartmentId, nodeSubnet, ads)

	// AFTER COMPLETION: create tutorial on GitHub and Notion
}

// getResourceID return a resource ID based on the filter of resource actionType and entityType
func getResourceID(resources []containerengine.WorkRequestResource, actionType containerengine.WorkRequestResourceActionTypeEnum, entityType string) *string {
	for _, resource := range resources {
		if resource.ActionType == actionType && strings.ToUpper(*resource.EntityType) == entityType {
			return resource.Identifier
		}
	}

	fmt.Println("cannot find matched resources")
	return nil
}

func getImageID(ctx context.Context, c core.ComputeClient) core.Image {
	request := core.ListImagesRequest{}
	request.CompartmentId = common.String(compartmentId)
	request.OperatingSystem = common.String("Oracle Linux")
	request.Shape = common.String(VMShape)
	request.DisplayName = common.String(OSImage)

	r, err := c.ListImages(ctx, request)
	helpers.FatalIfError(err)

	return r.Items[0]
}
