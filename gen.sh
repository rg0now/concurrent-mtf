#!/bin/bash

#make

SRC=uniform
LB=modulo
CPULIM=110

# M=10
# N=5

####################
# MTF
####################
DS=mtf
echo $DS
M=10000
N=5000000

echo -- Multicore: m=$M n=$N d=$DS s=$SRC l=$LB 
for i in $(seq 1 48); do ./cmtf -m $M -n $N -d $DS -s $SRC -l $LB -k $i; done

echo -- Singlecore: m=$M n=$N d=$DS s=$SRC l=$LB cpu=$CPULIM
(for i in $(seq 1 48); do cpulimit -f -q -l $CPULIM -- ./cmtf -m $M -n $N -d $DS -s $SRC -l $LB -k $i; done) | grep -v "Child process"

####################
# CACHE
####################
DS=cache
echo $DS
M=10000
N=500000

echo -- Multicore: m=$M n=$N d=$DS s=$SRC l=$LB 
for i in $(seq 1 48); do ./cmtf -m $M -n $N -d $DS -s $SRC -l $LB -k $i; done

echo -- Singlecore: m=$M n=$N d=$DS s=$SRC l=$LB cpu=$CPULIM
(for i in $(seq 1 48); do cpulimit -f -q -l $CPULIM -- ./cmtf -m $M -n $N -d $DS -s $SRC -l $LB -k $i; done) | grep -v "Child process"

####################
# SPLAY
####################
DS=splay
echo $DS
M=10000
N=5000000

echo -- Multicore: m=$M n=$N d=$DS s=$SRC l=$LB 
for i in $(seq 1 48); do ./cmtf -m $M -n $N -d $DS -s $SRC -l $LB -k $i; done

echo -- Singlecore: m=$M n=$N d=$DS s=$SRC l=$LB cpu=$CPULIM
(for i in $(seq 1 48); do cpulimit -f -q -l $CPULIM -- ./cmtf -m $M -n $N -d $DS -s $SRC -l $LB -k $i; done) | grep -v "Child process"

