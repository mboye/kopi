#!/bin/bash
echo "Small file hash:"
cat salt store/small-file.txt | sha1sum

echo "Large file hash:"
cat salt store/large-file.txt | sha1sum

echo "Large file hash 1:"
cat salt store/small-file.txt | sha1sum

echo "Large file hash 2:"
tmp=$(mktemp)
cat salt > "$tmp"
dd if=store/large-file.txt skip=64 bs=1 >> "$tmp" 2>/dev/null
sha1sum < "$tmp"
