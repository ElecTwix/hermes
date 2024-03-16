#!/bin/sh -l

echo "Hello $1"
time=$(date)
./usr/local/bin/app $1
echo "time=$time" >>$GITHUB_OUTPUT
