sudo su
curl -O https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz
tar xvf go1.14.3.linux-amd64.tar.gz
sudo chown -R root:root ./go
sudo mv go /usr/local
mkdir -p $HOME/work
export GOPATH=$HOME/work
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
mkdir -p $HOME/work/src
cd $HOME/work/src

# get all the dependencies

go get github.com/docker/docker/api/types
go get github.com/docker/docker/client
go get github.com/go-git/go-git
go get github.com/google/uuid
go get github.com/sirupsen/logrus
go get gopkg.in/redis.v5
go get k8s.io/api/apps/v1
go get k8s.io/api/core/v1
go get k8s.io/apimachinery/pkg/api/resource
go get k8s.io/apimachinery/pkg/apis/meta/v1
go get k8s.io/client-go/kubernetes
go get k8s.io/client-go/tools/clientcmd
go get k8s.io/client-go/util/homedir

# get the user repository
git clone https://github.com/shudaniel/RuntimeAutoDeploy.git

docker login

# redis

sudo apt update
sudo apt install redis-server -y
