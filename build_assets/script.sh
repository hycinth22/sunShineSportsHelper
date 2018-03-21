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
while read user passwd distance
do
    echo "$user $passwd $distance"
      execJKWX $user $passwd $distance
done < accounts.list

# recover previous workDir
popd
