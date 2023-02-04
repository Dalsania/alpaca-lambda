#!/bin/bash


GOOS=linux go build -o main alpaca.go
echo Linux go build completes
zip function.zip main
echo Zipping build completes
# echo Hello $1

# while getopts u:a:f: flag
# do
#     case "${flag}" in
#         u) username=${OPTARG};;
#         a) age=${OPTARG};;
#         f) fullname=${OPTARG};;
#     esac
# done
# echo "Username: $username";
# echo "Age: $age";
# echo "Full Name: $fullname";