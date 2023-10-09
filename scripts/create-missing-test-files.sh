#!/usr/bin/env bash

file_names=$(find ./tests -name "*.go" | grep -v "_test.go" | grep -v "_moq.go" | grep -v "ginkgo" | sed 's/.go//g')

# if a filename_test.go does not exist, create it
for file_name in $file_names; do
	if [ ! -f "$file_name"_test.go ]; then
		package_name=$(grep -m 1 "package" "$file_name".go | awk '{print $2}')
		echo "package $package_name" >>"$file_name"_test.go
	fi
done
