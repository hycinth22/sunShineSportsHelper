#!/bin/bash

# save previous workDir and switch to directory of the script 
wd="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd $wd

source ./common.sh

sleep $1
execJKWX $2 $3 $4

# recover previous workDir
popd