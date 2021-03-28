#!/usr/bin/make -f

builder:
	go install -v ./cmd/podbuild/.
	podbuild builder

ubuntu:
	./prereqs/ubuntu.sh

ubuntuglfw:
	./prereqs/ubuntu_glfw.sh

fedora:
	./prereqs/fedora.sh

