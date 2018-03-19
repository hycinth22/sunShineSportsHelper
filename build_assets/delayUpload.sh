#!/bin/bash

source ./common.sh

sleep $1
execJKWX $2 $3 $4

# recover previous workDir
popd