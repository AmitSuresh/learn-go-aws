#!/bin/bash
GOOS=linux go build -tags lambda.norpc -o bootstrap main.go
/c/"Program Files"/Go/bin/"build-lambda-zip.exe" -o myFunction.zip bootstrap