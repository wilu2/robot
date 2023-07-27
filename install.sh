
FR_WEB_CONTAINER_NAME='fr-web'
FR_SERVICE_CONTAINER_NAME='fr-service'


FUN_CHECK_ARG() {
    if [ ! -n '$1' ];then
        echo '参数错误！'
        exit
    fi
}

FUN_INSTALL_WEB() {
    echo '开始安装财报前端...'
    read -p '输入财报前端包名称：' fileName
    read -e -p '输入财报前端所使用的端口号：' -i "8080" port
    read -e -p '输入后端服务API Host与端口:' -i "10.4.18.26:8081" domain
    echo '开始load $fileName ...'
    eval 'docker load -i ${fileName}'
    read -e -p '输入上面一行Loaded image提示的版本信息:' -i "fr-web:v4.0.0" version
    eval 'docker rm -f ${FR_WEB_CONTAINER_NAME} >> /dev/null 2>&1'
    eval 'docker run -it --name ${FR_WEB_CONTAINER_NAME} -d -p ${port}:80 -e VUE_APP_BASE_URL=http://${domain}/ ${version}'
}

FUN_INSTALL_SERVER() {
    echo '开始安装财报后端...'
    read -p '输入后端包名称：' fileName
    read -e -p '输入所使用的端口号：'            -i "8081" port
    read -e -p '输入NLP引擎API Host与端口:'     -i "10.4.18.26:8082" nlpDomain
    read -e -p '输入Mysql 数据库名称:'          -i "financial_statement" mysqlDbName
    read -e -p '输入Mysql Host：'               -i "10.4.16.16" mysqlHost
    read -e -p '输入Mysql Port：'               -i "3306" mysqlPort
    read -e -p '输入Mysql UserName：'           -i "root" mysqlUserName
    read -e -p '输入Mysql PASSWORD：'           -i "intsig" mysqlPassword
    read -e -p '输入Redis Host：'               -i "10.4.16.17" redisHost
    read -e -p '输入Redis Port：'               -i "6379" redisPort
    echo '开始load $fileName ...'
    eval 'docker load -i ${fileName}'
    read -e -p '输入上面一行Loaded image提示的版本信息(例如：financial_statement:v2.0.0):' -i "financial_statement:v2.0.0" version
    eval 'docker rm -f ${FR_SERVICE_CONTAINER_NAME} >> /dev/null 2>&1'
    eval 'docker run -it --name ${FR_SERVICE_CONTAINER_NAME} -p ${port}:8080 -d -v /usr/local/fr-files/:/usr/local/fr-files/ -e APP_AUTH_AUTH_TYPE=url-token -e APP_OCR_RECOGNIZE_TABLE_API=http://${nlpDomain}/ai/service/v2/recognize/table -e APP_MYSQL_HOST=${mysqlHost} -e APP_MYSQL_PORT=${mysqlPort} -e APP_MYSQL_USERNAME=${mysqlUserName} -e APP_MYSQL_PASSWORD=${mysqlPassword} -e APP_MYSQL_DATABASE=${mysqlDbName} -e APP_REDIS_HOST=${redisHost} -e APP_REDIS_PORT=${redisPort} -e GIN_MODE=release ${version}'
}

FUN_UNINSTALL_WEB() {
    echo '开始卸载财报前端...'
    eval 'docker rm -f ${FR_WEB_CONTAINER_NAME}'
}

FUN_UNINSTALL_SERVER() {
    echo '开始卸载财报后端...'
    eval 'docker rm -f ${FR_SERVICE_CONTAINER_NAME}'
}

echo '菜单'
echo '1.安装财报后端(请提前准备好Mysql Redis OCR引擎信息)           2.安装财报前端(需要提前安装好财报后端)'
echo '3.卸载财报后端        4.卸载财报前端'
read -p '请选择：' cmd
FUN_CHECK_ARG cmd
if [ $cmd == '1' ]
then
    FUN_INSTALL_SERVER
elif [ $cmd == '2' ]
then    
    FUN_INSTALL_WEB
elif [ $cmd == '3' ]
then
    FUN_UNINSTALL_SERVER
elif [ $cmd == '4' ]
then    
    FUN_UNINSTALL_WEB
fi