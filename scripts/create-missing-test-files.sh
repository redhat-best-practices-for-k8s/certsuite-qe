#!/usr/bin/env bash

# Lists only files with package name.
file_names=$(
	find ./tests \
		-name "*.go" \
		! -name "*_test.go" \
		! -name "*_moq.go" \
		! -name "*ginkgo*" |
		sed 's/.go//g'
)

# Creates the file if it doesn't exist.
skp=0
add=0
for file_name in $file_names; do
	[ -f "$file_name"_test.go ] && {
		skp=$((skp + 1))
		continue
	}
	package_name=$(grep -m 1 package "$file_name".go | awk '{print $2}')
	echo "package $package_name" >>"$file_name"_test.go
	add=$((add + 1))
done
printf >&2 'Added %d and skipped %d file(s).\n' $add $skp
