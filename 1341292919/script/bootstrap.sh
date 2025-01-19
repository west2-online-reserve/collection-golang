#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=todolist
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}