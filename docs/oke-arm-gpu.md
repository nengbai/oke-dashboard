# OKE 适配Arm-based 和 GPU Nodes

## 1. OKE Arm-based 架构高并发么应用支持

OKE version 1.19.7开始支持ARM-base 架构高并发么应用，需要为OKE集群定义特定ARM-base 节点池，且选择ARM-base 类型实例和支持ARM-base 操作系统的image.

1. 注意项目如下

    * 定义节点池，增加Arm-based 类型实例节点。
    * Arm-based 类型实例节点需要Kubernetes version 1.19.7 或更高。
    * 节点节点池中不能混合Arm-based类型实例和非Arm-based类型实例。
    * Arm-based 类型实例节点需要选择兼容Oracle Linux ARM-base操作系统image。

2. 定义运行在Arm-based nodes的pod

    ```bash

    apiVersion: v1
    kind: Pod
    metadata:
    name: nginx
    labels:
        env: test
    spec:
    containers:
    - name: nginx
        image: nginx
        imagePullPolicy: IfNotPresent
    nodeSelector:
        kubernetes.io/arch: arm64
    ```

## 2. OKE GPU 支持

OKE version 1.19.7开始支持GPU 功能，需要为OKE集群定义特定GPU节点池，且选择GPU 类型实例和支持GPU 操作系统的image.

1. 注意项目如下

    * 定义节点池，增加GPU 类型实例节点。
    * GPU 类型实例节点需要Kubernetes version 1.19.7 或更高。
    * 节点节点池中不能混合GPU 类型实例和非GPU 类型实例。
    * GPU 类型实例节点需要选择兼容Oracle Linux GPU操作系统image。
    * 不是所有可用区域都有GPU 类型实例节点。

2. 当应用程序运行在GPU 类型OKE 节点上，下面会挂在到pod

* requested number of GPU devices.
* node's CUDA library

3. 定义运行在GPU类型Node的Pod

    ```bash
    apiVersion: v1
    kind: Pod
    metadata:
    name: test-with-gpu-workload
    spec:
    restartPolicy: OnFailure
    containers:
        - name: cuda-vector-add
        image: k8s.gcr.io/cuda-vector-add:v0.1
        resources:
            limits:
            nvidia.com/gpu: 1
    ```

4. 定义运行在No-GPU类型Node的Pod

    ```bash

    apiVersion: v1
    kind: Pod
    metadata:
    name: test-with-non-gpu-workload
    spec:
    restartPolicy: OnFailure
    containers:
        - name: test-with-non-gpu-workload
        image: "oraclelinux:8"
    ```
