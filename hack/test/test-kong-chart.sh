#!/bin/bash

# Update submodule to the latest commit
git submodule update --init --recursive --remote --merge

[[ `which yq` ]] && YQ=yq || \
    (
        wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O ./yq
        chmod +x ./yq
        YQ=./yq
    )

# Loop through each directory in charts/charts and run helm dependency update
for dir in charts/charts/*; do
    if [ -d "$dir" ]; then
        cd "$dir"
        $YQ e '.dependencies[] | .name + " " + .repository' Chart.lock | while read line; do 
            helm repo add $line --force-update
        done
        helm dependency build
        cd -
    fi
done

# Run make test.golden and check return code
make -C charts/ test.golden
if [ $? -eq 0 ]; then
    echo "Tests passed successfully."; exit 0
else
    echo "Tests failed."; exit 1
fi