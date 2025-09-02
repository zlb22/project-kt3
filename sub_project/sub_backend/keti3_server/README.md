### 蓬托斯脚手架
pontos-cli 依赖脚手架初始化项目

#### 简要说明
- config  `配置文件目录，环境依靠目录区分，默认配置文件名：app.toml`
  - development 
  - gray 
  - online 
  - test 
  - local `本地开发环境配置，git仓库忽略。`
  > 对于接入配置中心的项目， 配置文件可以直接删除。避免维护两处配置
- internal
  - application 
    - business `业务逻辑总目录，下面按照模块划分， 每个子目录下建议有interface.go定义该模块的接口`
    - wrapper(非必须) `各个business之间禁止相互调用。 如果在多个business层有共同代码时，为方便复用和防止循环引用，建议和business同级增加wrapper目录，其中增加对应的包，将business进行组合或者编排。`
    - constant `常量目录，尽可能的定义常量`
    - controller `只处理参数校验，具体实现都在business中处理`
    - errcode `定义本应用的错误码`
    - helper (非必须) `定义本应用的一些通用方法， 比如redis key的管理`
    - model `业务模型总目录，对应到各个表的CRUD，根据表的关联关系可聚合和拆分目录`
    - schema `定义本应用的数据表对应的结构体，也可拆分子目录`
    - proto `request、response等结构体`
    - service `一些下游依赖`
- pkg（非必须） 
- router `本应用所有的路由信息`
- script `执行各种构建、安装、分析等操作的脚本。`
- version `编译启动服务时，打印该服务的版本号等信息`

#### 项目编译
- 此项目为golang语言研发，本地需至少安装`go1.19版本`
- 进入项目根目录 `cd base_keti`
- 编译本地版本 `make build`
- 编译北师大版本[北师大为ubuntu是amd架构，若本机架构不同则需要交叉编译] `make build-beishida`
- 编译后在`bin`文件夹下会看到一个`keti3`的可执行文件

#### 项目启动
- 拷贝配置文件`app.toml`到`keti3`同目录，两者需要处于同一文件夹
- 执行命令`./keti3 -env=online`，即可启动课题3服务

#### 项目部署
- 将编译后的`keti3`可执行文件和`app.toml`配置文件放到服务器上的同一文件夹下启动即可


#### 重要文件介绍
- `./www/app.toml`，项目配置文件：可配置服务端口、mysql、redis链接
- `./router/`，项目路由文件：此项目提供的接口地址
- `./internal/application/schema/`，项目用到的所有数据表的结构体
- 接口`/oss/auth`为获取oss鉴权信息接口，以便前端进行文件上传，目前用的好未来自己的云存储，此接口务必更改为自己的存储