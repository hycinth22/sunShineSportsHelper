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

# function signature: exec(user, pwd, dis)
function execJKWX()
{
    user=$1
    pwd=$2
    dis=$3
    log="${logFile}_${user}"
    echo $user + `date` | tee -a ${log}
    $main -u $user -login -p $pwd | tee -a ${log}
    randomSleep 15 30
    # $main -status -u $user | tee -a ${log}
    $main -u $user -upload -q -distance $dis | tee -a ${log}
    echo "----------------" &>> ${log}
	randomSleep 15 360
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
	echo "Sleep ${rnd}s" |tee -a $logFile
	sleep ${rnd}'s'
}
