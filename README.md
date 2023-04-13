# 项目说明
## 关于项目
当前项目是一个go-micro的register插件，当我们在用k8s部署go-micro的grpc server的时候，如果用k8s 内置的service
去部署grpc server的pod，是无法进行grpc的http2.0的，所以我采用headless service(无头service)方式去部署，把service当作
dns用，利用go的net包去做dns解析，然后获取service endpoints的podip记录返回，实现服务发现的功能。


## 如何使用
```go
package main

import (
	"fmt"
	"github.com/go-micro-v4-demo/frontend/handler"
	helloworldPb "github.com/go-micro-v4-demo/helloworld/proto"
	userPb "github.com/go-micro-v4-demo/user/proto"
	mgrpc "github.com/go-micro/plugins/v4/client/grpc"
	mhttp "github.com/go-micro/plugins/v4/server/http"
	"github.com/gorilla/mux"
	k8sHeadlessSvc "github.com/gsmini/k8s-headless-svc"
	"go-micro.dev/v4/logger"
	"net/http"
)

var (
	service = "frontend"
	version = "latest"
)

const K8sSvcName = "user-svc"

const UserSvcName = "user-svc"        //user模块在k8s service中的metadata.name的名字
const HelloWordSvcName = "helloworld" //user模块在k8s service中的metadata.name的名字
func main() {
	UserSvc := &k8sHeadlessSvc.Service{Namespace: "default", SvcName: UserSvcName, PodPort: 8080}
	//HelloWordSvc := &k8sHeadlessSvc.Service{Namespace: "default", SvcName: HelloWordSvcName, PodPort: 9090}
	reg := k8sHeadlessSvc.NewRegistry([]*k8sHeadlessSvc.Service{UserSvc})
	// 当前frontend调用依赖多个grpc 上游服务器的情况
	//reg := k8sHeadlessSvc.NewRegistry([]*k8sHeadlessSvc.Service{UserSvc},[]*k8sHeadlessSvc.Service{HelloWordSvcName})
	srv := micro.NewService()
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Address("0.0.0.0:8080"), 
		micro.Registry(reg),//和内置registry一样注册我们的k8sHeadlessSvc registry
	)
	# 省略不重要业务到吗
}
```

## 核心原理代码
```go
package main
import (
	"fmt"
	"net"
)

func main() {
	ipRecords, err := net.LookupIP("www.baidu.com")
	if err != nil {
		panic(err)
	}
	for _, value := range ipRecords {
		fmt.Println(value.String())
	}
}
```
> 利用net.LookupIP解析host 返回服务ip
```shell
14.119.104.254
14.119.104.189
```
输出ip 效果如果nslookup www.baidu.com一样：
```shell
Server:         10.96.0.10
Address:        10.96.0.10#53

Non-authoritative answer:
www.baidu.com   canonical name = www.a.shifen.com.
Name:   www.a.shifen.com
Address: 14.119.104.189
Name:   www.a.shifen.com
Address: 14.119.104.254
```

在部署grpc server的时候Service需要配置sessionAffinity(session亲和性),保证grpc server收得到消息后能正常返回
```yaml
apiVersion: v1
kind: Service
metadata:
  name: user-svc
  namespace: default
spec:
  clusterIP: None
  ports:
    - port: 8080
  selector:
    app:  user

  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 3600
```
## examples
具体见examples下的fronted和user项目

```shell
kubectl apply -f k8s.yaml
```
> 执行部署命令去部署

## 
```shell
root@hecs-410147:# kubectl  get svc
NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
frontend-svc     ClusterIP   10.108.199.130   <none>        80/TCP     40h
```
> 查看frontend-svc clusterIp
```shell

curl http://10.108.199.130/index
```
> 请求frontend-svc clusterIp

```shell
root@hecs-410147:~# kubectl logs user-5cdd5697f-vr5db
2023-04-09 22:22:21  file=build/main.go:33 level=info Starting [service] user
2023-04-09 22:22:21  file=v4@v4.9.0/service.go:96 level=info Transport [http] Listening on [::]:8080
2023-04-09 22:22:21  file=v4@v4.9.0/service.go:96 level=info Broker [http] Connected to 127.0.0.1:33039
2023-04-09 22:22:21  file=server/rpc_server.go:832 level=info Registry [memory] Registering node: user-defaaa6b-7314-4757-bb47-9a1ea6043d0d
2023-04-11 20:46:35  file=handler/user.go:16 level=info Received User.Call request: name:"gsmini@sina.cn"
2023-04-11 21:23:35  file=handler/user.go:16 level=info Received User.Call request: name:"gsmini@sina.cn"
2023-04-11 21:25:00  file=handler/user.go:16 level=info Received User.Call request: name:"gsmini@sina.cn"
2023-04-11 21:35:39  file=handler/user.go:16 level=info Received User.Call request: name:"gsmini@sina.cn"
2023-04-11 21:35:49  file=handler/user.go:16 level=info Received User.Call request: name:"gsmini@sina.cn"
```
> 查看user pod 的日志