/* Requires the Docker Pipeline plugin */
pipeline {
    agent none
    environment {
        WSSO_ADMIN = credentials('wsso_admin_creds')
    }
    stages {
        stage('LDAP') {
            agent {
                docker { 
                    image 'osixia/openldap:latest'
                    args '--network=skinnywsso -p 389:389 -p 636:636 -v ldap_data:/var/lib/ldap -v ldap_slapd:/etc/ldap/slapd.d -e LDAP_ORGANISATION="Skinny WSSO" -e LDAP_DOMAIN="skinny.wsso" -e LDAP_ADMIN_PASSWORD="$WSSO_ADMIN_PWD"'
                }
                steps {
                    sh('slapcat')
                }
            }
        }

        stage('WSSO') {
            agent {
                docker { 
                    image 'golang:latest'
                    args '--network=skinnywsso -p 443:443 -p 80:80 -e WSSO_ADMIN_PASSWORD="$WSSO_ADMIN_PWD" -e WSSO_ADMIN_USERNAME="$WSSO_ADMIN_USR" -e USE_HTTPS="true" -e CERT_PATH="/opt/skinnywsso/tls/cert.pem" -e KEY_PATH="/opt/skinnywsso/tls/key.pem" -e JWT_PRIVATE_KEY="/opt.skinnywsso/tls/priv.pem" -e JWT_PUBLIC_KEY="/opt.skinnywsso/tls/pub.pem"'
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
