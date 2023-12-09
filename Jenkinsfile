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
                    sh 'ssh -o StrictHostKeyChecking=no skinnywsso-dev ls'
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