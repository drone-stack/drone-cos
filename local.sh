#!/bin/bash

go build -o ./dist/cos cmd/main.go

./dist/cos --bucket pkg --accesskey <xxx> --secretkey <xxxx> --endpoint cos.ap-shanghai.myqcloud.com --source /Users/ysicing/go/src/gitea.svc.gw/pangu/qcadmin/dist  --target tools/cli/test/ --autotime --timeformat 20060102 