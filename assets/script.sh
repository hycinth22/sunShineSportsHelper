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

# function signature: exec(user)
function execJKWX()
{
    user=$1
    pwd=$2
    dis=$3
    echo $user + `date` &>>$logFile
    $main -u $user -login -p $pwd &>> $logFile
    randomSleep 5 15
    # $main -status -u $user &>> $logFile
    $main -u $user -upload -q -distance $dis &>> $logFile
    echo "----------------" &>> $logFile
	randomSleep 20 60
}
# function signature: rand(min, max)
function rand(){  
    min=$1  
    max=$(($2-$min+1))  
    num=$(cat /dev/urandom | head -n 10 | cksum | awk -F ' ' '{print $1}')  
    echo $(($num%$max+$min))  
}  
# function signature: randomSleep(min, max)
function randomSleep()
{
    min=$1  
    max=$2
    rnd=$(rand $min $max)
	echo "Sleep ${rnd}s" &>> $logFile
	sleep ${rnd}'s'
}

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