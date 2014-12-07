#!/bin/bash

OUTFILE='res/outfile/monkey.jpg'
INFILE='res/infile/infile.jpg'
T_WIDTH=16
T_HEIGHT=16
O_HEIGHT=2400
O_WIDTH=1600

./imgtitles -overwrite -dir res/infiles/ -tile_width $T_WIDTH -tile_height $T_HEIGHT -in $INFILE -out $OUTFILE -height $O_HEIGHT -width $O_WIDTH

if [ -f $OUTFILE ]
then
    okular $OUTFILE
fi;
