#!/bin/bash

# save previous workDir and switch to directory of the script 
wd="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd $wd

source ./common.sh

# begin working
echo -e "\n" | tee -a $logFile
echo "开始执行" | tee -a $logFile
echo "Date:" + `date` | tee -a $logFile
echo "----------------" | tee -a $logFile

# read each and call execJKWX
while read school user passwd distance
do
    echo "$school $user $passwd $distance"|tee -a $logFile
    execJKWX $school $user $passwd $distance
done < accounts.list

# recover previous workDir
popd
