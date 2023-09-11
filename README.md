# mybili
基于gin框架的视频网站

## 介绍

gin+vue 全栈制作一个视频网站。

网站整体功能模仿B站设计和制作，仅作为学习交流使用。

## 目录结构
### 后端源码
```shell
├─  .gitignore
│  go.mod // 项目依赖
│  go.sum
│  LICENSE
│  main.go //主程序
│  README.md
│  Dockerfile
|  .env // 项目配置文件,可以参考env.example  
│          
├─api //接口服务层
├─model // 数据模型层
├─service //业务层
├─serializer //序列化器
│   
├─conf 
│ └─conf.go  // 初始化函数
├─cache // 缓存数据库redis
├─log  // 项目日志
├─middleware  // 中间件
├─server
│  └─router.go // 路由入口     
└─utils // 项目公用工具库      
```

### 前端源码
仓库地址:https://github.com/Mjiarong/mybiliweb
```shell      
├─docker //容器部署
│   ├─Dockerfile
│   ├─default.conf //nginx配置文件
│   ├─nginx.conf //nginx配置文件
│   └─docker-compose.yml
│   
├─public 
│ └─index.html  //打包目标文件
├─src // vue工程源文件
├─jbabel.config.js
├─jsconfig.json 
├─package.json  
│─package-lock.json  
└─vue-config.json
```
## 实现功能

1.  主页视频展示
2.  视频播放
3.  视频上传功能
4.  用户登录与用户注册
5.  用户后台管理视频


## 技术栈

- 后端
  - Golang
  - Gin web framework
  - gorm
  - jwt-go
  - 日志框架:ogrus+lfshook+file-rotatelogs
  - gin-contrib/cors
- 前端
  - MVVM框架:vue 2.x
  - vue router
  - ui框架:ElementUI
  - axios
- 数据库
  - MySQL
  - redis

## 项目预览

- 主页
  ![](https://github.com/Mjiarong/mybili/blob/main/preview/main.jpg)

- 视频播放页面
  ![](https://github.com/Mjiarong/mybili/blob/main/preview/video.jpg)

- 视频评论页面

  ![](https://github.com/Mjiarong/ginblog/blob/main/upload/comment.png)

- 投稿页面

  ![](https://github.com/Mjiarong/ginblog/blob/main/upload/upload.png)

 - 个人视频管理页面

  ![](https://github.com/Mjiarong/ginblog/blob/main/upload/edit.png)
