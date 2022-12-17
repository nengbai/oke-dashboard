# OKE Dashboard for above version 1.24

Kubernetes 从v1.24.0 开始使用的安装方式较之前有很大不同，OKE 容器服务完全兼容Kubernetes，当使用版本v1.24.0 or above v1.24.0 版本，需要参考该文档。

## 1. OKE Dashboard环境准备

### 1.1 下载OKE Dashboard 安装包

下载前，请确认下载机器可以访问https://github.com.

```bash
$ <copy> curl -o recommended.yaml https://github.com/nengbai/oke-dashborad/blob/main/dashboard/recommended.yaml
curl -o oke-admin.yaml  https://github.com/nengbai/oke-dashborad/blob/main/dashboard/oke-admin.yaml
curl -o dashboard-ingress.yaml https://github.com/nengbai/oke-dashborad/blob/main/dashboard/dashboard-ingress.yaml
curl -o create_cert.sh https://github.com/nengbai/oke-dashborad/blob/main/dashboard/create_cert.sh </copy>
```

### 1.2. 生成证书

1. 根据域名生成证书
注意：用您自己域名替换 example.com

``` bash
$ <copy> bash create_cert.sh oke-admin dashboard example.com kubernetes-dashboard </copy> 
```

2. 您会看到自动生成证书，并在您的OKE集群创建Secret

```bash
$<copy>kubectl get secret -n kubernetes-dashboard|grep oke </copy>
oke-admin    kubernetes.io/service-account-token   4      56d
```

## 2.部署 OKE Dashboard
### 2.1 安装OKE Dashboard

1. 执行安装OKE Dashboard

``` bash
$ <copy>  kubectl apply -f recommended.yaml </copy> 
```

2. 检查 OKE Dashboard状态

``` bash
$ <copy> kubectl -n kubernetes-dashboard get pod,svc </copy> 
NAME                                            READY   STATUS    RESTARTS   AGE
pod/dashboard-metrics-scraper-8c47d4b5d-xt8h4   1/1     Running   0          56d
pod/kubernetes-dashboard-67bd8fc546-2d6rt       1/1     Running   0          56d

NAME                                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/dashboard-metrics-scraper   ClusterIP   10.96.201.40    <none>        8000/TCP   56d
service/kubernetes-dashboard        ClusterIP   10.96.161.240   <none>        443/TCP    56d
```

### 2.2 创建OKE 用户及其资源授权

``` bash
$ <copy>  kubectl apply -f oke-admin.yaml </copy> 
```

## 3. 验证
### 3.1 本地验证

1. 启动本地Proxy

在您本地终端执行下面命令（需要kubectl环境）
    ``` bash
    $ <copy> kubectl proxy </copy> 
    Starting to serve on 127.0.0.1:8001
    ````


2. 获取登录token.

本地在您本地终端执行下面命令，并记录获取到"token"：

``` bash
$ <copy> curl 127.0.0.1:8001/api/v1/namespaces/kubernetes-dashboard/serviceaccounts/oke-admin/token -H "Content-Type:application/json" -X POST -d '{}' </copy> 
example:
 "token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImN0czdrSU1zcVo0c1R3YU9mTzZMMHdjY2JydTdJekt5dXdpQ1Z2d3lRbncifQ.eyJhdWQiOlsiYXBpIl0sImV4cCI6MTY2NjA2ODgyMiwiaWF0IjoxNjY2MDY1MjIyLCJpc3MiOiJodHRwczovL2t1YmVybmV0ZXMuZGVmYXVsdC5zdmMuY2x1c3Rlci5sb2NhbCIsImt1YmVybmV0ZXMuaW8iOnsibmFtZXNwYWNlIjoia3ViZXJuZXRlcy1kYXNoYm9hcmQiLCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoib2tlLWFkbWluIiwidWlkIjoiMmI0NGMxM2QtNzBkNS00MWI3LTk5MTUtNzQ2MjQxMDFkYzBlIn19LCJuYmYiOjE2NjYwNjUyMjIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlcm5ldGVzLWRhc2hib2FyZDpva2UtYWRtaW4ifQtExqrAwIJJ8WRFrNbH4BnSUDK2P0XBAizJafruSfBgksh_ivJrj6TzaTk1UgY6zFfw_fGQ9mB5nWMLVR1yMHTFpAjsUfnEoFU5alv2MBFVJ5mPGBhznoDVi7ZdU29hKr6LLUr2EbOWVHPkeLFtjivGe38S9wpzaL8iGN_bPtV2usJt8ZoYoVtc-jy0stPDm-2idi5aonAjqwKfsyS75WpLdq8Gx25Ge3Rw64diUo5-WgA3aSng7BGvWGR4FWvLKN3VKVVEuyzb5wmIcMb8MAOko1C8lYvma-L0OHA87DmOFGAo1GHQf7O8dtjBjCqVsnWA"
```

3. 网页验证

在您本地终端电脑上(与上面2步骤执行同一台电脑)，用浏览器(firefox or chrome)中打开下面网址

``` text
<copy> http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ </copy>
```

### 3.2 增加外网访问方式

前提是: 已经在OKE集群中部署OKE Ingress控制器,如果没有准备好OKE Ingress控制器，请参照下面link部署：</br>
<https://minqiaowang.github.io/oci-k8s-cn/workshops/freetier/?lab=deploy-complex-app#Task1:OKEIngress>

1. 增加 OKE Dashboard ingress 访问方式

Dashboard缺省只能通过Proxy转发本地访问（第6和第7步骤），为了方便管理，需要给OKE Dashboard增加https访问的Ingress.</br>

调整访问域名,编辑dashboard-ingress.yaml， 替换 example.com 为您在第2步 生成证书时对应的域名。

```bash
$ <copy> vim dashboard-ingress.yaml </copy>
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: oke-dashboard-ingress
  namespace: kubernetes-dashboard
  annotations:
    # 开启use-regex，启用path的正则匹配
    nginx.ingress.kubernetes.io/use-regex: "true"
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - oke-dashboard.example.com
    secretName: oke-admin
  rules:
  - host: "oke-dashboard.example.com"
    http:
      paths:
      - pathType: Prefix
        path: "/"
        backend:
          service:
            name: kubernetes-dashboard
            port:
              number: 443
```

2. 增加 Ingress

``` 
$ <copy> kubectl apply -f dashboard-ingress.yaml </copy> 
```

3. 外网访问验证 OKE Dashboard

浏览器(firefox or chrome)中打开下面网址

``` text
    https://oke-dashboard.example.com
```
