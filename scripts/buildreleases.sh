#!/usr/bin/bash

version="v0.1.0"

oss=(linux darwin windows)
archs=(amd64 arm64 arm)
cmds=(dmarcfetch dmarcsqltoxls)

cd ..
for os in ${oss[@]}
do
    for arch in ${archs[@]}
    do
        if [[ ${os}_${arch} == "windows_arm" ]]; then
            continue
        fi
        mkdir -p releases/buildtemp/${os}_${arch}/data
        for cmd in ${cmds[@]}
        do
            echo "Building ${cmd} for ${os} ${arch}"
            mkdir -p releases/buildtemp/${os}_${arch}/${cmd}

            pushd cmd/${cmd}

            env GOOS=${os} GOARCH=${arch} go build -o ../../releases/buildtemp/${os}_${arch}/${cmd}/${cmd}
            cp config.yml ../../releases/buildtemp/${os}_${arch}/${cmd}/config.yml
           
            popd
        done
        cp data/template.xlsx releases/buildtemp/${os}_${arch}/data/template.xlsx
        cp README.md releases/buildtemp/${os}_${arch}/README.md
        pushd releases/buildtemp/${os}_${arch}
        zip -r9 ../../${os}_${arch}_${version}_dmarcanalyze.zip *
        cd ..
        rm -rf ${os}_${arch}
        popd
    done
done
rm -rf releases/buildtemp