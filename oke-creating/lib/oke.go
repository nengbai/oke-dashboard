package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/containerengine"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"
	"github.com/oracle/oci-go-sdk/v65/identity"
)
func GetVcn(vcnDisplayName, cidrBlock, compartmentId, clusterName string) (core.Vcn,error){
	log.Println("CREATE VCN")
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	fmt.Println(clerr)

	ctx := context.Background()
	request := core.CreateVcnRequest{}
	request.CidrBlock = common.String(cidrBlock)
	request.CompartmentId = common.String(compartmentId)
	request.DisplayName = common.String(vcnDisplayName)

	r, err := c.CreateVcn(ctx, request)
	helpers.FatalIfError(err)

	return r.Vcn,nil
}

func CreateVcn(vcnDisplayName, cidrBlock, compartmentId, clusterName string) core.Vcn {
	log.Println("CREATE VCN")
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	fmt.Println(clerr)
	ctx := context.Background()

	request := core.CreateVcnRequest{}
	request.CidrBlock = common.String(cidrBlock)
	request.CompartmentId = common.String(compartmentId)
	request.DisplayName = common.String(vcnDisplayName)
	request.DnsLabel = common.String("vcndns")

	r, err := c.CreateVcn(ctx, request)
	helpers.FatalIfError(err)
	return r.Vcn
}

func CreateInternetGateway(vcnId, compartmentId string) core.InternetGateway {
	log.Println("CREATE INTERNET GATEWAY ", vcnId)
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	createInternetGatewayRequest := core.CreateInternetGatewayRequest{
		CreateInternetGatewayDetails: core.CreateInternetGatewayDetails{
			CompartmentId: common.String(compartmentId),
			VcnId:         common.String(vcnId),
			IsEnabled:     common.Bool(true),
			DisplayName:   common.String("OKE-Internet-Gateway"),
		},
	}

	igw, err := c.CreateInternetGateway(ctx, createInternetGatewayRequest)
	helpers.FatalIfError(err)

	return igw.InternetGateway
}

func CreateNatGateway(vcnId, compartmentId string) core.NatGateway {
	log.Println("CREATE NAT GATEWAY ", vcnId)
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	createNatGatewayRequest := core.CreateNatGatewayRequest{
		CreateNatGatewayDetails: core.CreateNatGatewayDetails{
			CompartmentId: common.String(compartmentId),
			VcnId:         common.String(vcnId),
			DisplayName:   common.String("OKE-Nat-Gateway"),
		},
	}

	nat, err := c.CreateNatGateway(ctx, createNatGatewayRequest)
	helpers.FatalIfError(err)

	return nat.NatGateway
}

func CreateServiceGateway(vcnId, compartmentId string) core.ServiceGateway {
	log.Println("CREATE SERVICE GATEWAY ", vcnId)
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	createServiceGatewayRequest := core.CreateServiceGatewayRequest{
		CreateServiceGatewayDetails: core.CreateServiceGatewayDetails{
			CompartmentId: common.String(compartmentId),
			VcnId:         common.String(vcnId),
			Services:      []core.ServiceIdRequestDetails{}, DisplayName: common.String("OKE-Svc-Gateway"),
			//RouteTableId: common.String("all-nrt-services-in-oracle-services-network"),
		},
		OpcRetryToken:   new(string),
		OpcRequestId:    new(string),
		RequestMetadata: common.RequestMetadata{},
	}

	sgw, err := c.CreateServiceGateway(ctx, createServiceGatewayRequest)
	helpers.FatalIfError(err)
	//sgw.ServiceGateway.Services = append(sgw.ServiceGateway.Services, common.String("all-nrt-services-in-oracle-services-network"))
	//fmt.Println(sgw.ServiceGateway.Services)

	return sgw.ServiceGateway
}

func UpdateServiceGateway(serviceId string) error {
	// Create a default authentication provider that uses the DEFAULT
	// profile in the configuration file.
	// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	// Create a request and dependent object(s).
	req := core.UpdateServiceGatewayRequest{IfMatch: common.String("OKE-Svc-Gateway"),
		ServiceGatewayId: common.String(serviceId),
		UpdateServiceGatewayDetails: core.UpdateServiceGatewayDetails{BlockTraffic: common.Bool(false),
			Services: []core.ServiceIdRequestDetails{core.ServiceIdRequestDetails{ServiceId: common.String("All NRT Services In Oracle Services Network")}}}}

	// Send the request using the service client
	resp, err := client.UpdateServiceGateway(context.Background(), req)
	if err != nil {
		return err
	}

	// Retrieve value from the response.
	fmt.Println(resp)
	return nil
}

func CreatePrivateRouteTable(vcnId, compartmentId, natId, serviceId string) core.RouteTable {
	log.Println("CREATE PRIVATE ROUTE TABLE ", vcnId)
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	createRouteTableRequest := core.CreateRouteTableRequest{
		CreateRouteTableDetails: core.CreateRouteTableDetails{
			CompartmentId: common.String(compartmentId),
			VcnId:         common.String(vcnId),
			DisplayName:   common.String("private-route-table"),
			RouteRules: []core.RouteRule{
				{
					NetworkEntityId: common.String(natId), // NAT gateway id
					Description:     common.String("traffic to the internet"),
					Destination:     common.String("0.0.0.0/0"),
					DestinationType: core.RouteRuleDestinationTypeCidrBlock,
				},
				{
					NetworkEntityId: common.String(serviceId), // service gateway id
					Description:     common.String("traffic to OCI services"),
					Destination:     common.String("all-nrt-services-in-oracle-services-network"),
					DestinationType: core.RouteRuleDestinationTypeServiceCidrBlock,
				},
			},
		},
	}

	rt, err := c.CreateRouteTable(ctx, createRouteTableRequest)
	helpers.FatalIfError(err)

	return rt.RouteTable
}

func CreatePublicRouteTable(vcnId, compartmentId, igwId string) core.RouteTable {
	log.Println("CREATE PUBLIC ROUTE TABLE ", vcnId)
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	createRouteTableRequest := core.CreateRouteTableRequest{
		CreateRouteTableDetails: core.CreateRouteTableDetails{
			CompartmentId: common.String(compartmentId),
			VcnId:         common.String(vcnId),
			DisplayName:   common.String("public-route-table"),
			RouteRules: []core.RouteRule{
				{
					NetworkEntityId: common.String(igwId), // internet gateway id
					Description:     common.String("traffic to/from internet"),
					Destination:     common.String("0.0.0.0/0"),
					DestinationType: core.RouteRuleDestinationTypeCidrBlock,
				},
			},
		},
	}

	rt, err := c.CreateRouteTable(ctx, createRouteTableRequest)
	helpers.FatalIfError(err)

	return rt.RouteTable
}

func CreateK8sSecurityList(vcnId, compartmentId string) core.SecurityList {
	log.Println("CREATE K8S SECURITY LIST ", vcnId)
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	createSecurityListRequest := core.CreateSecurityListRequest{
		CreateSecurityListDetails: core.CreateSecurityListDetails{
			CompartmentId: common.String(compartmentId),
			VcnId:         common.String(vcnId),
			DisplayName:   common.String("k8s-security-list"),
			IngressSecurityRules: []core.IngressSecurityRule{
				{
					Protocol:    common.String("1"), // ICMP
					Source:      common.String("192.168.10.0/24"),
					Description: common.String("Path discovery"),
					IcmpOptions: &core.IcmpOptions{
						Type: common.Int(3),
						Code: common.Int(4),
					},
				},
				{
					Protocol:    common.String("6"), // TCP
					Source:      common.String("0.0.0.0/0"),
					Description: common.String("External access to Kubernetes API endpoint"),
					TcpOptions: &core.TcpOptions{
						DestinationPortRange: &core.PortRange{
							Max: common.Int(6443),
							Min: common.Int(6443),
						},
					},
				},
				{
					Protocol:    common.String("6"), // TCP
					Source:      common.String("192.168.10.0/24"),
					Description: common.String("Kubernetes worker to Kubernetes API endpoint communication"),
					TcpOptions: &core.TcpOptions{
						DestinationPortRange: &core.PortRange{
							Max: common.Int(6443),
							Min: common.Int(6443),
						},
					},
				},
				{
					Protocol:    common.String("6"), // TCP
					Source:      common.String("192.168.10.0/24"),
					Description: common.String("Kubernetes worker to control plane communication"),
					TcpOptions: &core.TcpOptions{
						DestinationPortRange: &core.PortRange{
							Max: common.Int(12250),
							Min: common.Int(12250),
						},
					},
				},
			},
			EgressSecurityRules: []core.EgressSecurityRule{
				{
					Protocol:        common.String("6"), // TCP
					Description:     common.String("Allow Kubernetes Control Plane to communicate with OKE"),
					Destination:     common.String("all-nrt-services-in-oracle-services-network"),
					DestinationType: core.EgressSecurityRuleDestinationTypeServiceCidrBlock,
					TcpOptions: &core.TcpOptions{
						DestinationPortRange: &core.PortRange{
							Max: common.Int(443),
							Min: common.Int(443),
						},
					},
				},
				{
					Protocol:        common.String("6"), // TCP
					Description:     common.String("All traffic to worker nodes"),
					Destination:     common.String("192.168.10.0/24"),
					DestinationType: core.EgressSecurityRuleDestinationTypeCidrBlock,
				},
				{
					Protocol:        common.String("1"), // ICMP
					Description:     common.String("Path discovery"),
					Destination:     common.String("192.168.10.0/24"),
					DestinationType: core.EgressSecurityRuleDestinationTypeCidrBlock,
					IcmpOptions: &core.IcmpOptions{
						Type: common.Int(3),
						Code: common.Int(4),
					},
				},
			},
		},
	}

	s, err := c.CreateSecurityList(ctx, createSecurityListRequest)
	helpers.FatalIfError(err)

	return s.SecurityList
}

func CreateNodeSecurityList(vcnId, compartmentId string) core.SecurityList {
	log.Println("CREATE NODE SECURITY LIST ", vcnId)
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	createSecurityListRequest := core.CreateSecurityListRequest{
		CreateSecurityListDetails: core.CreateSecurityListDetails{
			CompartmentId: common.String(compartmentId),
			VcnId:         common.String(vcnId),
			DisplayName:   common.String("node-security-list"),
			IngressSecurityRules: []core.IngressSecurityRule{
				{
					Protocol:    common.String("all"), // ICMP
					Source:      common.String("192.168.10.0/24"),
					Description: common.String("Allow pods on one worker node to communicate with pods on other worker nodes"),
				},
				{
					Protocol:    common.String("1"), // ICMP
					Source:      common.String("192.168.0.0/28"),
					Description: common.String("Path discovery"),
					IcmpOptions: &core.IcmpOptions{
						Type: common.Int(3),
						Code: common.Int(4),
					},
				},
				{
					Protocol:    common.String("6"), // TCP
					Source:      common.String("192.168.0.0/28"),
					Description: common.String("TCP access from Kubernetes Control Plane"),
				},
				{
					Protocol:    common.String("6"), // TCP
					Source:      common.String("0.0.0.0/0"),
					Description: common.String("Inbound SSH traffic to worker nodes"),
					TcpOptions: &core.TcpOptions{
						DestinationPortRange: &core.PortRange{
							Max: common.Int(22),
							Min: common.Int(22),
						},
					},
				},
			},
			EgressSecurityRules: []core.EgressSecurityRule{
				{
					Protocol:    common.String("all"),
					Description: common.String("Allow pods on one worker node to communicate with pods on other worker nodes"),
					Destination: common.String("192.168.10.0/24"),
				},
				{
					Protocol:    common.String("6"), // TCP
					Description: common.String("Access to Kubernetes API Endpoint"),
					Destination: common.String("192.168.0.0/28"),
					TcpOptions: &core.TcpOptions{
						DestinationPortRange: &core.PortRange{
							Max: common.Int(6443),
							Min: common.Int(6443),
						},
					},
				},
				{
					Protocol:    common.String("6"), // TCP
					Description: common.String("Kubernetes worker to control plane communication"),
					Destination: common.String("192.168.0.0/28"),
					TcpOptions: &core.TcpOptions{
						DestinationPortRange: &core.PortRange{
							Max: common.Int(12250),
							Min: common.Int(12250),
						},
					},
				},
				{
					Protocol:    common.String("1"), // ICMP
					Description: common.String("Path discovery"),
					Destination: common.String("192.168.0.0/28"),
					IcmpOptions: &core.IcmpOptions{
						Type: common.Int(3),
						Code: common.Int(4),
					},
				},
				{
					Protocol:        common.String("6"), // TCP
					Description:     common.String("Allow nodes to communicate with OKE to ensure correct start-up and continued functioning"),
					Destination:     common.String("all-nrt-services-in-oracle-services-network"),
					DestinationType: core.EgressSecurityRuleDestinationTypeServiceCidrBlock,
					TcpOptions: &core.TcpOptions{
						DestinationPortRange: &core.PortRange{
							Max: common.Int(443),
							Min: common.Int(443),
						},
					},
				},
				{
					Protocol:    common.String("1"), // ICMP
					Description: common.String("ICMP Access from Kubernetes Control Plane"),
					Destination: common.String("0.0.0.0/0"),
					IcmpOptions: &core.IcmpOptions{
						Type: common.Int(3),
						Code: common.Int(4),
					},
				},
				{
					Protocol:    common.String("all"),
					Description: common.String("Worker Nodes access to Internet"),
					Destination: common.String("0.0.0.0/0"),
				},
			},
		},
	}

	s, err := c.CreateSecurityList(ctx, createSecurityListRequest)
	helpers.FatalIfError(err)

	return s.SecurityList
}

func CreateSubnet(displayName *string, compartmentId string, cidrBlock *string, dnsLabel *string, availableDomain *string, vcn core.Vcn, routeTable core.RouteTable, securityList *core.SecurityList, envVarKey bool) core.Subnet {
	log.Println("CREATE SUBNET ", *displayName)
	c, clerr := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	ctx := context.Background()

	request := core.CreateSubnetRequest{}
	request.AvailabilityDomain = availableDomain
	request.CompartmentId = common.String(compartmentId)
	request.CidrBlock = cidrBlock
	request.DisplayName = displayName
	request.DnsLabel = dnsLabel
	request.RequestMetadata = helpers.GetRequestMetadataWithDefaultRetryPolicy()
	request.VcnId = vcn.Id
	request.ProhibitInternetIngress = common.Bool(envVarKey)
	request.RouteTableId = routeTable.Id
	if securityList != nil {
		log.Println("SUBNET SECURITY LIST ID:", *securityList.Id)
		request.SecurityListIds = []string{*securityList.Id}
	}

	r, err := c.CreateSubnet(ctx, request)
	helpers.FatalIfError(err)
	log.Println("Subnet created")

	// retry condition check, stop until return true
	pollUntilAvailable := func(r common.OCIOperationResponse) bool {
		if converted, ok := r.Response.(core.GetSubnetResponse); ok {
			log.Println(converted.LifecycleState)
			return converted.LifecycleState != core.SubnetLifecycleStateAvailable
		}
		return true
	}

	pollGetRequest := core.GetSubnetRequest{
		SubnetId:        r.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(pollUntilAvailable),
	}

	// wait for lifecyle become running
	_, pollErr := c.GetSubnet(ctx, pollGetRequest)
	helpers.FatalIfError(pollErr)

	return r.Subnet
}

/*
	func CreateCluster(vcnId, compartmentId, svcSubnetId, displayName, kubernetesVersion string) containerengine.CreateClusterResponse {
		log.Println("CREATE CLUSTER")
		ctx := context.Background()
		c, clerr := containerengine.NewContainerEngineClientWithConfigurationProvider(common.DefaultConfigProvider())
		helpers.FatalIfError(clerr)
		createClusterRequest := containerengine.CreateClusterRequest{
			CreateClusterDetails: containerengine.CreateClusterDetails{
				Name:              common.String(displayName),
				CompartmentId:     common.String(compartmentId),
				VcnId:             common.String(vcnId),
				KubernetesVersion: common.String(kubernetesVersion),
				Options: &containerengine.ClusterCreateOptions{
					ServiceLbSubnetIds: []string{svcSubnetId},
				},
			},
		}
		resp, err := c.CreateCluster(ctx, createClusterRequest)
		helpers.FatalIfError(err)

		return resp
	}
*/

func GetCluster() error{
	return nil
}

func CreateCluster(vcnId, compartmentId, svcSubnetId, k8sSubnetId, displayName, kubernetesVersion string) containerengine.CreateClusterResponse {
	log.Println("CREATE CLUSTER")

	ctx := context.Background()
	c, clerr := containerengine.NewContainerEngineClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	createClusterRequest := containerengine.CreateClusterRequest{
		CreateClusterDetails: containerengine.CreateClusterDetails{
			Name:              common.String(displayName),
			CompartmentId:     common.String(compartmentId),
			VcnId:             common.String(vcnId),
			KubernetesVersion: common.String(kubernetesVersion),
			Options: &containerengine.ClusterCreateOptions{
				ServiceLbSubnetIds: []string{svcSubnetId},
			},
		},
	}
	resp, err := c.CreateCluster(ctx, createClusterRequest)
	helpers.FatalIfError(err)

	return resp
}

func MigrateToVcnNativeCluster(clusterId, k8sSubnetId string) containerengine.ClusterMigrateToNativeVcnResponse {
	log.Println("MIGRATE CLUSTER TO NATIVE-VCN")
	ctx := context.Background()
	c, clerr := containerengine.NewContainerEngineClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)
	migrateRequest := containerengine.ClusterMigrateToNativeVcnRequest{
		ClusterId: &clusterId,
		ClusterMigrateToNativeVcnDetails: containerengine.ClusterMigrateToNativeVcnDetails{
			EndpointConfig: &containerengine.ClusterEndpointConfig{
				SubnetId:          &k8sSubnetId,
				IsPublicIpEnabled: common.Bool(true),
			},
		},
	}
	resp, err := c.ClusterMigrateToNativeVcn(ctx, migrateRequest)
	helpers.FatalIfError(err)

	return resp
}

func CreateNodePool(VMShape, nodePoolName, kubernetesVersion, clusterId, imageId, compartmentId string, subnet core.Subnet, ads identity.ListAvailabilityDomainsResponse) {
	log.Println("CREATE NODE POOL ", nodePoolName)
	size := len(ads.Items)
	createNodePoolRequest := containerengine.CreateNodePoolRequest{
		CreateNodePoolDetails: containerengine.CreateNodePoolDetails{
			CompartmentId:     &compartmentId,
			ClusterId:         &clusterId,
			Name:              &nodePoolName,
			KubernetesVersion: &kubernetesVersion,
			NodeShape:         common.String(VMShape),
			NodeShapeConfig: &containerengine.CreateNodeShapeConfigDetails{
				Ocpus: common.Float32(4),
			},
			NodeConfigDetails: &containerengine.CreateNodePoolNodeConfigDetails{
				Size:             &size,
				PlacementConfigs: make([]containerengine.NodePoolPlacementConfigDetails, 0, size),
			},
			InitialNodeLabels: []containerengine.KeyValue{{Key: common.String("name"), Value: common.String(nodePoolName)}},
			NodeSourceDetails: containerengine.NodeSourceViaImageDetails{ImageId: common.String(imageId)},
		},
	}

	for i := 0; i < len(ads.Items); i++ {
		createNodePoolRequest.NodeConfigDetails.PlacementConfigs = append(createNodePoolRequest.NodeConfigDetails.PlacementConfigs, containerengine.NodePoolPlacementConfigDetails{
			AvailabilityDomain: ads.Items[i].Name,
			SubnetId:           subnet.Id,
		})
	}

	c, clerr := containerengine.NewContainerEngineClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clerr)

	ctx := context.Background()
	createNodePoolResponse, err := c.CreateNodePool(ctx, createNodePoolRequest)
	helpers.FatalIfError(err)
	fmt.Println("creating nodepool")

	WaitUntilWorkRequestComplete(c, createNodePoolResponse.OpcWorkRequestId)
	fmt.Println("nodepool created")
}

func WaitUntilWorkRequestComplete(client containerengine.ContainerEngineClient, workReuqestID *string) containerengine.GetWorkRequestResponse {
	// retry GetWorkRequest call until TimeFinished is set
	shouldRetryFunc := func(r common.OCIOperationResponse) bool {
		return r.Response.(containerengine.GetWorkRequestResponse).TimeFinished == nil
	}

	getWorkReq := containerengine.GetWorkRequestRequest{
		WorkRequestId:   workReuqestID,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(shouldRetryFunc),
	}

	getResp, err := client.GetWorkRequest(context.Background(), getWorkReq)
	helpers.FatalIfError(err)
	return getResp
}
