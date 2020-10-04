#!/bin/bash

package=$1
major=$2
minor=$3
patch=$4

if [[ -z "$package" || -z "$major" || -z "$minor" || -z "$patch" ]]; then
  echo "usage: $0 <package-name> <major-ver-num> <patch-ver-num> <patch-ver-num>"
  #  Example: ./build-releases.sh dvembed 1 0 8
  exit 1
fi

# add more architectures
platforms=("windows/amd64" "linux/amd64" "darwin/amd64")
zip=""

for platform in "${platforms[@]}"
do
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}

  path=releases/v"$major"/"$minor"/"$patch"/"$GOOS"/"$package"

  if [ $GOOS = "windows" ]; then
    path+='.exe'
  fi

  env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$path" "$package"

  if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting the script execution...'
    exit 1
  fi

  zip+="$GOOS"" "
done

tar -czf releases/v"$major"/"$minor"/"$patch"/dvembed.tar.gz -C releases/v"$major"/"$minor"/"$patch" $zip
#echo tar -czf releases/v"$major"/"$minor"/"$patch"/dvembed.tar.gz "$zip"
