pipeline {
  agent any
  stages {
    currentBuild.result = "SUCCESS"

    stage('Checkout') {
      checkout scm

      echo "Current build result: ${currentBuild.result}"
    }

    stage('Build Images') {
      environment {
        COVERALLS_TOKEN = "not-a-token"
      }
      sh 'printenv | sort > /tmp/env'
      sh 'make test-docker'
      sh 'make test-captplanet-docker'
      sh 'make images'

      echo "Current build result: ${currentBuild.result}"
    }

    stage('Build Binaries') {
      sh 'make binaries'

      echo "Current build result: ${currentBuild.result}"
    }

    if (env.BRANCH_NAME == "master") {
      stage('Push Images') {
        echo 'Push to Repo'
        sh 'make push-images'
        echo "Current build result: ${currentBuild.result}"
      }

      /* This should only deploy to staging if the branch is master */
      stage('Deploy') {
        sh 'make deploy'
        sh 'make binaries-upload'
        echo "Current build result: ${currentBuild.result}"
      }
    }

  }
  post {
    error {
      echo "Caught errors! ${err}"
      echo "Setting build result to FAILURE"
      currentBuild.result = "FAILURE"
    }
    always {
      sh 'make test-docker-clean clean-images'
      deleteDir()
    }
  }
}
