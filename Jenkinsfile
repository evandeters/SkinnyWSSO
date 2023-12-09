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
        /**stage('LDAP') {
            steps {
                    sh '''
                        rm -rf /var/lib/ldap/* 
                        cp -R /root/ldap_backup/* /var/lib/ldap/
                        chown -R openldap:openldap /var/lib/ldap/
                        systemctl restart slapd
                        ldapadd -x -w $WSSO_ADMIN_PSW -H ldapi:/// -D cn=admin,dc=skinny,dc=wsso -f ./wsso.ldif
                        echo test
                    '''
                }
        }
    **/

        stage('Compile') {
            steps {
                sh 'go build'
            }
        }

        stage('Release') {
            when {
                buildingTag()
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
