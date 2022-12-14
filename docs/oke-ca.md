# OKE Node 自动弹性伸缩(Cluster Autoscaler)操作手册

OKE完全兼容Kubernetes Cluster Autoscaler. OKE Cluster Autoscaler 根据应用程序工作负载需求,自动调整集群节点池的大小。同时，通过自动调整集群节点池的大小，可以确保应用程序可用性情况下，优化您的成本。

## 1. OKE Cluster Autoscaler 设计规划

### 1.1 集群自动扩容或者缩放条件

1. 扩容：由于资源不足，某些Pod无法在任何当前节点上进行调度.
2. 缩容: Node节点资源利用率较低时，且此node节点上存在的pod都能被重新调度到其他node节点上运行.

### 1.2 什么时候集群节点不会被 CA 删除

1. 节点上有Pod被 PodDisruptionBudget 控制器限制。
2. 节点上有命名空间是 kube-system 的pods。
3. 节点上的Pod不是被控制器创建，例如不是被Deployment,Replica Set,Job,Stateful Set创建。
4. 节点上有Pod使用了本地存储
5. 节点上Pod驱逐后无处可去，即没有其他node能调度这个pod
6. 节点有注解："cluster-autoscaler.kubernetes.io/scale-down-disabled": "true"(在CA 1.0.3或更高版本中受支持)

### 1.3 如何保证节点不受OKE Cluster Autoscaler删除

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

Kubernetes Cluster Autoscaler暂不支持参数：
*--node-group-auto-discovery : 不支持节点池自动发现
*--node-autoprovisioning-enabled=true : Not supported.
*--gpu-total : 不支持GPU.
*--expander=price : 不支持根据成本自动调度.

## 2. OCI OKE Cluster Autoscaler 权限

### 2.1 创建dynamic group

1. 创建dynamic group
    a. Identity & Security. Under Identity, click Dynamic Groups.
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

1. 定制 OKE Cluster Autoscaler配置 cluster-autoscaler.yaml

    ```text
    <copy>
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
        resources: ["storageclasses", "csinodes"]
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
    </copy>
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

2. 定义一个 Nginx 应用

    ```text
    <copy>
    apiVersion: apps/v1
        kind: Deployment
        metadata:
            name: nginx-deployment
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
                image: nginx:latest
                ports:
            - containerPort: 80
                resources:
                requests:
                    memory: "500Mi"
    </copy>
    ```

3. 部署Nginx 应用

    ```bash
    $<copy> kubectl create -f nginx.yaml </copy>
    ```

4. 增加 deployment pods数量 从10 到100

    ```bash
    $<copy> kubectl scale deployment nginx-deployment --replicas=100 </copy>
    ```

5. 观察deployment 状态

    ```bash
    $<copy> kubectl get deployment nginx-deployment --watch </copy>
    ```

6. 检查worker nodes 数量

    ```bash
    $<copy> kubectl get nodes </copy>
    ```

### 3.4 OKE Node 弹性自动缩容

1. 调整Nginx应用 deployment pods 数量

    ```bash
    $<copy> kubectl scale deployment nginx-deployment --replicas=10 </copy>
    ```
    或者删除Nginx 应用：

    ```bash
    $<copy> kubectl delete deployment nginx-deployment </copy>
    ```

2. 确认worker nodes数量减少到初始数量

    ```bash
    $<copy> kubectl get nodes </copy>
    ```
