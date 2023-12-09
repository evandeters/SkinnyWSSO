#!/bin/bash
export wssoadminusername=$1
export wssoadminpassword=$2
export ldapadminpassword=$3
export jwtprivatekey=/opt/skinnywsso/tls/priv.pem
export jwtpublickey=/opt/skinnywsso/tls/pub.pem
export CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
orig_dir=$(pwd)
cp -R $orig_dir /opt/skinnywsso
cd /opt/skinnywsso
apt install docker-compose -y
mkdir tls
cd /opt/skinnywsso/tls
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=US/ST=CA/L=Pomona/O=SkinnyWSSO/OU=SkinnyWSSO/CN=skinny.wsso"
openssl genrsa -out $jwtprivatekey 2048
openssl rsa -in $jwtprivatekey -pubout -out $jwtpublickey
docker-compose up -d
sleep 3
docker exec skinnywsso_ldap_1 ldapadd -Y EXTERNAL -H ldapi:/// -f /opt/skinny_wsso/wsso.ldif
