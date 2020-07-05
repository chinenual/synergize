#!/bin/bash

set -x

iterations=${1:-100}

bash testbin/serialtests.sh ${iterations} output/windows-386/Synergize-cmd.exe -port COM1 -baud 9600 2>&1 | tee serialtests-windows.log

