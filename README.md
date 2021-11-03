### 在go.mod 里添加引用
```
因为etcd和grpc本身的兼容性问题，所以需要添加某些包的替换

replace (
    github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
    github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
    google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
```

### 安装包
`go get github.com/outreach-golang/etcd`
### 初始化etcd客户端
```
import "github.com/outreach-golang/etcd"

if err = etcd.GEtcd.InitEtcd(
   //环境变量值（DefaultEV,TestingEV,ReleaseEV,ProductionEV）
   etcd.EnvVar(etcd.DefaultEV),
   //etcd地址
   etcd.Points([]string{"127.0.0.1:2379"}),
   //ssl、ca等证书存放地址
   etcd.DirPath("./configs/default/k8s_keys/"),
); err != nil {
    log.Fatal(err.Error())
}
```
### 将服务注册到etcd
```
etcdRegister := etcd.NewServiceRegister(etcd.GEtcd.GetCli())

if err = etcd.ServiceRegisterHandler.InitServiceRegister(
    context.Background(),
    etcdRegister,
    "server.AppName",
    "server.port"); err != nil {
    log.Fatal(err.Error())
}
```
