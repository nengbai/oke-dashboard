# OKE 证书管理器(OKE cert-manager)
cert-manager 将证书和证书颁发者作为资源类型添加到 Kubernetes 集群中，并简化了这些证书的获取、更新和使用过程。
它支持从各种来源颁发证书，包括 Let's Encrypt (ACME)、HashiCorp Vault 和 Venafi TPP / TLS Protect Cloud，以及本地集群内颁发。
cert-manager 还确保证书保持有效和最新，尝试在到期前的适当时间更新证书以降低中断风险并消除工作量。
![cert-manager high level overview diagram](https://cert-manager.io/images/high-level-overview.svg)

## 文档使用说明

可以在【cert-manager.io】(https://cert-manager.io/docs/) 找到 cert-manager 的文档。
有关为 Ingress 资源自动颁发 TLS 证书的常见用例，请参阅【cert-manager nginx-ingress】(https://cert-manager.io/docs/tutorials/acme/nginx-ingress/) 快速入门指南。
有关颁发您的第一个证书的更全面的指南，请参阅我们的入门指南(https://cert-manager.io/docs/getting-started/)。

### 安装文档

安装支持多种支持的方法：

1.默认静态安装
您不需要对 cert-manager 安装参数进行任何调整。默认静态配置可以安装如下：
Install kubectl version >= v1.19.0. (otherwise, you'll have issues updating the CRDs - see v0.16 upgrade notes)

```bash
$ <copy>kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml </copy>
```
默认情况下，cert-manager 将安装到cert-manager 命名空间中。可以在不同的命名空间中运行 cert-manager，但您需要对部署清单进行修改。

2. 持续部署(CI/CD)
您知道如何配置您的证书管理器设置并希望将其自动化。
📖 helm：直接将cert-manager Helm chart与 Flux、ArgoCD 和 Anthos 等系统一起使用。
📖 helm 模板：用helm template生成自定义的证书管理器安装清单。有关详细信息，请参阅使用 helm 模板输出 YAML 。这个模板化的证书管理器清单可以通过管道传输到您首选的部署工具中。

### 验证
检查证书管理器 API
1. 首先，确保安装了 cmctl。
cmctl 对 Kubernetes 集群执行试运行证书创建检查。The cert-manager API is ready如果成功，将显示消息。
```
$ <copy> cmctl check api </copy>
The cert-manager API is ready
```

该命令也可用于等待检查成功。这是在安装 cert-manager 的同时运行命令的输出示例：
```
$ <copy> cmctl check api --wait=2m </copy>
Not ready: the cert-manager CRDs are not yet installed on the Kubernetes API server
Not ready: the cert-manager CRDs are not yet installed on the Kubernetes API server
Not ready: the cert-manager webhook deployment is not ready yet
Not ready: the cert-manager webhook deployment is not ready yet
Not ready: the cert-manager webhook deployment is not ready yet
Not ready: the cert-manager webhook deployment is not ready yet
The cert-manager API is ready
```
cert-manager、cert-manager-cainjector和 cert-manager-webhookpod 处于某种Running状态。webhook 可能需要比其他人更长的时间才能成功配置。

2. 创建一个Issuer测试 webhook 是否正常工作

```
$ <copy> cat <<EOF > test-resources.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cert-manager-test
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: test-selfsigned
  namespace: cert-manager-test
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: selfsigned-cert
  namespace: cert-manager-test
spec:
  dnsNames:
    - example.com
  secretName: selfsigned-cert-tls
  issuerRef:
    name: test-selfsigned
EOF
</copy>
```

创建测试资源:

```
$ <copy> kubectl apply -f test-resources.yaml </copy>
```

检查新创建的证书的状态。在 cert-manager 处理证书请求之前，您可能需要等待几秒钟。

```
$ <copy> kubectl describe certificate -n cert-manager-test </copy>

---
Spec:
  Common Name:  example.com
  Issuer Ref:
    Name:       test-selfsigned
  Secret Name:  selfsigned-cert-tls
Status:
  Conditions:
    Last Transition Time:  2019-01-29T17:34:30Z
    Message:               Certificate is up to date and has not expired
    Reason:                Ready
    Status:                True
    Type:                  Ready
  Not After:               2019-04-29T17:34:29Z
Events:
  Type    Reason      Age   From          Message
  ----    ------      ----  ----          -------
  Normal  CertIssued  4s    cert-manager  Certificate issued successfully
```

3. 清理环境

```
$ <copy> kubectl delete -f test-resources.yaml </copy>
```

## 故障排除
如果您在使用 cert-manager 时遇到任何问题，我们可以通过多种方式获得帮助：
我们网站上的[故障排除指南](https://cert-manager.io/docs/faq/troubleshooting/) 。
我们的官方[Kubernetes Slack频道](https://cert-manager.io/docs/contributing/#slack) - 最快的提问方式！
搜索现[已知问题](https://github.com/cert-manager/cert-manager/issues)。
如果你认为你已经找到了一个错误并且找不到现有的问题，请随时打开一个新问题！请务必包含尽可能多的关于您的环境的信息。

### 社区
Google cert-manager-devGroup 用于项目范围内的公告和开发协调。任何人都可以通过访问[此处](https://groups.google.com/forum/#!forum/cert-manager-dev) 并单击“加入群组”来加入群组。加入该群组需要一个 Google 帐户。
