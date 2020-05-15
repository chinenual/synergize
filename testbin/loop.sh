#!/bin/bash

# usage:
#    loop.sh <N> <command and args>

n=$1
shift
cmd=$*

passcount=0
failcount=0
i=0

while [ $i -lt ${n} ]; do
    i=$((i + 1))
    echo
    echo ${cmd}
    echo
    ${cmd}
    if [ $? != 0 ]; then
	echo $i: FAIL
	failcount=$((failcount + 1))
    else 
	echo $i: PASS
	passcount=$((passcount + 1))
    fi
done

echo $passcount PASS, $failcount FAIL : ${cmd}

if [ $failcount != 0 ]; then
    exit 1
fi
exit 0
