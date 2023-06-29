#!/bin/bash
go test -run Basic > trace_Basic.txt
go test -run OneFailure > trace_OneFailure.txt
go test -run ManyFailures > trace_ManyFailures.txt