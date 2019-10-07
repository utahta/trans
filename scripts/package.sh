#!/bin/bash

curdir=$(dirname $0)
distdir=${curdir}/../dist

source ${curdir}/common.sh

export GOFLAGS="-trimpath"
gox -ldflags="-X main.Version=$(version)" github.com/utahta/trans/cmd/trans/

rm -rf ${distdir}
mkdir -p ${distdir}
mv trans_* ${distdir}/

cd ${distdir}
for name in trans_*; do
    if [ "${name#*.}" = "exe" ]; then
        command="trans.exe"
    else
        command="trans"
    fi

    mv "${name}" "${command}"
    zip "${name}.zip" "${command}"
    rm "${command}"
done
