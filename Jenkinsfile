#!/usr/bin/env groovy

pipeline{
    //agent any
    agent any

    stages{
        stage("Test"){
            steps{
                make test
            }
        }
        stage("Build"){
            steps{
                make build
            }
        }
        stage("Create Image"){
            steps{
                make docker-build
            }
        }
        stage("Push Image"){
            steps{
                make docker-push
            }
        }
        stage("Archive"){
            steps{
                echo "Archive and cleanup..."
            }
        }
    }
    post {
        always {
            echo 'This will always run'
            slackSend   channel: '#jenkins',
                        color: 'good',
                        message: "The pipeline ${currentBuild.fullDisplayName} completed successfully."
        }
        success {
            echo 'This will run only if successful'
        }
        failure {
            echo 'This will run only if failed'
        }
        unstable {
            echo 'This will run only if the run was marked as unstable'
        }
        changed {
            echo 'This will run only if the state of the Pipeline has changed'
            echo 'For example, if the Pipeline was previously failing but is now successful'
        }
    }
}
