#!/bin/bash
yes|sudo pacman -Syyu
yes|sudo pacman -S yay
yes|yay -Syyu
yes|yay --needed -S bash wget curl git make base-devel xclip
