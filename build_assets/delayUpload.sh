#!/bin/bash

# save previous workDir and switch to directory of the script 
wd="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd $wd

# files
main=$wd/jkwx
logFile=$wd/log/script_`date +%Y-%m-%d`.log

# create files
touch $logFile
chmod 777 $logFile

sleep $2
$main -u $1 -upload -q

# recover previous workDir
popd