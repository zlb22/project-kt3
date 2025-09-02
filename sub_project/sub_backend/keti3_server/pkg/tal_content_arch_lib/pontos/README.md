## Pontos 蓬托斯

框架命名[蓬托斯](https://baike.baidu.com/item/%E8%93%AC%E6%89%98%E6%96%AF/4452836) ，是古希腊神话中的海洋之神。

### 初始化服务
- 默认加载各组件和系统中间件
> 默认加载智能学习部门的TCM命名空间，SLS的Project等，非智能学习部门请勿使用该方式，请使用下面自定义组件+中间件的方式。
```golang 
package main

func main() {
  if err := pontos.InitApp(); err != nil {
    panic(err)
  }
  
  // register router
  router.RegisterRouter(pontos.Server.GinEngine())

  err := pontos.Server.ListenAndServe()
  if err != nil {
      log.Fatalf("Server stop err:%v\n", err)
  } else {
      log.Println("Server exit")
  }
}
```

- 自定义组件+中间件
```golang 
package main

func main() {
  cfgFile := "./config/local/cfg.toml"
  fmt.Println("config file path:", cfgFile)

  if err = pontos.NewApp(&pontos.ChainOptions{
      ConfigFile: cfgFile,
  }); err != nil {
      panic(err)
  }

  // register logger
  mainLogger, err := logger.NewPLogger("logs/xx.log")
  if err != nil {
      log.Fatal(err)
  }
  pontos.Logger = mainLogger

  pontos.Server.UseMiddleware(middleware.SetSkipLogKey())

  // register router
  router.RegisterRouter(pontos.Server.GinEngine())

  err = pontos.Server.ListenAndServe()
  if err != nil {
      log.Fatalf("Server stop err:%v\n", err)
  } else {
      log.Println("Server exit")
  }
}
```


### 验证服务
```
curl http://127.0.0.1:8096/pontos/ping
```

### Tag 版本说明
- beta   将一些新特性或者新功能集成之后的版本归为此类。更新频率较快
- bugfix 针对线上某些case的问题，需要紧急修复的情况下。更新频率：不固定
- stable 对于beta或者bugfix经过多个项目线上使用之后会升级到一个稳定版本。更新频率：较慢

### Tag 升级记录 >> [Change Log](CHANGELOG.md)

### 简易脚手架说明
#### 安装使用方式
* 具体操作请参考[pontos-cli](https://pkg/tal_content_arch_lib/pontos-cli/blob/master/README.md)

#### 命令列表
* 新建项目```pontos new myapp```
* 编译命令 ```pontos build```
* 监听项目文件，自动编译重启, ```pontos run```
* 生成数据库对应表的schema```pontos model table1 table2```
* 查看版本号 ```pontos version```

### 主要目录说明
- bootstrap  
包装GIN框架server，增加前置和后置方法，以及平滑重启等 


- config  
包装[Viper](https://github.com/spf13/viper)  


- constant  
框架所用到的常量定义 


- curl   
包装http请求
  

- kit  
工具箱，常用加解密，context等方法


- logger  
初始化logger以及hooks


- middleware   
框架所使用的基础中间件


- orm  
  包装[GORM](https://gorm.io/zh_CN/docs/index.html) 并实现其Logger的接口，用于格式化SQL语句打印


- redis   
  包装[go-redis](https://github.com/go-redis/redis/v8) 

### Gin 操作文档
- https://gin-gonic.com/zh-cn/docs/
- https://github.com/gin-gonic/gin

### gorm 操作文档
- https://gorm.io/zh_CN/docs/index.html
