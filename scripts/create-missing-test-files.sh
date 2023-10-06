#!/usr/bin/env bash
# shellcheck disable=SC2043
for folder in tests; do
	echo "./$folder"

	FILE_NAMES=$(find ./$folder -name "*.go" | grep -v "_test.go" | grep -v "_moq.go" | grep -v "ginkgo" | sed 's/.go//g')

	# if a filename_test.go does not exist, create it
	for FILE_NAME in $FILE_NAMES; do
		if [ ! -f "$FILE_NAME"_test.go ]; then
			touch "$FILE_NAME"_test.go

			PACKAGE_NAME=$(grep -m 1 "package" "$FILE_NAME".go | awk '{print $2}')
			echo "package $PACKAGE_NAME" >>"$FILE_NAME"_test.go
		fi
	done
done
rm -f {}_test.go
