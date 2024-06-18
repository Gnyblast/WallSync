#!/bin/bash

newVersion=$1
currentDir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
versioningFile="${currentDir}/../internal/constants/version.go"

if [ -z "$newVersion" ]; then
    echo "Please pass version as 1st parameter to this script."
    exit 1
fi

echo "New version: ${newVersion}"
versionNumber=$(grep VERSION "${versioningFile}" | cut -d '=' -f 2)
sed -i "s/${versionNumber}/ \"v${newVersion}\"/" "${versioningFile}"
