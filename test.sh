#!/bin/bash

OUTFILE='res/outfile/monkey.jpg'
INFILE='res/infile/infile.jpg'

./imgtitles -dir res/infiles/ -in $INFILE -out $OUTFILE -height 800 -width 1200

if [ -f $OUTFILE ]
then
    okular $OUTFILE
fi;
