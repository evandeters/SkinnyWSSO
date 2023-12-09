/* Requires the Docker Pipeline plugin */
pipeline {
    agent { docker { image 'golang:latest' } }
    stages {
        stage('build') {
            steps {
                withCredentials([string(credentialsId: 'wsso_admin_creds', usernameVariable: 'WSSO_ADMIN_USERNAME', passwordVariable: 'WSSO_ADMIN_PASSWORD')]) {
                    sh '''
                        chmod +x ./install.sh
                        ./install.sh ${env.WSSO_ADMIN_USERNAME} ${env.WSSO_ADMIN_PASSWORD} ${env.WSSO_ADMIN_PASSWORD}
                    '''
                }
            }
        }
    }
}
