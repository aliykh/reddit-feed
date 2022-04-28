def buildDeployNx() {
    echo "building image and uploading it into docker hosted repo on the Nexus server"
    // withCredentials([usernamePassword(credentialsId: 'nexus-user', passwordVariable: 'PASS', usernameVariable: 'USER')]) {
    //     sh 'docker build --build-arg PROXY=http://172.17.0.1:8081/repository/go-proxy/ -t 0.0.0.0:8083/reddit-feed:1.1 .'
    //     sh "echo $PASS | docker login -u $USER --password-stdin 0.0.0.0:8083"
    //     sh 'docker push 0.0.0.0:8083/reddit-feed:1.1'
    // }

    sh 'docker build --build-arg PROXY=$NEXUS_GO_PROXY -t $NEXUS_DOCKER_HOST/reddit-feed:1.1 .'
    sh 'echo $NEXUS_RM_CREDS_PSW | docker login -u $NEXUS_RM_CREDS_USR --password-stdin $NEXUS_DOCKER_HOST'
    sh 'docker push $NEXUS_DOCKER_HOST/reddit-feed:1.1'

}

def buildDeployDocker() {
      echo "building image and uploading it into docker hub"
    // withCredentials([usernamePassword(credentialsId: 'docker-hub', passwordVariable: 'PASS', usernameVariable: 'USER')]) {
    //     sh 'docker build --build-arg PROXY=http://172.17.0.1:8081/repository/go-proxy/ -t alioy/reddit-feed:1.1 .'
    //     sh "echo $PASS | docker login -u $USER --password-stdin"
    //     sh 'docker push alioy/reddit-feed:1.1'
    // }

    sh 'docker build --build-arg PROXY=$NEXUS_GO_PROXY -t alioy/reddit-feed:1.1 .'
    sh 'echo $DOCKER_HUB_CREDS_PSW | docker login -u $DOCKER_HUB_CREDS_USR --password-stdin'
    sh 'docker push alioy/reddit-feed:1.1'
}

return this