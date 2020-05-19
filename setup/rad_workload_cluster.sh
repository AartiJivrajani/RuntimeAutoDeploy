#bin/bash
set -e

echo -e "###################################### Initialization For Common Providers ######################################### \n"

# export these as ENVIRONMENT variables
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=AKIAUYHYI3WA6IZSPBHAYL6L4
export AWS_SECRET_ACCESS_KEY=NtNRlP6WO

clusterawsadm alpha bootstrap create-stack

echo "#############Sleeping for 120 seconds###########"
sleep 120

export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm alpha bootstrap encode-aws-credentials)
clusterctl init --infrastructure aws

sleep 300

echo -e "****************************************************************************************************************** \n"


echo -e "###################################### Creating First Workload Cluster ############################################# \n"

export AWS_REGION=us-east-1
export AWS_CONTROL_PLANE_MACHINE_TYPE=t2.medium
export AWS_SSH_KEY_NAME=aarti
export AWS_NODE_MACHINE_TYPE=t2.medium

clusterctl config cluster capi-quickstart --kubernetes-version v1.18.0 --control-plane-machine-count=1 --worker-machine-count=1 > capi-quickstart.yaml
kubectl apply -f capi-quickstart.yaml


echo -e "****************************************************************************************************************** \n"

echo -e "###################################### Accessing The Workload Cluster ############################################## \n"

kubectl get cluster --all-namespaces
kubectl get kubeadmcontrolplane --all-namespaces

echo "######Sleeping for 20 minutes for making sure all the resources to come to ready state#######"

sleep 1200

echo -e "****************************************************************************************************************** \n"

echo -e "################################# Retrieving The Workload Cluster KubeConfig ####################################### \n"

kubectl --namespace=default get secret/capi-quickstart-kubeconfig -o jsonpath={.data.value} \
  | base64 --decode \
  > ./capi-quickstart.kubeconfig

echo -e "****************************************************************************************************************** \n"

echo -e "######################################### Deploying a CNI solution ################################################## \n"

kubectl --kubeconfig=./capi-quickstart.kubeconfig \
  apply -f https://docs.projectcalico.org/v3.12/manifests/calico.yaml

sleep 120

echo "######Hold on tight, we are almost there :)#######"

kubectl --kubeconfig=./capi-quickstart.kubeconfig get nodes

#Deleting Workload cluster : kubectl delete -f capi-quickstart.yaml

echo -e "****************************************************************************************************************** \n"

kubectl get cluster --all-namespaces
kubectl get kubeadmcontrolplane --all-namespaces

echo -e "########************************************* WORKLOAD CLUSTER IS READY ***********************************########## \n"
