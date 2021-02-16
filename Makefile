#!/usr/bin/make -f

stroy:
	go install -v ./stroy/.
	stroy stroy

ubuntu:
	sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev
	go install -v ./stroy/.
	stroy stroy

