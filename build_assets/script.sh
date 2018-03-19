#!/bin/bash

source ./common.sh

# begin working
echo -e "\n" >> $logFile
echo "开始执行" &>> $logFile
echo "Date:" + `date` &>>$logFile
echo "----------------" &>> $logFile

# read each and call execJKWX
while read user passwd distance
do
    echo $(execJKWX $user $passwd $distance)
    echo "$user $passwd $distance"
done < accounts.list

# recover previous workDir
popd