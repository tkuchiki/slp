#!/bin/bash

bin="slp"
version="${1}"

set -e

cd build
for file in ./${bin}-* ; do
    dir=$(echo ${file} | awk -F'-' -v bin=${bin} '{print bin"_"$(NF-1)"_"$(NF)}')
    mkdir ${dir}
    mv ${file} ${dir}/${bin}
    # zip
    zip ${dir}.zip -j ${dir}/${bin} ../README.md ../LICENSE
    # tar
    cd ${dir}
    cp ../../README.md .
    cp ../../LICENSE .
    tar czf ${dir}.tar.gz ${bin} README.md LICENSE
    mv ${dir}.tar.gz ../
    cd ../
    rm -rf ${dir}
done

shasum -a 256 *.zip *.tar.gz > ${bin}_${version}_checksums.txt
