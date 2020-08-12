::public
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 7200 -key ca.key -out ca.crt

::server
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 7200 -in server.csr -out server.crt

::client
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 7200 -in client.csr -out client.crt

goto start

ca.key           // ca私钥【保密】
ca.crt         // ca公钥
ca.srl            // 签名期间生成【保密】
client.csr      // 证书签名请求文件【可删除】
client.key     // 私钥
client.crt   // 公钥
server.csr     // 【可删除】
server.key
server.crt

:start