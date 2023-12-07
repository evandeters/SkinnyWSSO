cd /opt/
git clone https://github.com/evanjd711/skinnywsso.git
cd /opt/skinnywsso
apt install docker-compose -y
docker exec skinnywsso-ldap-1 ldapadd -Y EXTERNAL -H ldapi:/// -f /opt/skinny_wsso/wsso.ldif
mkdir tls
cd /opt/skinnywsso/tls
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=US/ST=CA/L=Pomona/O=SkinnyWSSO/OU=SkinnyWSSO/CN=skinny.wsso"