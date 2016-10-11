#!/usr/bin/env bash

csvName=$1
gnpName=$2
getfile=$3


ab -k -n 5000 -c 10 -e ${csvName} -g ${gnpName} ${getfile}