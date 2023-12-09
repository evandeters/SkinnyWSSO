/* Requires the Docker Pipeline plugin */
pipeline {
    agent none
    environment {
        WSSO_ADMIN = credentials('wsso_admin_creds')
    }
    stages {
        stage('Initialize') {
            agent any
            steps {
                sh '''
                    /usr/bin/docker volume create ldap_data
                    /usr/bin/docker volume create ldap_slapd
                    /usr/bin/docker network create skinnywsso
                '''
            }
        }
        stage('LDAP') {
            agent {
                environment {
                    LDAP_ORGANISATION: "Skinny WSSO"
                    LDAP_DOMAIN: "skinny.wsso"
                    LDAP_ADMIN_PASSWORD: '$WSSO_ADMIN_PWD'
                }
                docker { 
                    image 'osixia/openldap:latest'
                    args '--network=skinnywsso -p 389:389 -p 636:636 -v /data/ldap:/var/lib/ldap -v /data/slapd.d:/etc/ldap/slapd.d -e LDAP_ORGANISATION="$LDAP_ORGANISATION" -e LDAP_DOMAIN="$LDAP_DOMAIN" -e LDAP_ADMIN_PASSWORD="$LDAP_ADMIN_PASSWORD"'
                }
            }
            steps {
                sh('chmod +x ./install.sh; ./install.sh $WSSO_ADMIN_USR $WSSO_ADMIN_PSW $WSSO_ADMIN_PSW')
            }
        }

        stage('WSSO') {
            agent {
                environment {
                    WSSO_ADMIN_PASSWORD: '$WSSO_ADMIN_PWD'
                    WSSO_ADMIN_USERNAME: '$WSSO_ADMIN_USR'
                    USE_HTTPS: 'true'
                    CERT_PATH: "/opt/skinnywsso/tls/cert.pem"
                    KEY_PATH: "/opt/skinnywsso/tls/key.pem"
                    JWT_PRIVATE_KEY: "/opt.skinnywsso/tls/priv.pem"
                    JWT_PUBLIC_KEY: "/opt.skinnywsso/tls/pub.pem"
                }
                docker { 
                    image 'golang:latest'
                    args '--network=skinnywsso -p 443:443 -p 80:80 -e WSSO_ADMIN_PASSWORD="$WSSO_ADMIN_PASSWORD" -e WSSO_ADMIN_USERNAME="$WSSO_ADMIN_USERNAME" -e USE_HTTPS="$USE_HTTPS" -e CERT_PATH="$CERT_PATH" -e KEY_PATH="$KEY_PATH" -e JWT_PRIVATE_KEY="$JWT_PRIVATE_KEY" -e JWT_PUBLIC_KEY="$JWT_PUBLIC_KEY"'
                }
            }
            steps {
                sh '''
                    orig_dir=$(pwd)
                    cp -R $orig_dir /opt/skinnywsso
                    cd /opt/skinnywsso
                    mkdir tls
                    cd /opt/skinnywsso/tls
                    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=US/ST=CA/L=Pomona/O=SkinnyWSSO/OU=SkinnyWSSO/CN=skinny.wsso"
                    openssl genrsa -out $jwtprivatekey 2048
                    openssl rsa -in $jwtprivatekey -pubout -out $jwtpublickey
                    cd /opt/skinnywsso
                    go run .
                '''
            }
        }
    }
}
