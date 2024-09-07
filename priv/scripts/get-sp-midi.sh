#!/bin/sh

SP_VERS_NUM=4.5.1
SP_VERS=v${SP_VERS_NUM}
SP_URL= https://github.com/sonic-pi-net/sonic-pi/archive/refs/tags/${SP_VERS}.zip
SP_DL_DIR=tmp
SP_UNZIP_DIR=sonic-pi-${SP_VERS_NUM}
SP_DIR=sp_midi
SP_BUILD_DIR=${SP_DIR}/build

function download () {
    echo "Downloading Sonic Pi MIDI NIF source ..."
    mkdir -p $SP_DL_DIR
    cd $SP_DL_DIR && \
        curl --silent -L -O $SP_URL && \
        unzip ${SP_VERS}.zip
}

function prep-build () {
    echo "Setting up MIDI NIF build dir ..."
    mkdir -p $SP_DIR
    mv ${SP_UNZIP_DIR}/app/external/sp_midi .
    mkdir -p $SP_BUILD_DIR
}

function build () {
    echo "Building MIDI NIF ..."
    cd $SP_BUILD_DIR && \
        cmake && \
        make
    cp ${SP_BUILD_DIR}/libsp_midi.dylib src
}

function post-build () {
    echo "Cleaning up MIDI NIF temporary and build directories ..."
    rm -rf $SP_DL_DIR $SP_DIR
}

download
pre-build
build
post-build

