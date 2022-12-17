# OKE Node 自动弹性伸缩(Cluster Autoscaler)操作手册

    OKE完全兼容Kubernetes Cluster Autoscaler. OKE Cluster Autoscaler 根据应用程序工作负载需求,自动调整集群节点池的大小。同时，通过自动调整集群节点池的大小，可以确保应用程序可用性情况下，优化您的成本。

## 1. OKE Cluster Autoscaler 设计规划

### 1.1 OKE集群节点自动扩容或者缩放条件

    1. 扩容：由于资源不足，某些Pod无法在任何当前节点上进行调度.
    2. 缩容: Node节点资源利用率较低时，且此node节点上存在的pod都能被重新调度到其他node节点上运行.

### 1.2 什么时候OKE集群节点不会被Cluster Autoscaler删除

    1. 节点上有Pod被 PodDisruptionBudget 控制器限制。
    2. 节点上有命名空间是 kube-system 的pods。
    3. 节点上的Pod不是被控制器创建，例如不是被Deployment,Replica Set,Job,Stateful Set创建。
    4. 节点上有Pod使用了本地存储
    5. 节点上Pod驱逐后无处可去，即没有其他node能调度这个pod
    6. 节点有注解："cluster-autoscaler.kubernetes.io/scale-down-disabled": "true"(在CA 1.0.3或更高版本中受支持)

### 1.3 如何保证OKE集群节点不受OKE Cluster Autoscaler删除

    特定标签保护：从CA 1.0开始，节点可以打上以下标签：

    ```text
    "cluster-autoscaler.kubernetes.io/scale-down-disabled": "true"
    ```

    使用 kubectl 将其添加到节点(或从节点删除)：

    ```bash
    $<copy> kubectl annotate node cluster-autoscaler.kubernetes.io/scale-down-disabled=true </copy>
    ```

### 1.4 OKE Cluster Autoscaler 最佳实践

    1. 采用多个节点池(放在多个可用AD)，按节点池区划分功能，
    2. 同一个节点池中的节点:有相同的容量、CPU架构和操作系统版本。
    3. 不需要Cluster Autoscaler管理节点池，至少有一个节点池，不要把不需要自动调度的节点放到Cluster Autoscaler管理的节点池中.
    4. 不要手动修改Cluster Autoscaler管理的节点组及节点，容器CPU,内存配置。
    5. Pod定义多个副本，声明中带有 requests 资源限制。
    6. 对于不需要调度Pod,使用PodDisruptionBudgets属性，防止Pod删除。
    7. 设置节点池指定最小/最大之前，请检查您租户的配额是否足够大。
    8. 不要同一个集群中运行多个Cluster Autoscaler。
    9. Cluster Autoscaler自动调度max-node-provision-time 25分钟

    Kubernetes Cluster Autoscaler暂不支持参数：</br>
    *--node-group-auto-discovery : 不支持节点池自动发现. </br>
    *--node-autoprovisioning-enabled=true : Not supported.</br>
    *--gpu-total : 不支持GPU. </br>
    *--expander=price : 不支持根据成本自动调度. </br>

## 2. 赋予OKE Cluster Autoscaler 操作相关资源权限

### 2.1 创建dynamic group

1. 创建dynamic group

    a. Identity & Security. click Dynamic Groups.

    b. 选择所在的隔离区间（Compartment）

    创建dynamic group,输入dynamic group name，例如：oke-cluster-autoscaler-dyn-grp

    ```text
    <copy>
    ALL {instance.compartment.id = '<compartment-ocid>'} where <compartment-ocid>
    </copy>
    ```

    例如:

    ```text
    ALL {instance.compartment.id = 'ocid1.compartment.oc1..aaaaaaaa23______smwa'}
    ```

### 2.2 创建Policy策略，授权资源管理权限

1. 点击Identity & Security下面的Policies
2. 创建policy, 并给一个policy name
   例如： oke-cluster-autoscaler-dyn-grp-policy

    ```text
    <copy>
    Allow dynamic-group oke-cluster-autoscaler-dyn-grp to manage cluster-node-pools in compartment <compartment-name>
    Allow dynamic-group oke-cluster-autoscaler-dyn-grp to manage instance-family in compartment <compartment-name>
    Allow dynamic-group oke-cluster-autoscaler-dyn-grp to use subnets in compartment <compartment-name>
    Allow dynamic-group oke-cluster-autoscaler-dyn-grp to read virtual-network-family in compartment <compartment-name>
    Allow dynamic-group oke-cluster-autoscaler-dyn-grp to use vnics in compartment <compartment-name>
    Allow dynamic-group oke-cluster-autoscaler-dyn-grp to inspect compartments in compartment <compartment-name>
    </copy>
    ```

## 3. OKE Cluster Autoscaler部署

   为了实现自动缩放pod实现自动缩放，您需要部署Kubernetes Metrics Server(参见<https://nengbai.github.io/oke-dashboard/?lab=oke-metrics>)，以从集群中的每个工作节点收集资源度量。部署Kubernetes Metrics Server之后，您可以参照下面步骤部署OKE Cluster Autoscaler：

### 3.1 OKE Cluster Autoscaler配置

1. 下载 OKE Cluster Autoscaler cluster-autoscaler.yaml

    ```bash
    $ <copy> curl -o cluster-autoscaler.yaml https://github.com/nengbai/oke-dashborad/blob/main/cluster-autoscaler/cluster-autoscaler.yaml
    curl -o nginx.yaml https://github.com/nengbai/oke-dashborad/blob/main/cluster-autoscaler/nginx.yaml
    </copy>
    ```

2. 编辑 OKE Cluster Autoscaler配置 cluster-autoscaler.yaml
   
    ```bash
    $ <copy> vim cluster-autoscaler.yaml </copy>
    ---
    apiVersion: v1
    kind: ServiceAccount
    metadata:
    labels:
        k8s-addon: cluster-autoscaler.addons.k8s.io
        k8s-app: cluster-autoscaler
    name: cluster-autoscaler
    namespace: kube-system
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
    name: cluster-autoscaler
    labels:
        k8s-addon: cluster-autoscaler.addons.k8s.io
        k8s-app: cluster-autoscaler
    rules:
    - apiGroups: [""]
        resources: ["events", "endpoints"]
        verbs: ["create", "patch"]
    - apiGroups: [""]
        resources: ["pods/eviction"]
        verbs: ["create"]
    - apiGroups: [""]
        resources: ["pods/status"]
        verbs: ["update"]
    - apiGroups: [""]
        resources: ["endpoints"]
        resourceNames: ["cluster-autoscaler"]
        verbs: ["get", "update"]
    - apiGroups: [""]
        resources: ["nodes"]
        verbs: ["watch", "list", "get", "patch", "update"]
    - apiGroups: [""]
        resources:
        - "pods"
        - "services"
        - "replicationcontrollers"
        - "persistentvolumeclaims"
        - "persistentvolumes"
        verbs: ["watch", "list", "get"]
    - apiGroups: ["extensions"]
        resources: ["replicasets", "daemonsets"]
        verbs: ["watch", "list", "get"]
    - apiGroups: ["policy"]
        resources: ["poddisruptionbudgets"]
        verbs: ["watch", "list"]
    - apiGroups: ["apps"]
        resources: ["statefulsets", "replicasets", "daemonsets"]
        verbs: ["watch", "list", "get"]
    - apiGroups: ["storage.k8s.io"]
        resources: ["storageclasses", "csinodes","csidrivers","csistoragecapacities"]
        verbs: ["watch", "list", "get"]
    - apiGroups: ["batch", "extensions"]
        resources: ["jobs"]
        verbs: ["get", "list", "watch", "patch"]
    - apiGroups: ["coordination.k8s.io"]
        resources: ["leases"]
        verbs: ["create"]
    - apiGroups: ["coordination.k8s.io"]
        resourceNames: ["cluster-autoscaler"]
        resources: ["leases"]
        verbs: ["get", "update"]
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: Role
    metadata:
    name: cluster-autoscaler
    namespace: kube-system
    labels:
        k8s-addon: cluster-autoscaler.addons.k8s.io
        k8s-app: cluster-autoscaler
    rules:
    - apiGroups: [""]
        resources: ["configmaps"]
        verbs: ["create","list","watch"]
    - apiGroups: [""]
        resources: ["configmaps"]
        resourceNames: ["cluster-autoscaler-status", "cluster-autoscaler-priority-expander"]
        verbs: ["delete", "get", "update", "watch"]

    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
    name: cluster-autoscaler
    labels:
        k8s-addon: cluster-autoscaler.addons.k8s.io
        k8s-app: cluster-autoscaler
    roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: cluster-autoscaler
    subjects:
    - kind: ServiceAccount
        name: cluster-autoscaler
        namespace: kube-system

    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
    name: cluster-autoscaler
    namespace: kube-system
    labels:
        k8s-addon: cluster-autoscaler.addons.k8s.io
        k8s-app: cluster-autoscaler
    roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: Role
    name: cluster-autoscaler
    subjects:
    - kind: ServiceAccount
        name: cluster-autoscaler
        namespace: kube-system

    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
    name: cluster-autoscaler
    namespace: kube-system
    labels:
        app: cluster-autoscaler
    spec:
    replicas: 3
    selector:
        matchLabels:
        app: cluster-autoscaler
    template:
        metadata:
        labels:
            app: cluster-autoscaler
        annotations:
            prometheus.io/scrape: 'true'
            prometheus.io/port: '8085'
        spec:
        serviceAccountName: cluster-autoscaler
        containers:
            - image: iad.ocir.io/oracle/oci-cluster-autoscaler:{{ image tag }}
            name: cluster-autoscaler
            resources:
                limits:
                cpu: 100m
                memory: 300Mi
                requests:
                cpu: 100m
                memory: 300Mi
            command:
                - ./cluster-autoscaler
                - --v=4
                - --stderrthreshold=info
                - --cloud-provider=oci-oke
                - --max-node-provision-time=25m
                - --nodes=1:5:{{ node pool ocid 1 }}
                - --nodes=1:5:{{ node pool ocid 2 }}
                - --scale-down-delay-after-add=10m
                - --scale-down-unneeded-time=10m
                - --unremovable-node-recheck-timeout=5m
                - --balance-similar-node-groups
                - --balancing-ignore-label=displayName
                - --balancing-ignore-label=hostname
                - --balancing-ignore-label=internal_addr
                - --balancing-ignore-label=oci.oraclecloud.com/fault-domain
            imagePullPolicy: "Always"
            env:
            - name: OKE_USE_INSTANCE_PRINCIPAL
                value: "true"
            - name: OCI_SDK_APPEND_USER_AGENT
                value: "oci-oke-cluster-autoscaler"  
    ```

2. 增加 cluster-autoscaler containers command 特定参数
   a. 在cluster-autoscaler.yaml中增加cloud-provider为oci-oke

    ```text
    - --cloud-provider=oci-oke
    ```

   b.  对于Kubernetes version 1.23 or earlier 增加cloud-provider为oci

    ```text
    - --cloud-provider=oci
    ```

3. 调整cluster-autoscaler.yaml 容器Image Path 为 OCI Docker Registry

    参照<https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengusingclusterautoscaler.htm#Using_the_Kubernetes_Cluster_Autoscaler>

    ```text
    - image: iad.ocir.io/oracle/oci-cluster-autoscaler:{{ image tag }}
    ````

4. 增加Kubernetes Cluster Autoscaler管理的节点池

    ```text
    --nodes=<min-nodes>:<max-nodes>:<nodepool-ocid>
    ```

    例如：

    ```text
    - --nodes=1:5:{{ node pool ocid 1 }}
    ```

### 3.2 部署OKE Cluster Autoscaler

1. 部署 Kubernetes Cluster Autoscaler

    ```bash
    $<copy> kubectl apply -f cluster-autoscaler.yaml </copy>
    serviceaccount/cluster-autoscaler created
    clusterrole.rbac.authorization.k8s.io/cluster-autoscaler created
    role.rbac.authorization.k8s.io/cluster-autoscaler created
    clusterrolebinding.rbac.authorization.k8s.io/cluster-autoscaler created
    rolebinding.rbac.authorization.k8s.io/cluster-autoscaler created
    deployment.apps/cluster-autoscaler created
    ```

2. 检查Kubernetes Cluster Autoscaler日志信息

    ```bash
    $<copy> kubectl -n kube-system logs -f deployment.apps/cluster-autoscaler </copy>
    ```

3. 监控Kubernetes Cluster Autoscaler pods数量

    ```bash
    $<copy> kubectl -n kube-system get lease </copy>
    ```

4. 检查Kubernetes Cluster Autoscaler 状态

    ```bash
    $<copy> kubectl -n kube-system get cm cluster-autoscaler-status -o yaml </copy>
    ```

### 3.3 OKE Node 弹性自动扩容

1. 确定worker nodes数

    ```bash
    $<copy> kubectl get nodes  </copy>
    ```

2. 创建测试 Namespace 

    例如：nginx

    ```bash
    $<copy> kubectl create ns nginx </copy>
    ```

3. 为测试 Namespace 准备Secret key
    为了能安全正常从OCI Docker Registry拉取容器镜像，需要使用该集群OCI账号和 auth token 在OKE集群中该Namespace中增加Secret Key。例如：为Namespace kuboard 增加 Secret Key。

    ```bash
    <copy>kubectl create secret docker-registry ocisecret --docker-server=<icn.ocir.io> --docker-username='<oci username>' --docker-password='<auth token>' --docker-email='<email address>' -n kuboard </copy>
    ```

4. 下载image重定向到OCI Docker Registry
    a. 下载image

    ```bash
    $<copy> docker pull nginx:latest </copy>
    Trying to pull repository docker.io/library/nginx ... 
    latest: Pulling from docker.io/library/nginx
    025c56f98b67: Pull complete 
    ec0f5d052824: Pull complete 
    cc9fb8360807: Pull complete 
    defc9ba04d7c: Pull complete 
    885556963dad: Pull complete 
    f12443e5c9f7: Pull complete 
    Digest: sha256:75263be7e5846fc69cb6c42553ff9c93d653d769b94917dbda71d42d3f3c00d3
    Status: Downloaded newer image for nginx:latest
    nginx:latest
    ```

    b. Tag 容器镜像为 OCI Docker Registry, 例如：“icn.ocir.io/cnxcypamq98c/devops-repos”

    ```bash
    $<copy> docker tag docker.io/library/nginx:latest icn.ocir.io/cnxcypamq98c/devops-repos/nginx:latest </copy>
    ```
    c. 上传到OCI Docker Registry镜像库存储（例如：icn.ocir.io/cnxcypamq98c/devops-repos/）

    ```bash
    $ <copy> docker push icn.ocir.io/cnxcypamq98c/devops-repos/nginx:latest </copy>
    The push refers to repository [icn.ocir.io/cnxcypamq98c/devops-repos/nginx]
    e83791f03918: Pushed 
    10e506a84718: Pushed 
    9485bb85a132: Pushed 
    47064e0edc59: Pushed 
    5678f6b95362: Pushed 
    b5ebffba54d3: Pushed 
    latest: digest: sha256:d586384381a0e6834cef73d432b1486f0b86334cb92e54256def62dd403f82ab size: 1570
    ```

5. 调整 Nginx 应用部署文件nginx.yaml，注意调整 image path为 OCI Docker Registry镜像库存储路径（例如：icn.ocir.io/cnxcypamq98c/devops-repos/nginx:latest）

    ```text
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx-deployment
      namespace: nginx
    spec:
      selector:
        matchLabels:
        app: nginx
      replicas: 2
      template:
        metadata:
          labels:
            app: nginx
        spec:
          containers:
          - name: nginx
            image: icn.ocir.io/cnxcypamq98c/devops-repos/nginx:latest
            imagePullPolicy: IfNotPresent
            ports:
              - containerPort: 80
                resources:
                  requests:
                    memory: "500Mi"
          imagePullSecrets:
            - name: ocisecret
    ```

3. 部署Nginx 应用

    ```bash
    $<copy> kubectl create -f nginx.yaml </copy>
    deployment.apps/nginx-deployment created
    ```

4. 增加 deployment pods数量 从10 到100

    ```bash
    $<copy> kubectl get pod -n nginx｜wc -l 
    kubectl scale deployment nginx-deployment --replicas=100 -n nginx </copy>
    ```

    删除异常Pod:

    ```bash
    $<copy> kubectl -n nginx get pod|grep "0/1" |awk '{print $1}'|xargs -I {} kubectl -n nginx delete pod {} </copy>
    ```

5. 观察deployment 状态

    ```bash
    $<copy> kubectl get deployment nginx-deployment --watch -n nginx </copy>
    ```

6. 检查worker nodes 数量

    ```bash
    $<copy> kubectl get nodes </copy>
    ```

### 3.4 OKE Node 弹性自动缩容

1. 调整Nginx应用 deployment pods 数量

    ```bash
    $<copy> kubectl scale deployment nginx-deployment --replicas=10 -n nginx </copy>
    ```
    或者删除Nginx 应用：

    ```bash
    $<copy> kubectl delete deployment nginx-deployment -n nginx </copy>
    ```

2. 检查Kubernetes Cluster Autoscaler日志信息

    ```bash
    $<copy> kubectl -n kube-system logs -f deployment.apps/cluster-autoscaler </copy>
    I1215 08:16:23.104671       1 pre_filtering_processor.go:57] Node 10.0.10.73 should not be processed by cluster autoscaler (no node group config)
    I1215 08:16:23.104725       1 scale_down.go:448] Node 10.0.10.211 - cpu utilization 0.146000
    I1215 08:16:23.105208       1 scale_down.go:509] Scale-down calculation: ignoring 3 nodes unremovable in the last 5m0s
    I1215 08:16:23.105339       1 static_autoscaler.go:522] 10.0.10.211 is unneeded since 2022-12-15 08:12:40.621793266 +0000 UTC m=+5355.825155070 duration 3m42.479860286s
    I1215 08:16:23.105389       1 static_autoscaler.go:533] Scale down status: unneededOnly=false lastScaleUpTime=2022-12-15 07:21:50.551646995 +0000 UTC m=+2305.755008799 lastScaleDownDeleteTime=2022-12-15 08:12:30.343455121 +0000 UTC m=+5345.546816965 lastScaleDownFailTime=2022-12-15 05:43:47.576028691 +0000 UTC m=-3577.220609485 scaleDownForbidden=false isDeleteInProgress=false scaleDownInCooldown=false
    I1215 08:16:23.105438       1 static_autoscaler.go:546] Starting scale down
    I1215 08:16:23.105503       1 scale_down.go:828] *10.0.10.211 was unneeded for *3m42.479860286s
    I1215 08:16:23.105563       1 scale_down.go:917] No candidates for scale down
    ```
3. 确认worker nodes数量减少到初始数量(需要等待，大约在25分钟左右，优先删除是autoscaler增加的节点)

    ```bash
    $<copy> kubectl get nodes </copy>
    NAME          STATUS                     ROLES   AGE    VERSION
    10.0.10.12    Ready                      node    58d    v1.24.1
    10.0.10.211   Ready,SchedulingDisabled   node    59m    v1.24.1
    10.0.10.240   Ready                      node    7h2m   v1.24.1
    10.0.10.37    Ready                      node    7h2m   v1.24.1
    10.0.10.68    Ready                      node    58d    v1.24.1
    10.0.10.73    Ready                      node    58d    v1.24.1
    10.0.10.88    Ready                      node    7h2m   v1.24.1
    ```
