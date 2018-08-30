pipeline {
  agent any
  stages {
    stage('Checkout') {
      steps {
        checkout scm

        echo "Current build result: ${currentBuild.result}"
      }
    }

    stage('Build Images') {
      environment {
        COVERALLS_TOKEN = "not-a-token"
      }
      steps {
        sh 'printenv | sort > /tmp/env'
        sh 'make test-docker'
        sh 'make test-captplanet-docker'
        sh 'make images'

        echo "Current build result: ${currentBuild.result}"
      }
    }

    stage('Build Binaries') {
      sh 'make binaries'

      echo "Current build result: ${currentBuild.result}"
    }

    stage('Push Images') {
      when {
        branch 'master'
      }
      steps {
        echo 'Push to Repo'
        sh 'make push-images'
        echo "Current build result: ${currentBuild.result}"
      }
    }

    /* This should only deploy to staging if the branch is master */
    stage('Deploy') {
      when {
        branch 'master'
      }
      steps {
        sh 'make deploy'
        echo "Current build result: ${currentBuild.result}"
      }
    }

  }
  post {
    failure {
      steps {
        echo "Caught errors! ${err}"
        echo "Setting build result to FAILURE"
      }
    }
    cleanup {
      steps {
        sh 'make test-docker-clean clean-images'
        deleteDir()
      }
    }
  }
}
