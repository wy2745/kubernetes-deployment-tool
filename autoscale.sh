#!/usr/bin/env bash


IP=192.168.6.10
kubectl=/home/administrator/kubernetes/cluster/ubuntu/binaries/kubectl

ssh ${IP} "${kubectl} autoscale rc nginx --max=9 --min=3 --cpu-percent=80"