#!/usr/bin/env groovy

pipeline{
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
}
