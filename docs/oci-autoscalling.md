# OCI AutoScalling 基于系统负载实现自动扩容或缩容

本操作是在有VCN 网络环境操作情况实施，如果没VCN，请参照VCN 章节创建VCN网络环境。

## 1. 创建负载均衡器(Load Balancer)和增加安全策略(Security List）

### 1.1 创建负载均衡器（Load Balancer)

1. 在OCI Services menu下的 Networking中 选择 Load Balancers

2. 创建一个负载均衡器Load Balancer

    详细参数参照如下：</br>
    * LOAD BALANCER NAME: 输入load balancer名字</br>
    * CHOOSE VISIBILITY TYPE: Public  选择类型为Public</br>
    * CHOOSE THE MAXIMUM TOTAL BANDWIDTH: 选择带款 SMALL 100Mbps， 可选择10Mbps～8Gbps之间</br>

    VIRTUAL CLOUD NETWORK: 选择您的 Virtual Cloud Network </br>
    SUBNET: Choose the Public Subnet 选择一个您的公有子网 </br>

3. 增加 Backends 应用配置

   此处暂时不用增加，在2.1 配置 VM instance pool 章节配置。

4. 增加健康监测策略：HEALTH CHECK POLICY

    ```text
    PROTOCOL:         HTTP
    Port:             80 
    URL PATH (URI):    /
    ```

   
5. 配置监听器-Configure Listener

    ```text
    选择 LISTENER 协议: HTTP
    输入 LISTENER Ingress Port: 80
    ```
### 1.2 增加安全策略（Security List)

在OCI Services menu下的 Networking->Virtual Cloud Networks，

1. 点击对应的VCN name， 进入VCN Details
2. 在左下方的中点击 Security Lists， 选择对应的 Security List
3. 增加 Ingress 规则 和egress 规则

    ```text
    Source Type:            CIDR
    Source CIDR:            0.0.0.0/0
    IP Protocol:            TCP 
    Source Port Range:       All
    Destination Port Range:  80 (the listener port)
    ```

## 2. 配置 VM instance pool 和 auto scaling 策略

## 2.1 配置 VM instance pool

1. 在OCI services menu->Compute 下 点击 Instances

2. 点击创建 Instance.  

    填写下面信息:

    a. Name your instance: 实例名称

        Create in Compartment: 选择 compartment

    b. Placement：

        Availability Domain: 选这一个可用AD (缺省AD 1)

    c. Image and shape: 选择采用的操作系统，推荐使用最新版本 Oracle Linux available

3. 选择实例 Network 和Storage（可选项）:

    a. 选择您comparment下的VCN网络：Virtual cloud network in devops

    b. 选择IP类型： Public IPv4 address 

    c. Add SSH Keys：Upload public key files (.pub) </br>

    d. Boot volume </br>
        *Specify a custom boot volume size（可选项）:缺省系统盘50GB </br>
        *Use in-transit encryption：缺省选用Oracle 加密算法加密      </br>
        *Encrypt this volume with a key that you manage： 或选择您拥有加密算法加密 </br>

    e. Show advanced options

    在表格的Management项中，增加cloud-init script:

    ```text
    #cloud-config
    packages:
    - httpd
    - stress

    runcmd:
    - [sh, -c, echo "<html>Web Server IP `hostname --ip-address`</html>" > /var/www/html/index.html]
    - [firewall-offline-cmd, --add-port=80/tcp]
    - [systemctl, start, httpd]
    - [systemctl, restart, firewalld]
    ```

3. 点击"Create"

### 2.2 创建 VM Instance Pool

1. 点击 Instance name，进入Instance details页面

2. actions 中选择"Create Instance Configuretion"

3. 配置实例基础信息

4. 配置实例池

### 2.3 auto scaling 策略

1. Autoscaling configurations

## 3: 验证测试

自动弹性伸/缩策略是：当CPU 使用率在300秒内一直大于 80%，将自动创建1个实例. 一旦CPU 使用率在300秒内一直小于40%，将自动删除1个实例，实例池中至少有一个实例。

1. 登录VM 实例

    ```bash
    ssh -i <private_key> opc@<PUBLIC_IP_OF_COMPUTE>
    ```

2. 启动 CPU stress，进行压力测试

    ```bash
    sudo stress --cpu 4 --timeout 350
    ```

3. 检查CPU使用率：在Instance Pool Details下，选择VM instance Name

    当CPU utilization > 80%，自动增加1个VM 实例；当CPU utilization < 40%,从实例池中减少1个VM 实例。
