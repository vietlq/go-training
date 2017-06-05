#!/usr/bin/bash

if [ -z $GOPATH ]; then
    export GOPATH=`pwd`:$HOME/projects/gocode
else
    export GOPATH=$GOPATH:`pwd`:$HOME/projects/gocode
fi
