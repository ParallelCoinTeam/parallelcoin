#!/usr/bin/make -f

stroy:
	go install -v ./stroy/.
	stroy stroy
