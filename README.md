# mybili
基于gin框架的视频网站

## 介绍

gin+vue 全栈制作一个视频网站。

网站整体功能模仿B站设计和制作，仅作为学习交流使用。

网站链接，欢迎访问交流！
[mybili](http://www.mybili.fun)

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
- 其它
  - 腾讯云对象存储(COS)
  - nginx

## 项目预览

- 主页
  ![](https://github.com/Mjiarong/mybili/blob/main/preview/main.jpg)

- 视频播放页面
  ![](https://github.com/Mjiarong/mybili/blob/main/preview/video.jpg)

- 视频评论页面

  ![](https://github.com/Mjiarong/mybili/blob/main/preview/comment.jpg)

- 投稿页面

 ![](https://github.com/Mjiarong/mybili/blob/main/preview/upload.jpg)

 - 个人视频管理页面

  ![](https://github.com/Mjiarong/mybili/blob/main/preview/edit.jpg)

##  Docker部署
### 一、创建后端项目镜像
```shell
#首先拉取项目仓库代码到你的机器上
$ git clone https://github.com/Mjiarong/mybili)https://github.com/Mjiarong/mybili

#进入项目主目录，执行docker build指令，生成docker镜像
$ docker build -t mybili:v1 .

# 将镜像推送到镜像仓库，这里以阿里云镜像仓库为例
$ docker tag mybili:v1 registry.cn-guangzhou.aliyuncs.com/yourname/mybili:v1
$ docker push registry.cn-guangzhou.aliyuncs.com/yourname/mybili:v1
```

### 二、创建前端项目镜像
```shell
#首先拉取项目仓库代码到你的机器上
$ git clone https://github.com/Mjiarong/mybiliweb)https://github.com/Mjiarong/mybiliweb

#进入项目主目录下的docker目录，把vue打包生成的dist文件夹copy到当前目录下
$ cd ./docker.
$ cp -r ../dist ./

#执行docker build指令，生成docker镜像
$ docker build -t mybili-vue:v1 . 

# 将镜像推送到镜像仓库，这里以阿里云镜像仓库为例
$ docker tag mybili-vue:v1 registry.cn-guangzhou.aliyuncs.com/yourname/mybili-vue:v1
$ docker push registry.cn-guangzhou.aliyuncs.com/yourname/mybili-vue:v1
```

### 三、使用docker-compose启动工程
```shell
#首先在服务器/usr目录下创建一个文件夹mybili-project
$ cd /usr
$ mkdir -p mybili-project

#进入mybili-project文件夹，创建三个文件夹，分别名为compose  nginx  redis
#进入mybili-project/nginx文件夹,继续创建conf文件夹
#将你修改好的redis.conf文件放到mybili-project/redis文件夹下,作挂载使用
#将前端工程目录mybiliweb/docker下的nginx.conf文件和default.conf文件复制到mybili-project/nginx文件夹下,作挂载使用
#将前端工程目录mybiliweb/docker下的docker-compose.yml复制到mybili-project/compose文件夹下
#进入mybili-project/compose目录，执行docker-compose.yml文件启动工程
$ cd mybili-project/compose
$ docker compose up -d
```
- 文件目录树
![](https://github.com/Mjiarong/mybili/blob/main/preview/dir-tree.png)


### 四、在服务器上安装nginx
```shell
#配置Centos 7 Nginx Yum源仓库，安装Nginx
$ rpm -Uvh http://nginx.org/packages/centos/7/noarch/RPMS/nginx-release-centos-7-0.el7.ngx.noarch.rpm
$ yum install nginx -y

#注释掉/etc/nginx/nginx.conf 中的server默认配置
http {
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 4096;

    include             /etc/nginx/mime.types;
    default_type        application/octet-stream;

    # Load modular configuration files from the /etc/nginx/conf.d directory.
    # See http://nginx.org/en/docs/ngx_core_module.html#include
    # for more information.
    include /etc/nginx/conf.d/*.conf;

#    server {
#        listen       80;
#        listen       [::]:80;
#        server_name  _;
#        root         /usr/share/nginx/html;
#
#        # Load configuration files for the default server block.
#        include /etc/nginx/default.d/*.conf;
#
#       error_page 404 /404.html;
#       location = /404.html {
#        }
#
#        error_page 500 502 503 504 /50x.html;
#        location = /50x.html {
#        }
#    }

# Settings for a TLS enabled server.
#
#    server {
#        listen       443 ssl http2;
#        listen       [::]:443 ssl http2;
#        server_name  _;
#        root         /usr/share/nginx/html;
#
#        ssl_certificate "/etc/pki/nginx/server.crt";
#        ssl_certificate_key "/etc/pki/nginx/private/server.key";
#        ssl_session_cache shared:SSL:1m;
#        ssl_session_timeout  10m;
#        ssl_ciphers HIGH:!aNULL:!MD5;
#        ssl_prefer_server_ciphers on;
#
#        # Load configuration files for the default server block.
#        include /etc/nginx/default.d/*.conf;
#
#        error_page 404 /404.html;
#            location = /40x.html {
#        }
#
#        error_page 500 502 503 504 /50x.html;
#            location = /50x.html {
#        }
#    }


#创建自己的配置文件
$ cd /etc/nginx/conf.d/
$ vi mybili.conf

#输入以下内容，保存
server {
        #服务监听端口
        listen       80;
        #域名
        server_name  mybili.fun;
        #请求的地址
        location / {
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_pass http://127.0.0.1:3001;
        }

        location /api {
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_pass http://127.0.0.1:3000;
        }
}

#启动Nginx,并设置开启自启动
$ systemctl start nginx
$ systemctl enable nginx
```
## 项目运行架构
### 总体架构
![](https://github.com/Mjiarong/mybili/blob/main/preview/jiagou.jpg)

### COS存储数据流程
- 下载：前端每次请求前都需要通过后端获取签名，后端使用固定密钥或临时密钥计算签名返回至前端，前端收到有效预签名URL下载对象。
- 上传：前端每次请求前都需要通过后端获取签名，后端通过获取临时密钥给到前端，前端计算签名。
![](https://github.com/Mjiarong/mybili/blob/main/preview/cos.jpg)
