#!/bin/bash
if [ "$1" = "cover" ]; then
  echo "Running Tests With Coverage..."

  go test ./domain ./usecase ./handler ./repository -v -covermode=atomic -coverprofile=cover.out
  go tool cover -html=cover.out -o cover.html
else
  echo "Running Tests..."

  go test ./domain ./usecase ./handler ./repository -v
fi