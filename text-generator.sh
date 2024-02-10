#!/bin/bash

# Define the file name
file="entries.txt"

# Generate 1 million entries and write them to the file
for ((i=1; i<=1000000; i++)); do
    echo "Entry $i" >> "$file"
done

echo "File generation complete."
