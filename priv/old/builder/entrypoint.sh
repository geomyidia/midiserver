#!/bin/bash

set -e

echo "Got environment:"
env

if [[ -z "$GITHUB_WORKSPACE" ]]; then
  echo "Set the GITHUB_WORKSPACE env variable."
  exit 1
fi
echo "----> Using GITHUB_WORKSPACE $GITHUB_WORKSPACE"
ls -al $GITHUB_WORKSPACE

if [[ -z "$GITHUB_REPOSITORY" ]]; then
  echo "Set the GITHUB_REPOSITORY env variable."
  exit 1
fi
echo "----> Using GITHUB_REPOSITORY $GITHUB_REPOSITORY"

root_path="/go/src/github.com/$GITHUB_REPOSITORY"
base_path="$(dirname $root_path)"
release_path="$root_path/builds"
repo_name="$(echo $GITHUB_REPOSITORY | cut -d '/' -f2)"
targets=${@-"darwin/amd64 darwin/386arm64 linux/amd64 linux/386 windows/amd64 windows/386"}

echo "----> Using root_path $root_path"
echo "----> Using release_path $release_path"
echo "----> Using repo_name $repo_name"
echo "----> Using GOPATH $GOPATH"
echo "----> Setting up Go repository"
mkdir -p $base_path

echo "cp -a $GITHUB_WORKSPACE to $base_path"
cp -a $GITHUB_WORKSPACE $root_path
mkdir -p $release_path
echo "cd to $base_path"
cd $base_path
ls -al ./
echo "cd $root_path"
cd $root_path
ls -al ./

for target in $targets; do
  os="$(echo $target | cut -d '/' -f1)"
  arch="$(echo $target | cut -d '/' -f2)"
  output="${release_path}/${repo_name}-${os}-${arch}"

  echo "----> Building project for: $target"
  echo "current directory:"
  pwd
  echo "directory listing:"
  ls -al 
  # GOOS=$os GOARCH=$arch CGO_ENABLED=0 OUTPUT=$output make cross-compile
  OUTPUT=$output make cross-compile

  if [[ -n "$COMPRESS_FILES" ]]; then
    zip -j $output.zip $output > /dev/null
    rm $output
  fi
done

echo "----> Build is complete. List of files at $release_path:"
cd $release_path
ls -al
cp -a $release_path $GITHUB_WORKSPACE
