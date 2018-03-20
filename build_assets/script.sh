#!/bin/bash

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
