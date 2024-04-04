#!/bin/bash

# Update submodule to the latest commit
git submodule update --init --recursive

# Loop through each directory in charts/charts and run helm dependency update
for dir in charts/charts/*; do
    if [ -d "$dir" ]; then
        cd "$dir"
        helm dependency update
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