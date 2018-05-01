#!/bin/bash
set -eu

/home/server &
ID=$! # ID of webserver process, so we can kill it

tests_passed=true
expected="Hello From Adidas."
output=$(curl -s localhost:8080)
if [[ $output == *"$expected"* ]]; then
  echo "Test Success"
else
  echo "Test Failure"
  echo "$expected != $output"
  tests_passed=false
fi


kill $ID

if [[ "$tests_passed" == "true" ]]; then
  echo "Passed Tests"
else 
  echo "Failed Tests"
  exit 1
fi



