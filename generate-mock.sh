#!/bin/bash

if [[ $1 != "" ]]; then

for d in $(find . -path '*/'"$1"'/*' -name '*.go' -not -path '*/mock/*' ! -name '*_test.go' ! -name '*_mock.go' ! -name 'base.go' ! -name 'query.go'); do
    fileName="${d##*/}"
    dir="$(dirname "$d")"

    # Check if the file contains an interface definition
    if grep -q "interface" "$d"; then
        name=${fileName%.*}
        fileGen="${name}_mock.go"
        destination="${dir}/mock"

        # Create the mock directory if it doesn't exist
        mkdir -p "$destination"

        printf "generating mock for $fileGen ...\n"

        # Generate the mock using mockgen
        mockgen -source="$d" -destination="$destination/$fileGen" -package=mock

        printf "finished\n"
    fi

done
printf "mock $1 successfully generated ...\n"

else

printf "directory must be filled, no mock generated ...\n"

fi
