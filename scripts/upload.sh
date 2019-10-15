#!/bin/bash

if [ ! -d "build" ]; then
  echo -e "error! build not exist"
  exit 1
fi

for filename in build/*.deb; do
  echo -e "uploading $filename"
  curl -s -X POST -F "file=@"$filename "http://packages.debian.jiangxingai.com:8081/api/files/jxtoolset"
  echo -e
  curl -s -X POST "http://packages.debian.jiangxingai.com:8081/api/repos/jxtoolset/file/jxtoolset?forceReplace=1"
  echo -e
done

curl -s -X PUT \
  -H "Content-Type: application/json" \
  -d '{"SourceKind": "local", "Sources": [{"Name": "jxtoolset"}], "ForceOverwrite": true}' \
  "http://packages.debian.jiangxingai.com:8081/api/publish/local_jxtoolset/stable"

echo -e