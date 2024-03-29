#!/bin/bash

# files
main=$wd/jkwx
today=`date +%Y-%m-%d`
logDir=$wd/log/$today
logFile=$logDir/summary.log
userLogPathPattern=${logDir}/%s.log

# prepare files
echo working at $wd
echo today $today
echo logDir $logDir
echo logFile $logFile
mkdir -p -m=777 $logDir
touch $logFile
chmod 664 $logFile

# function signature: exec(user, pwd, dis)
function execJKWX()
{
    school=$1
    user=$2
    pwd=$3
    dis=$4
	log=`printf ${userLogPathPattern} ${user}`
	echo user\'s logFile $log
    echo $user + `date` | tee -a ${log}

    $main -u $user -login -p $pwd -s $school | tee -a ${log}
    randomSleep 15 30
    # $main -status -u $user | tee -a ${log}
    $main -u $user -upload -q -distance $dis | tee -a ${log}
    echo "----------------" >> ${log}

	if [[ `cat ${log}` =~ "已达标" ]]; then
		echo $user + " done." >> ${log}.done
	fi
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
