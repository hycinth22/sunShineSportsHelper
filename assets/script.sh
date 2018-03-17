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

# create function and export it
# function signature: exec(user)
function execJKWX()
{
user=$1
pwd=$2
$main -login -u $user -p $pwd &>> $logFile
sleep 5s
# $main -status -u $user &>> $logFile
$main -upload -u $user -q &>> $logFile
echo "----------------" &>> $logFile
}
export main
export logFile
export -f execJKWX

# begin working
echo -e "\n" >> $logFile
echo "开始执行" &>> $logFile
echo "Date:" + `date` &>>$logFile
echo "----------------" &>> $logFile
# read each and call execJKWX
awk '{print $1 " " $2;"execJKWX " $1 " " $2|getline result;print result}' < $wd/accounts.list

# recover previous workDir
popd