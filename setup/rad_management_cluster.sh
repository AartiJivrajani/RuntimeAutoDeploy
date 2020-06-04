#bin/bash
set -e

echo -e "################################################ Installing Kubectl ################################################ \n"

sudo su

#curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/darwin/amd64/kubectl"

curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl

#kubectl version


#if kubectl version
#then
#    echo -e "kubectl version is :" kubectl version | awk '{print $1 "\t" $5}' | head -n 1
#    exit 0
#else
#	echo "kubectl not installed. Exiting"
#	exit 1
#fi


echo -e "****************************************************************************************************************** \n"

echo -e "################################################ Installing Kind ################################################# \n"

curl -Lo ./kind "https://kind.sigs.k8s.io/dl/v0.8.0/kind-$(uname)-amd64"
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
kind version

echo -e "****************************************************************************************************************** \n"

echo -e "############################################ Installing Docker Engine ############################################ \n"

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
apt-cache policy docker-ce
sudo apt-get install -y docker-ce
#sudo systemctl status docker

echo -e "****************************************************************************************************************** \n"

echo -e "############################################### Create Kind Cluster ################################################ \n"

kind create cluster
kubectl cluster-info --context kind-kind
#Run "kind delete cluster", to re-run shell again

echo -e "****************************************************************************************************************** \n"

echo -e "############################################# Installing Clusterctl ############################################## \n"

## for mac:
## curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.3.5/clusterctl-darwin-amd64 -o clusterctl
curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.3.5/clusterctl-linux-amd64 -o clusterctl
chmod +x ./clusterctl
sudo mv ./clusterctl /usr/local/bin/clusterctl
clusterctl version

echo -e "****************************************************************************************************************** \n"

echo -e "###################################### Initialize the management cluster ########################################### \n"

clusterctl init

echo -e "****************************************************************************************************************** \n"

echo -e "############################################ Install clusterawsadm ################################################# \n"

# for mac OS
# curl -L https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/v0.5.3/clusterawsadm-darwin-amd64 -o clusterawsadm
curl -L https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/v0.5.3/clusterawsadm-linux-amd64 -o clusterawsadm
chmod +x clusterawsadm
sudo mv ./clusterawsadm /usr/local/bin/clusterawsadm
clusterawsadm version

echo -e "****************************************************************************************************************** \n"

