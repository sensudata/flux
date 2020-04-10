#!/bin/bash
#

#
# This directory contains some colm programs that can be used standalone, or in
# some cases, embedded into the Go codebase.
#
#  * flux: contains the primary flux grammar and a parsing driver
#
#  * influxql: contains timeboxed experiment to translate influxql -> flux. Very
#    far from complete.
#
#  * tableflux: contains a translator from TableFlux (proposed early 2020) to Flux.
#    This code can be included in the Go codebase.
#
# In the case of tableflux, by default, a stubbed interface will be included in
# the Go project. To make it functional you must first install colm, then
# configure and build in this directory with the go interface enabled. It will
# generate the appropriate Go/C files. See tableflux/call.*.in. These are
# rewritten by make.
#
# If you have trouble building, the first thing to do is get the latest colm on
# the master branch. New features may be added to colm to support the work
# here.
#

ORIG_PWD=$PWD

#
# 1. Install colm
#

DEVEL=$HOME/devel
INST=$HOME/pkgs

cd $DEVEL
git clone https://github.com/adrian-thurston/colm.git

cd colm
./autogen.sh
./configure --prefix=$INST/colm --disable-manual
make install

#
# 2. Configure in this directory, giving location of colm install, then make
#

cd $ORIG_PWD
./configure --with-colm=$INST/colm --enable-go-interface
make

#
# After this you should be able to build flux or influxdb (with go mod rewrite
# to this flux repos) and it should pick up the TableFlux implementation.
#
# To disable the implementation after enabling it, reconfigure without the
# --enable-go-interface flag and rebuild.
#


