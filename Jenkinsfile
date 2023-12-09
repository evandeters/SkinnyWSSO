/* Requires the Docker Pipeline plugin */
pipeline {
    agent any

    options {
    // Only keep the 10 most recent builds
    buildDiscarder(logRotator(numToKeepStr:'10'))
    }

    environment {
        WSSO_ADMIN = credentials('wsso_admin_creds')
        LDAP_ADMIN_PASSWORD = credentials('ldap_admin_password')
    }

    tools { go '1.21.5'}

    stages {
        stage('Compile') {
            steps {
                sh 'go build'
            }
        }

        stage('Prepare Tests') {
            steps {
                sh '''
                    export jwtprivatekey=/var/jenkins/workspace/SkinnyWSSO/priv.pem
                    export jwtpublickey=/var/jenkins/workspace/SkinnyWSSO/pub.pem
                    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=US/ST=CA/L=Pomona/O=SkinnyWSSO/OU=SkinnyWSSO/CN=skinny.wsso"
                    openssl genrsa -out $jwtprivatekey 2048
                    openssl rsa -in $jwtprivatekey -pubout -out $jwtpublickey
                    rm -rf /var/lib/ldap/*
                    cp -R /root/ldap_backup/* /var/lib/ldap/
                    chown -R openldap:openldap /var/lib/ldap/
                    systemctl restart slapd
                '''
                sh 'ldapadd -x -H ldapi:/// -f ~/wsso.ldif -D cn=admin,dc=skinny,dc=wsso -w $LDAP_ADMIN_PASSWORD'
            }
        }

        stage('Unit Tests') {
            environment {
                JWT_PRIVATE_KEY = '/var/jenkins/workspace/SkinnyWSSO/priv.pem'
                JWT_PUBLIC_KEY = '/var/jenkins/workspace/SkinnyWSSO/pub.pem'
            }
            steps {
                sh '''
                    go test -v
                '''
            }
        }

        stage('Deploy to Dev') {
            when {
                branch 'development'
            }
            steps {
                sshagent(['ssh_key']) {
                    sh '''
                        ssh -o StrictHostKeyChecking=no skinnywsso-dev 'systemctl stop skinnywsso.service'
                        scp -o StrictHostKeyChecking=no SkinnyWSSO skinnywsso-dev:/opt/skinnywsso/
                        scp -o StrictHostKeyChecking=no -r templates/ skinnywsso-dev:/opt/skinnywsso/
                        scp -o StrictHostKeyChecking=no wsso.ldif skinnywsso-dev:/opt/skinnywsso/
                        ssh -o StrictHostKeyChecking=no skinnywsso-dev 'systemctl start skinnywsso.service'
                        ssh -o StrictHostKeyChecking=no skinnywsso-dev 'rm -rf /var/lib/ldap/*; cp -R /root/ldap_backup/* /var/lib/ldap/; chown -R openldap:openldap /var/lib/ldap/; systemctl restart slapd'
                    '''
                    sh 'ssh -o StrictHostKeyChecking=no skinnywsso-dev "ldapadd -x -H ldapi:/// -f /opt/skinnywsso/wsso.ldif -D cn=admin,dc=skinny,dc=wsso -w $LDAP_ADMIN_PASSWORD"'
            }
            }
        }

        stage('Release') {
            when {
                tag 'v*'
            }
            environment {
                GITHUB_TOKEN = credentials('github_token')
            }
            steps {
                sh 'curl -sL https://git.io/goreleaser | bash'
            }
        }
    }
}