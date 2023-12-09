/* Requires the Docker Pipeline plugin */
pipeline {
    agent any

    options {
    // Only keep the 10 most recent builds
    buildDiscarder(logRotator(numToKeepStr:'10'))
    }

    environment {
        WSSO_ADMIN = credentials('wsso_admin_creds')
    }

    stages {
        stage('LDAP') {
            steps {
                    sh '''
                        rm -rf /var/lib/ldap/* 
                        cp -R /root/ldap_backup/* /var/lib/ldap/
                        chown -R openldap:openldap /var/lib/ldap/
                        systemctl restart slapd
                        ldapadd -x -w $WSSO_ADMIN_PSW -H ldapi:/// -D cn=admin,dc=skinny,dc=wsso -f ./wsso.ldif
                    '''
                }
        }

        stage('WSSO') {
            agent {
                docker { 
                    image 'golang:latest'
                    args '--add-host host.docker.internal:host-gateway --network=skinnywsso -p 443:443 -p 80:80 -e WSSO_ADMIN_PASSWORD="$WSSO_ADMIN_PSW" -e LDAP_ADMIN_PASSWORD="$WSSO_ADMIN_PSW" -e WSSO_ADMIN_USERNAME="$WSSO_ADMIN_USR" -e USE_HTTPS="true" -e CERT_PATH="/opt/skinnywsso/tls/cert.pem" -e KEY_PATH="/opt/skinnywsso/tls/key.pem" -e JWT_PRIVATE_KEY="/opt/skinnywsso/tls/priv.pem" -e JWT_PUBLIC_KEY="/opt/skinnywsso/tls/pub.pem"'
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
                    openssl genrsa -out $JWT_PRIVATE_KEY 2048
                    openssl rsa -in $JWT_PRIVATE_KEY -pubout -out $JWT_PUBLIC_KEY
                    cd /opt/skinnywsso
                    go run .
                '''
            }
        }
    }
}
