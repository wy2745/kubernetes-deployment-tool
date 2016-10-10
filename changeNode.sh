#!/usr/bin/env bash

nodenum=$1

IP=192.168.6.10
kubectl=/home/administrator/kubernetes/cluster/ubuntu/binaries/kubectl
echo "haha"

ssh ${IP} "cd /home/administrator/resources/ && ${kubectl} delete rc hollow-node --namespace=kubemark"
echo "${kubectl} create -f hollow-node"+"${nodenum}.json --namespace=kubemark"

