# 财报机器人

# 私有云部署方式
## 一、部署表格识别引擎
详见表格识别引擎部署文档

## 二、部署财报NLP提取引擎

### 1. 镜像下载  

v0.0.7镜像下载链接 : https://dllf.intsig.net/download/2023/Solution/financial_statement_table_ner_v0.0.7.tar

--------
### 2. 部署方式
以v0.0.5版本部署为例：

* 加载NLP推理镜像`financial_statement_table_ner`
    ```
    docker load -i financial_statement_table_ner_v0.0.5_20221020.tar
    ```

* 部署NLP推理镜像
    ```
    docker run -it -d -p 30006:30006 -e OCR_ENGINE_URL=http://host:port/ai/service/v2/recognize/table?character=1 image_name -L 192.168.8.160:8885
    ```
    > * 其中环境变量`OCR_ENGINE_URL`为表格识别引擎地址，`host:port`修改为表格识别引擎的`host`与`port`,`character=1`必须要加
    > * 其中`-p 30006:30006` 第一个`30006`为宿主机端口号，可以自定义。第二个`30006`为容器内端口号，请勿修改。

## 三、部署数据库

### 部署 Mysql 数据库
> 部署Mysql数据库为可选项，如果行方内部有提供mysql服务，则可以直接使用行方的数据库，使用行方的数据库时，请直接在行方数据库中执行Mysql初始化脚本 `fr_mysql_init_v2.0.0.sql`

若要自己部署mysql数据库，则推荐使用docker部署方式，docker run命令如下：
```
docker run -it -d -p 3306:3306 -v /kingdee/caibao/20221027/mysql-config:/etc/mysql/conf.d -v /kingdee/caibao/20221027/mysql-data:/var/lib/mysql -v /kingdee/caibao/20221027/init_mysql.sql:/docker-entrypoint-initdb.d/init.sql -e MYSQL_ROOT_PASSWORD=intsig mysql:latest
```
> * 其中`-p 3306:3306`是mysql暴漏的端口号配置，可以根据需求修改第一个`3306`
> * 其中`-v /kingdee/caibao/20221027/mysql-config:/etc/mysql/conf.d`里的第一个路径为mysql的配置文件挂载目录
> * 其中`-v /kingdee/caibao/20221027/mysql-data:/var/lib/mysql`里的第一个路径为mysql数据持久化挂载目录
> * 其中`-v /kingdee/caibao/20221027/init_mysql.sql:/docker-entrypoint-initdb.d/init.sql`里的第一个路径为初始化项目时的sql脚本路径
> * 其中`-e MYSQL_ROOT_PASSWORD=intsig`中的intsig为mysql root账号的密码
### 部署 达梦 数据库
> 部署达梦数据库为可选项，如果行方内部有提供达梦数据库服务，则可以直接使用行方的数据库，使用行方的数据库时，请直接在行方数据库中执行达梦初始化脚本 `fr_dm8_init_v2.0.0.sql`

若要自己部署达梦数据库，则推荐使用docker部署方式，步骤如下：

1、下载达梦数据库docker image
可以在企业微信云盘里下载，也可以在 https://eco.dameng.com/download/ 下载  

2、加载达梦数据库docker image
```
docker load -i dm8_20220822_rev166351_x86_rh6_64_ctm.tar
```
3、启动达梦数据库
```
docker run -d -p 5236:5236 --restart=always --name dm8 --privileged=true -e PAGE_SIZE=16 -e LD_LIBRARY_PATH=/opt/dmdbms/bin -e UNICODE_FLAG=1 -e INSTANCE_NAME=dm8_01 -v /data/dm8_01:/opt/dmdbms/data dm8_single:v8.1.2.128_ent_x86_64_ctm_pack4
```
启动完成后，可通过日志检查启动情况，命令如下：
```
docker logs -f dm8_test
```
>**注意**   
>Docker 镜像中数据库默认用户名/密码为 SYSDBA/SYSDBA001。


## 四、部署Redis
> 部署Redis数据库为可选项，如果行方内部有提供redis服务，则可以直接使用行方的redis，要求redis使用时无需账号密码，财报将会使用redis db0 数据库

若要自己部署Redis服务，推荐使用docker方式部署，docker run命令如下：
```
docker run -it -d -p 6379:6379 redis:latest --slave-read-only no
```
> * 其中`-p 6379:6379`为对外暴漏的端口号，第一个6379可以根据需要自行修改

## 五、部署财报机器人后端
### 1. 版本对应列表
无
### 2. 部署方式
以上述版本列表中的v2.0.0版本部署为例：

* 加载image镜像
    ```
    docker load -i financial_statement_v2.0.0_20221020.tar
    ```
* 部署image镜像
    ```
    docker run -it --name fr-service -d -p 8070:8080  --restart=always  \
        -v your_path:/usr/local/fr-files/     \
        -e APP_AUTH_AUTH_TYPE=default     \
        -e APP_TASK_RESULT_ASYNC_URL=http://xxxxx.com/result   \
        -e APP_OCR_RECOGNIZE_TABLE_API=http://host:port/ai/service/v2/recognize/table   \
        -e APP_MYSQL_DATABASE=financial_statement   \
        -e APP_MYSQL_HOST=192.168.60.73 \
        -e APP_MYSQL_PORT=3306  \
        -e APP_MYSQL_USERNAME=textin    \
        -e APP_MYSQL_PASSWORD=is@SHWJC  \
        -e APP_REDIS_HOST=192.168.60.55   \
        -e APP_REDIS_PORT=6379  \
        financial_statement:v2.0.0
    ```
    > * 其中`APP_OCR_RECOGNIZE_TABLE_API`的环境变量的值里，`host:port`修改为NLP提取引擎的`host`与`port`
    > * `-v your_path:/usr/local/fr-files/`为财报数据存储挂载命令，命令中的`your_path`需要替换为需要挂载的目录地址

    支持的环境变量
    |  ENV   | 说明  |
    |  ----  | ----  |
    |  APP_SERVER_DB  | 数据库类型;default:mysql enum：mysql、dm                    |
    |  APP_SERVER_FILE_STORAGE  | 文件存储方式; default:local; enum：local,nfs     |
    |  APP_SERVER_RECOGNIZEMAXRETRY  | 识别任务最大尝试次数; default:10     |
    |  APP_AUTH_AUTH_TYPE  | 鉴权方式; default:default; enum：default,url-token,cas     |
    |  APP_AUTH_CHECK_TOKEN_API  | url-token鉴权时，鉴权url;   |
    |  APP_AUTH_CAS_HANDLER  | 如果是cas单点登录，这里指定使用哪个handler来处理cas认证流程 |
    |  APP_AUTH_CAS_QUERY_KEY  | 当sso认证方式为CAS协议时，这里输入前端传递TGC时的url参数key name |
    |  APP_AUTH_CAS_SERVER  | 当sso认证方式为CAS协议时，这里输入CAS的TGT校验服务地址 |
    |  APP_AUTH_CAS_UID_PATH  | 当sso认证方式为CAS协议时，这里输入ticket验证通过后用户信息json路径 |
    |  APP_AUTH_CAS_CLIENT_SERVER  | 当sso认证方式为CAS协议时，这里输入客户端在CAS服务里注册的service地址，一般是前端登录地址 |
    |  APP_TASK_RESULT_ASYNC_URL  | 识别结果回调的url，仅支持POST请求; default:"";      |
    |  APP_TASK_RESULT_ASYNC_SAVE_WITH_SYNC  | web端编辑完财报后是否回调，默认关闭，开启可设置为true;      |
    |  APP_TASK_RESULT_ASYNC_RE_IDENTIFY_WITH_SYNC  | web端选择任务重新识别后，是否需要发起识别结果回调，默认关闭，开启可设置为true;      |
    |  APP_FILE_CLEANER_ENABLE  | 文件清理功能,default:true enum: true,false;      |
    |  APP_LOCAL_STORAGE_SAVE_DIR  | 本地存储文件模式下，容器内的存储路径,default:/usr/local/fr-files/     |
    |  APP_OCR_RECOGNIZE_TABLE_API  | 表格识别引擎url     |
    |  APP_MYSQL_DATABASE  | mysql database数据库名称     |
    |  APP_MYSQL_HOST  | mysql host     |
    |  APP_MYSQL_PORT  | mysql port     |
    |  APP_MYSQL_USERNAME  | mysql username     |
    |  APP_MYSQL_PASSWORD  | mysql password     |
    |  APP_DAMENG_HOST  | 达梦数据库模式下的db host     |
    |  APP_DAMENG_PORT  | 达梦数据库模式下的db port     |
    |  APP_DAMENG_USERNAME  | 达梦数据库模式下的db username，例如可填写SYSDBA     |
    |  APP_DAMENG_PASSWORD  | 达梦数据库模式下的db password，例如填写SYSBA001     |
    |  APP_DAMENG_DATABASE  | 达梦数据库模式下的db database，例如可填写SYSDBA     |
    |  APP_REDIS_TYPE  | redis type   redis模式，default：default ；enum：default（默认单机） 、 sentinel（哨兵）、cluster（集群）  |
    |  APP_REDIS_HOST  | redis host     |
    |  APP_REDIS_PORT  | redis port     |
    |  APP_REDIS_USERNAME  | redis username     |
    |  APP_REDIS_PASSWORD  | redis password     |
    |  APP_REDIS_DATABASE  | redis db index 默认0     |
    |  APP_REDIS_MASTERNAME  | redis 哨兵模式的masterName     |
    |  APP_REDIS_SENTINELADDRS  | redis 哨兵模式的sentinelAddrs，多个使用`,`分割,例如：127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381    |
    |  APP_REDIS_SENTINELPASSWORD  | redis 哨兵模式的sentinelPassword     |
    |  APP_REDIS_CLUSTERADDRS  | redis Cluster模式的ip地址，多个使用`,`分割,例如：127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381    |

### 关于Redis配置说明
Redis支持`单机`、`哨兵`、`集群`三种部署方式，以下是三种部署方式要关注的环境变量（`部署时需按照真实情况配置`）：
``` bash
# 单机模式
-e APP_REDIS_TYPE=default
-e APP_REDIS_HOST=127.0.0.1
-e APP_REDIS_PORT=6739
-e APP_REDIS_USERNAME=userName  #若有
-e APP_REDIS_PASSWORD=password  #若有

# 哨兵模式
-e APP_REDIS_TYPE=sentinel
-e APP_REDIS_MASTERNAME=masterName
-e APP_REDIS_SENTINELADDRS=127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381
-e APP_REDIS_SENTINELPASSWORD=sentinelPassword  # 若有
-e APP_REDIS_USERNAME=userName  # 若有
-e APP_REDIS_PASSWORD=password  # 若有

# 集群模式
-e APP_REDIS_TYPE=cluster
-e APP_REDIS_CLUSTERADDRS=127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381
-e APP_REDIS_USERNAME=userName  #若有
-e APP_REDIS_PASSWORD=password  #若有
```

## 六、部署财报机器人Web端
```
docker run -it -d --name=fr-web --restart=always -p 8081:80 -e VUE_APP_BASE_URL=http://host:port/  image_name
```
支持的环境变量
|  ENV   | 说明  |
|  ----  | ----  |
|  VUE_APP_BASE_URL  | 财报后端服务地址                     |
|  VUE_APP_TITLE  | Web端HTML 显示的 Title，默认是“财报机器人”                     |
|  VUE_APP_FLAG_REMOVE_LOGO  | Web端是否展示合合信息的logo，'enable'为不保留, 其他为保留,默认开启                     |
|  VUE_APP_FLAG_SSO  | 是否开启sso登录功能，'enable'为开启, 其他为关闭,默认关闭                    |
|  VUE_APP_FLAG_ADMIN  | 是否开启sso保留admin登录，'enable'为开启, 其他为关闭                     |
|  VUE_APP_SSO_KEY  | 开启sso登录功能后地址栏的参数名称，默认为token                     |
|  VUE_BASE_PATH  | 二级域名配置字段，默认为for_financial_statements                     |
|  VUE_API_CONTENT  | 接口拼接配置字段，默认为/financial_statements_api/，记得前后加上/                     |
|  API_SSO_LOGIN_URL  | SSO登录地址                     |
|  API_SSO_LOGOUT_URL  | SSO模式下，系统退出登录时，跳转到的地址，一般填写SSO登录地址                     |
|  VUE_APP_NORM_CHANGE  | 是否启用编辑页的准则变更，'disable'为关闭，默认开启                     |
|  VUE_APP_ADJUST_GROUP  | 是否启用调整分组功能，'disable'为关闭，默认开启                     |
|  VUE_APP_CANCELLATION  | 是否启用财报作废功能，'enable'为开启，默认关闭                     |
|  VUE_APP_FR_EVENT  | 编辑页面的提交和作废按钮的postMessage事件通知开关，'enable'为开启，默认关闭；  提交的标识是： FR_EVENT_SUBMIT 作废的标识是： FR_EVENT_OBSOLETE  ，具体使用方式详见下方对该字段的说明                 |
|  VUE_APP_EXPORT_EXCEL  | 是否启用导出EXCEL功能，'disable'为关闭，默认开启                     |
|  VUE_APP_MODIFY_RECORD  | 是否启用修改记录功能，'disable'为关闭，默认开启                     |
|  VUE_APP_TOP_MENU_BAR  | 是否显示顶部菜单栏，'disable'为关闭，默认开启                     |
|  VUE_APP_REDIRECT_TO_LOGIN  | 页面鉴权失败后是否跳转到登录页，'disable'为关闭，默认开启,当设置为'disable'时，访问财报界面需要在url里带有X-TOKEN,其值为调用财报登录接口获取的token值                    |
|  VUE_APP_CUSTOM_NORM_LIST  | 是否支持自定义准则列表，'enable'为开启，默认关闭                     |
|  VUE_APP_FORMULA_AND_SUBJECT_MATCH_CHECK  | 财报编辑后点击提交，是否进行公式配平和科目匹配校验，'enable'为开启，默认关闭                 |
|  VUE_APP_MAX_UPLOAD_IMAGE_COUNT_LIMIT  | 上传文件数可配置选项 '默认一次性最多10个图片文件'                     |
|  VUE_APP_UPLOAD_PDF_IMAGE_COUNT  | 选择上传PDF图片数 '默认全部图片'                     |
|  VUE_APP_FILTER_ABLE  | 财报编辑页下拉框是否支持输入，'disable'为关闭，默认开启                     |
|  VUE_APP_TABLE_HEADER_SELECT  | 财报编辑页表头是否支持下拉修改，'enable'为开启，默认关闭                     |
|  VUE_APP_BALANCE_SHEET_HEADER_FIELD  | 资产负债表头配置字段，格式如：期末数,期初数。注意中间使用英文输入法的逗号隔开。                     |
|  VUE_APP_CASH_FLOW_HEADER_FIELD  | 现金流量表头配置字段，格式如：期末数,期初数。注意中间使用英文输入法的逗号隔开。                     |
|  VUE_APP_INCOME_HEADER_FIELD  | 利润表头配置字段，格式如：期末数,期初数。注意中间使用英文输入法的逗号隔开。                     |
|  VUE_APP_ONLY_ADMIN_MENU  | 是否开启只有admin才显示“控制管理台” 和 “账号管理”菜单，'enable'为开启，默认关闭  |

### 特殊字段说明
#### VUE_APP_FR_EVENT
该字段适用于业务方希望捕获编辑界面“提交”按钮事件的场景，若有该场景，请开启该功能，然后业务方需要将财报的编辑界面使用iframe进行嵌套，并在自己业务层的js里增加事件监听，示例代码如下：
```
<iframe id="frame1" src="http://xxx.xxx.xxx.xxx:pppp/financial_statements/show?task_id=652" frameborder="0" scrolling="no" width="100%" height="960px"></iframe>
<script >
    window.addEventListener('message', function(event) {
        console.log( "子页面传递了消息到父页面："+ event.data);
        // 处理后续的业务逻辑 START
        …………
        // 业务逻辑 END
    });
</script>
```
提交的事件标识是： FR_EVENT_SUBMIT
作废的事件标识是： FR_EVENT_OBSOLETE

### 暴漏的API与Web

### API
* http://domain/api/v2/task/create  创建财报任务

    >API文档定义详见API文档OpenApi
### Web
* http://domain/for_financial_statements/user/login     登录界面（默认账号密码：admin/intsig@888）
* http://domain/for_financial_statements/show?task_id=249   财报任务编辑界面

> 其中`domain`为前端`host:port`
