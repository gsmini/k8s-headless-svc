# 项目说明
## 关于项目
    当前项目是一个go-micro的register插件，当我们在用k8s部署go-micro的grpc server的时候，如果用k8s 内置的service
去部署grpc server的pod，是无法进行grpc的http2.0的，所以我采用headless service(无头service)方式去部署，把service当作
dns用，利用go的net包去做dns解析，然后获取service endpoints的podip记录返回，实现服务发现的功能。



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
```shell
14.119.104.254
14.119.104.189
```
> 输出ip 效果如果nslookup www.baidu.com一样

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
## example