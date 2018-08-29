node('node') {
  environment {
    COVERALLS_TOKEN = credentials('COVERALLS_TOKEN')
  }
  try {
    currentBuild.result = "SUCCESS"

    stage('Checkout') {
      checkout scm

      echo "Current build result: ${currentBuild.result}"
    }

    stage('Build Images') {
      sh 'make test-docker test-captplanet-docker images'

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
  catch (err) {
    echo "Caught errors! ${err}"
    echo "Setting build result to FAILURE"
    currentBuild.result = "FAILURE"

    /*
    mail body: "project build error is here: ${env.BUILD_URL}" ,
      from: 'build@storj.io',
      replyTo: 'build@storj.io',
      subject: 'project build failed',
      to: "${env.CHANGE_AUTHOR_EMAIL}"

      throw err
    */

  }
  finally {

    stage('Cleanup') {
      sh 'make test-docker-clean clean-images'
      deleteDir()
    }

  }
}
