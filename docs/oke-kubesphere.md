# KubeSphere on Oracle OKE

## $1. KubeSphere系统要求

1. KubeSphere 3.3 要求Kubernetes version 满足 v1.19.x, v1.20.x, v1.21.x, *v1.22.x,*v1.23.x, and * v1.24.x. 最优选择是 Kubernetes v1.21.x 或更高版本。
    ```bash
    $ <copy>kubectl version </copy>
    Client Version: version.Info{Major:"1", Minor:"25", GitVersion:"v1.25.4", GitCommit:"872a965c6c6526caa949f0c6ac028ef7aff3fb78", GitTreeState:"clean", BuildDate:"2022-11-09T13:36:36Z", GoVersion:"go1.19.3", Compiler:"gc", Platform:"darwin/amd64"}
    Kustomize Version: v4.5.7
    Server Version: version.Info{Major:"1", Minor:"24", GitVersion:"v1.24.1", GitCommit:"b13b16197b0e07f78f7ced71255ce69516fdd9e6", GitTreeState:"clean", BuildDate:"2022-05-30T10:16:45Z", GoVersion:"go1.18.2 BoringCrypto", Compiler:"gc", Platform:"linux/amd64"}
    ```

2. 可用资源 CPU > 1 Core 和 Memory > 2 G. x86_64 CPUs are supported, 当前不支持Arm CPUs。

    ```bash
    $ <copy> kubectl get nodes </copy>
    NAME         STATUS   ROLES   AGE   VERSION
    10.0.10.12   Ready    node    55d   v1.24.1
    10.0.10.68   Ready    node    55d   v1.24.1
    10.0.10.73   Ready    node    55d   v1.24.1
    ```

3. Kubernetes cluster有 default StorageClass。

    ```bash
    $ <copy> kubectl get sc </copy>
    NAME              PROVISIONER            RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
    oci               oracle.com/oci         Delete          Immediate               false                 55d
    oci-bv (default)  blockvolume.csi.oraclecloud.com   Delete  WaitForFirstConsumer   true                55d
    ```

## $2. 安装 KubeSphere on OKE

1. 安装 KubeSphere install初始化
   注意：下面<v3.3.1>需要根据OKE版本要求替换成对应版本，参见$1. KubeSphere系统要求.

    ```bash
    $ <copy> kubectl apply -f https://github.com/kubesphere/ks-installer/releases/download/v3.3.1/kubesphere-installer.yaml </copy>
    ```

2. 部署 KubeSphere

    kubectl apply -f https://github.com/kubesphere/ks-installer/releases/download/v3.3.1/cluster-configuration.yaml

3. 监控 KubeSphere状态

    ```bash
    $ <copy>kubectl logs -n kubesphere-system $(kubectl get pod -n kubesphere-system -l 'app in (ks-install, ks-installer)' -o jsonpath='{.items[0].metadata.name}') -f </copy>

    When the installation finishes, you can see the following message:

    #####################################################
    ###              Welcome to KubeSphere!           ###
    #####################################################

    Console: <http://10.0.10.12:30880>
    Account: admin
    Password: P@88w0rd

    NOTES：
      1. After logging into the console, please check the
        monitoring status of service components in
        the "Cluster Management". If any service is not
        ready, please wait patiently until all components 
        are ready.
      2. Please modify the default password after login.

    #####################################################
    https://kubesphere.io             20xx-xx-xx xx:xx:xx

    Access KubeSphere Console
    ```
4. 检查 KubeSphere svc状态

    kubectl get svc -n kubesphere-system

## $3. 部署 KubeSphere Ingress

1. 下载kubesphere-ingress.yaml

    ```bash
    $ <copy> curl -o kubesphere-ingress.yaml https://github.com/nengbai/oke-dashborad/blob/main/kubesphere/kubesphere-ingress.yaml </copy> 
    ```
2. 编辑 kubesphere-ingress.yaml,调整域名:example.com 为您拥有域名
3. 部署 kubesphere ingress

    ```bash
    $ <copy> kubectl apply -f kubesphere-ingress.yaml </copy> 
    ingress.networking.k8s.io/oke-kubesphere-ingress created
    ```
4. 检查Ingress状态

    ```bash 
    $ <copy> kubectl -n kubesphere-sys get ing </copy> 
    NAME                     CLASS   HOSTS                        ADDRESS             PORTS      AGE
    oke-kubesphere-ingress   nginx   oke-kubesphere.example.com   141.147.172.67        80       2m44s
    ```

## $3. 验证

1. 增加域名解释
   长期使用建议使用dns服务解释，如果是临时测试，建议在本地hosts中增加，下面以mac中增加域名解释为例。

    ```bash
    $ <copy> sudo vi /etc/hosts</copy> 
    141.147.172.67  oke-kubesphere.example.com
    ```
2. 浏览器访问 Kubesphere 验证
    在浏览器中打开链接<http://your-ingress>
    例如： <http://oke-kubesphere.example.com>
    输入初始用户名和密码，并登录
        用户名： admin
        密码： P@88w0rd
    第一次登录需要修改新密码。
