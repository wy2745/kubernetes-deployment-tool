#!/usr/bin/env bash

IP=192.168.6.10

kubectl=/home/administrator/kubernetes/cluster/ubuntu/binaries/kubectl

ssh ${IP} "cd /home/administrator/resources/ && ${kubectl} --kubeconfig=kubeconfig.loc delete rc nginx"