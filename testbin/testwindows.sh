#!/bin/bash

iterations=${1:-20}

bash testfiles/bin/serialtests.sh ${iterations} output/windows-386-cmd/Synergize-cmd.exe -port COM1 -baud 9600 2>&1 | tee serialtests-windows.log

