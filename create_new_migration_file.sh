#!/bin/bash
[ $# -ne 1 ] && echo "Usage: $0 \"<comment>\"" && exit 1
comment=`echo "$1" | sed 's/[^a-zA-Z0-9]/_/g' | awk '{print tolower($0)}'`
newfile=db/migrations/`date +"%Y%m%d%H%M%S"`_$comment.sql
echo "-- migrate:up" > $newfile
echo "-- write statements below this line" >> $newfile
echo "" >> $newfile
echo "" >> $newfile
echo "-- migrate:down" >> $newfile
echo "-- write rollback statements below this line" >> $newfile
echo "" >> $newfile
echo "Created ./$newfile"
