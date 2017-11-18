#!/bin/bash -
#
#

die(){
	echo "[DIE] $@"
	exit 1
}

go build tf.go

rm -f test.log
EXPECT=1400
for i in $(seq 200)
do
	./tf
	test $? -ne 1 && die "expect exitcode 1"
done

set -eu
trap 'rm -f test.log' EXIT

CNT=$(wc -l test.log | awk '{print $1}')
if test "X$CNT" != "X$EXPECT"
then
	echo "test failed, expected: $EXPECT but got $CNT"
	exit 1
else
	echo "check passed"
fi
