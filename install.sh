user=$1
wssopass=$2
ldappass=$3
export wssoadminusername=$user
export wssoadminpassword=$wssopass
export ldapadminpassword=$ldappass
export jwtprivatekey=/opt/skinnywsso/tls/priv.pem
export jwtpublickey=/opt/skinnywsso/tls/pub.pem
export CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
cd /opt/
git clone https://github.com/evanjd711/skinnywsso.git
cd /opt/skinnywsso
apt install docker-compose -y
mkdir tls
cd /opt/skinnywsso/tls
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=US/ST=CA/L=Pomona/O=SkinnyWSSO/OU=SkinnyWSSO/CN=skinny.wsso"
openssl genrsa -out $jwtprivatekey 2048
openssl rsa -in $jwtprivatekey -pubout -out $jwtpublickey
docker-compose up 
sleep 3
docker exec skinnywsso_ldap_1 ldapadd -Y EXTERNAL -H ldapi:/// -f /opt/skinny_wsso/wsso.ldif
