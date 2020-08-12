: cd当前目录
cd /d %~dp0

:（1）生成客户端私钥 （生成CA私钥）
openssl genrsa -out ca.key 2048

:（2）生成CA证书
openssl req -x509 -new -nodes -key ca.key -subj "/CN=" -days 5000 -out ca.crt

: 接下来，生成server端的私钥，生成数字证书请求，并用我们的ca私钥签发server的数字证书：

:（3）生成服务端私钥
openssl genrsa -out server.key 2048

:（4）生成证书请求文件
openssl req -new -key server.key -subj "/CN=localhost" -out server.csr

:（5）根据CA的私钥和上面的证书请求文件生成服务端证书
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000

：服务端使用：server.crt server.key
:  客户端使用ca.crt
