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

    tools { go '1.21.5'}

    stages {
        stage('Compile') {
            steps {
                sh 'go build'
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
                        ssh -o StrictHostKeyChecking=no skinnywsso-dev 'ldapadd -x -D cn=admin,dc=skinny,dc=wsso -w $WSSO_ADMIN_PSW -H ldapi:/// -f /opt/skinnywsso/wsso.ldif'
                        ssh -o StrictHostKeyChecking=no skinnywsso-dev 'cd /opt/skinnywsso/; nohup ./SkinnyWSSO > /opt/skinnywsso/log &'
                    '''
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