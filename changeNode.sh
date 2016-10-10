#!/usr/bin/env bash

nodenum=$1

IP=192.168.6.10

ssh ${IP} "cd /home/administrator/resources/ && kubectl delete rc hollow-node --namespace=kubemark && kubectl create -f hollow-node"+"${nodenum}.json --namespace=kubemark"

