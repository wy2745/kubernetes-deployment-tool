#!/usr/bin/env bash


for hatch_rate in 10 100 1000 10000
do
    for locust_count in 10 100 1000 10000
    do
        ./kubernetes-deployment-tool -l 3 ${locust_count} ${hatch_rate}
        done
    done