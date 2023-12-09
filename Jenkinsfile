/* Requires the Docker Pipeline plugin */
pipeline {
    agent { docker { image 'golang:latest' } }
    environment {
        WSSO_ADMIN_USERNAME = credentials('wsso_admin_creds').username
        WSSO_ADMIN_PASSWORD = credentials('wsso_admin_creds').password
    }
    stages {
        stage('build') {
            steps {
                sh '''
                    chmod +x ./install.sh
                    ./install.sh ${env.WSSO_ADMIN_USERNAME} ${env.WSSO_ADMIN_PASSWORD} ${env.WSSO_ADMIN_PASSWORD}
                '''
            }
        }
    }
}
