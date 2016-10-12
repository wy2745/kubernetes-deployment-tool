#!/usr/bin/env bash

sudo sysctl -w net.ipv4.tcp_fin_timeout=5

./kubernetes-deployment-tool -dp
./kubernetes-deployment-tool -dn

for nodenum in 5 10 20 40 80 160 320
do
    echo "创建${nodenum}个node..."
    ./kubernetes-deployment-tool -cn ${nodenum}
    echo "创建成功"
    echo "进行${nodenum}个node的pod创建删除时间实验..."
    ./kubernetes-deployment-tool -t 10
    for cnt in 1 2 3
        do
            ./kubernetes-deployment-tool -cpt ${nodenum} ${cnt}
            ./kubernetes-deployment-tool -t 60
    done

    for replic in 3 5 10 15 20 30
    do
        echo "在${nodenum}个node上创建${nodenum}创建${replic}倍pod..."
        ./kubernetes-deployment-tool -cp ${nodenum} ${replic}
        echo "在${nodenum}个node上创建${nodenum}做${replic}倍pod实验..."
        for count in 1 2 3
            do
                ./kubernetes-deployment-tool -ab ${nodenum} ${replic} ${count}
                done
        ./kubernetes-deployment-tool -t 60
        echo "在${nodenum}个node上创建${nodenum}删除${replic}倍pod..."
        ./kubernetes-deployment-tool -dp
        ./kubernetes-deployment-tool -t 60
        echo "在${nodenum}个node上创建${nodenum}删除${replic}倍pod成功..."
        done
    echo "删除${nodenum}个node..."
    ./kubernetes-deployment-tool -dn
    echo "删除${nodenum}个node成功..."
    done