#!/bin/bash

#usage: serialtest.sh <exe>

iterations=$1
shift
cmd=$*
status=0

bash testbin/loop.sh ${iterations} ${cmd} -SYNVER
SYNVER_status=$?
if [ $SYNVER_status != 0 ]; then
    status=1
fi

bash testbin/loop.sh ${iterations} ${cmd} -SAVESYN test.syn
SAVESYN_status=$?
if [ $SAVESYN_status != 0 ]; then
    status=1
fi

bash testbin/loop.sh ${iterations} ${cmd} -LOADSYN test.syn
LOADSYN_status=$?
if [ $LOADSYN_status != 0 ]; then
    status=1
fi

bash testbin/loop.sh ${iterations} ${cmd} -LOADCRT data/testfiles/INTERNAL.CRT
LOADCRT_INTERNAL_status=$?
if [ $LOADCRT_INTERNAL_status != 0 ]; then
    status=1
fi

bash testbin/loop.sh ${iterations} ${cmd} -LOADCRT data/testfiles/L4.CRT
LOADCRT_L4_status=$?
if [ $LOADCRT_L4_status != 0 ]; then
    status=1
fi

bash testbin/loop.sh ${iterations} ${cmd} -LOADCRT data/testfiles/VCART1.CRT
LOADCRT_VCART1_status=$?
if [ $LOADCRT_VCART1_status != 0 ]; then
    status=1
fi

function show () {
    if [ $1 -eq 0 ];then
	echo PASS: test $2 
    else
	echo FAIL: test $2 
    fi
}

echo ========================
show $SYNVER_status  SYNVER
show $SAVESYN_status  SAVESYN
show $LOADSYN_status  LOADSYN
show $LOADCRT_INTERNAL_status  LOADCRT_INTERNAL
show $LOADCRT_L4_status  LOADCRT_L4
show $LOADCRT_VCART1_status  LOADCRT_VCART1

if [ $status -eq 0 ]; then
    echo overall PASS
else
    echo overall FAIL
fi
exit $status
