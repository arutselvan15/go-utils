#!/usr/bin/env groovy

pipeline{
    agent { 
        docker { 
            image 'arutselvan15/golang:1.12.4'
        }
    }
    environment {
        GOCACHE = '/tmp'
    }
    stages{
        stage("Checkout"){
            steps{
                echo "Git checkout ..."
                checkout([
                    $class: 'GitSCM', 
                    branches: [[name: '*/master']
                ],
                userRemoteConfigs: [[
                    url: 'https://github.com/arutselvan15/go-utils.git']]
                ])
            }
        }
        stage("Test"){
            steps{
                echo "Run Test ..."
                sh "make test"
            }
        }
        stage("Build"){
            steps{
                echo "Run Build ..."
            }
        }
        stage("Create Image"){
            steps{
                echo "Create Image ..."
            }
        }
        stage("Push Image"){
            steps{
                echo "Push Image ..."
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
            echo 'Build End'
        }
        success {
            echo 'Build completed successfully'
            slackSend(channel: '#jenkins', color: 'good', message: "The pipeline ${currentBuild.fullDisplayName} completed successfully.")
        }
        failure {
            echo 'Builld failure'
            slackSend(channel: '#jenkins', color: 'red', message: "The pipeline ${currentBuild.fullDisplayName} build failed.")
        }
        unstable {
            echo 'Build unstable'
            slackSend(channel: '#jenkins', color: 'red', message: "The pipeline ${currentBuild.fullDisplayName} is unstable.")
        }
        changed {
            echo 'Build state changed'
            slackSend(channel: '#jenkins', color: 'good', message: "The pipeline ${currentBuild.fullDisplayName} state changed.")
        }
    }
}
