#!/bin/bash
input=".env"

# 1 is for running service
# 2 is for testing
if [ $1 == "1" ]
then
echo "--run_under='export "
while read -r l 
do echo "$l"
done < "$input";
echo " &&'";
elif [ $1 == "2" ]
then
while read -r l; do echo "--test_env $l"; done < "$input"
else
echo "not support"
fi
