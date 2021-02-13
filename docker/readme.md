# How to use these things

Herein are tools that work using Docker and the legacy parallelcoin repository at https://github.com/parallelcointeam/parallelcoin

The directories ending in .AppDir contain materials that when combined with the parallelcoin repository that is created by the docker in `legacy/` allows you to create an AppImage universal binary. It is built from ubuntu 14.04 so it should work on any 64 bit linux from 2014, and could be easily adapted to work with many other servers and wallets based on circa 2014-2016 bitcoin codebase.

To use the docker, you first need docker installed and the server running, then you can just `source init.sh` in the `legacy/` folder and then `halp` will show you all the short commands you can use and the long version they will invoke.

Then just copy those AppDir folders into a the `src/` subdirectory of the repository linked above, and for the Qt wallet you need to first run `linuxdeployqt-continuous-x86_64.AppImage` inside (it is a qmake dir, you can reinitialise like this using `qmake ../`) and then for the main currently just build, copy the binary in place and if necessary update the binaries in the `usr/lib` folder.

With these as a base it should be possible to create universal binaries that run everywhere on the same OS and ABI.