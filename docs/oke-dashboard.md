# oke dashboard for above version 1.24

Kubernetes 从v1.24.0 开始使用的安装方式于之前很很大不同，OKE 容器服务完全兼容Kubernetes，当使用版本v1.24.0 or above v1.24.0 版本，需要参考该文档。

## 1. 下载OKE Dashboard 安装包

```bash
$ <copy> curl -o recommended.yaml https://github.com/nengbai/oke-dashborad/blob/main/dashboard/recommended.yaml
curl -o oke-admin.yaml  https://github.com/nengbai/oke-dashborad/blob/main/dashboard/oke-admin.yaml
curl -o dashboard-ingress.yaml https://github.com/nengbai/oke-dashborad/blob/main/dashboard/dashboard-ingress.yaml
curl -o create_cert.sh https://github.com/nengbai/oke-dashborad/blob/main/dashboard/create_cert.sh </copy>
```

## 2.根据域名生成证书

    注意：用自己域名替换 example.com
``` bash
$ <copy> bash create_cert.sh oke-admin dashboard example.com kubernetes-dashboard </copy> 
```

## 3.部署 OKE Dashboard

``` bash
$ <copy>  kubectl apply -f recommended.yaml </copy> 
```

## 4. 检查 OKE Dashboard

``` bash
$ <copy> kubectl -n kubernetes-dashboard get pod,svc </copy> 
```

## 5. 创建OKE 用户及其资源授权

``` bash
$ <copy>  kubectl apply -f oke-admin.yaml </copy> 
```

## 6. 启动本地Proxy

``` bash
$ <copy> kubectl proxy </copy> 
````

## 7. 浏览器(firefox or chrome)中打开下面网址

``` text
<copy> http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ </copy>
```

## 8. 获取登录 token

``` bash
$ <copy> curl 127.0.0.1:8001/api/v1/namespaces/kubernetes-dashboard/serviceaccounts/oke-admin/token -H "Content-Type:application/json" -X POST -d '{}' </copy> 
        example:
 "token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImN0czdrSU1zcVo0c1R3YU9mTzZMMHdjY2JydTdJekt5dXdpQ1Z2d3lRbncifQ.eyJhdWQiOlsiYXBpIl0sImV4cCI6MTY2NjA2ODgyMiwiaWF0IjoxNjY2MDY1MjIyLCJpc3MiOiJodHRwczovL2t1YmVybmV0ZXMuZGVmYXVsdC5zdmMuY2x1c3Rlci5sb2NhbCIsImt1YmVybmV0ZXMuaW8iOnsibmFtZXNwYWNlIjoia3ViZXJuZXRlcy1kYXNoYm9hcmQiLCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoib2tlLWFkbWluIiwidWlkIjoiMmI0NGMxM2QtNzBkNS00MWI3LTk5MTUtNzQ2MjQxMDFkYzBlIn19LCJuYmYiOjE2NjYwNjUyMjIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlcm5ldGVzLWRhc2hib2FyZDpva2UtYWRtaW4ifQ.YA8sU6gyW7sTWHBoHO9jtExqrAwIJJ8WRFrNbH4BnSUDK2P0XBAizJafruSfBgksh_ivJrj6TzaTk1UgY6zFfw_fGQ9mB5nWMLVR1yMHTFpAjsUfnEoFU5alv2MBFVJ5mPGBhznoDVi7ZdU29hKr6LLUr2EbOWVHPkeLFtjivGe38S9wpzaL8iGN_bPtV2usJt8ZoYoVtc-jy0stPDm-2idi5aonAjqwKfsyS75WpLdq8Gx25Ge3Rw64diUo5-WgA3aSng7BGvWGR4FWvLKN3VKVVEuyzb5wmIcMb8MAOko1C8lYvma-L0OHA87DmOFGAo1GHQf7O8dtjBjCqVsnWA"
```

## 9. 增加 OKE Dashboard ingress 访问方式

``` 
$ <copy> kubectl apply -f dashboard-ingress.yaml </copy> 
```

## 10. 验证 OKE Dashboard

浏览器(firefox or chrome)中打开下面网址

``` text
    https://oke-dashboard.example.com
```
