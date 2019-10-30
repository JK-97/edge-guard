#!/bin/bash
set -e

package=$1
version=$2
user=$3

url="http://packages.debian.jiangxingai.com:8081/api/repos/user_$user/packages?q=$package%20(=%20$version)"
packages=$(curl -s --fail $url)
echo -e "Removing old packages from local repo user_$user: $packages"
data='{"PackageRefs": '$packages'}'
curl -s --fail -X DELETE \
  -H "Content-Type: application/json" \
  -d "$data" \
  "http://packages.debian.jiangxingai.com:8081/api/repos/user_$user/packages"
echo -e

url="http://packages.debian.jiangxingai.com:8081/api/repos/jxtoolset/packages?q=$package%20(=%20$version)"
packages=$(curl -s --fail $url)
echo -e "Adding packages to local repo user_$user: $packages"
data='{"PackageRefs": '$packages'}'
curl -s --fail -X POST \
  -H "Content-Type: application/json" \
  -d "$data" \
  "http://packages.debian.jiangxingai.com:8081/api/repos/user_$user/packages"
echo -e


echo -e "Updating published local repo user_$user"
curl -s -X PUT \
  -H "Content-Type: application/json" \
  -d "{}" \
  "http://packages.debian.jiangxingai.com:8081/api/publish/user_$user/stable"
echo -e
