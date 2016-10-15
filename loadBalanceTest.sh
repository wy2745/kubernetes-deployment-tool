#!/usr/bin/env bash

echo "开始测试"

teserIP=192.168.6.15
nginxIp=192.168.6.31
serviceUrl=http://192.168.6.22:30080/index.html
apiUrl=http://192.168.6.10:8080/api/v1/proxy/namespaces/default/services/nginx-svc:80/index.html
eloadbUrl=http://192.168.6.31/index.html
fileroot=lbtest

kubectl=/home/administrator/kubernetes/cluster/ubuntu/binaries/kubectl



for podnum in 1 2 4 8 16 32
do
    echo "replic 数量:${podnum}"
    kubectl scale --replicas=${podnum} rc/nginx
    while [ "$(kubectl get pods | grep nginx | wc -l)" != "${podnum}" ];do
        sleep 1
    done
    servicefileName=service${podnum}p
    apifileName=apiserver${podnum}p
    eloadfileName=eloadbalance${podnum}p


    sleep 10
    echo "进行server的ab测试"
    ssh ${teserIP} "cd ${fileroot} && ab -k -n 100000 -c 50 -e ${servicefileName}.csv -g ${servicefileName}.gnp ${serviceUrl} > ${servicefileName}.html"
    sleep 60
    echo "测试完成"
    echo "进行apiserver的ab测试"
    ssh ${teserIP} "cd ${fileroot} && ab -k -n 100000 -c 50 -e ${apifileName}.csv -g ${apifileName}.gnp ${apiUrl} > ${apifileName}.html"
    sleep 60
    echo "测试完成"
    arr=($(kubectl get pods | grep nginx | awk '{print $1};' | tail -n +1))
    server=""
    for s in ${arr}
    do
        server+="\tserver "$(kubectl describe pod ${s} | grep Node | awk '{{split($2,a,"/" ); print a[1]}}')":8888;\n"
        done
    echo "对nginx进行配置"
    ssh ${nginxIp} "cd /etc/nginx/conf.d && ./nginxProxy.sh \"${server}\" && echo incongruous | sudo -S service nginx restart"
    echo "配置完成"
    sleep 10
    echo "进行外部loadbalancer的ab测试"
    ssh ${teserIP} "cd ${fileroot} && ab -k -n 100000 -c 50 -e ${eloadfileName}.csv -g ${eloadfileName}.gnp ${eloadbUrl} > ${eloadfileName}.html"
    sleep 60
    echo "测试完成"
    done

echo "测试完成,准备删除rc"
kubectl delete rc nginx
while [ "$(kubectl get pods | grep nginx | wc -l)" != "0" ];do
        sleep 1
    done
kubectl delete svc nginx-svc
echo "完成"