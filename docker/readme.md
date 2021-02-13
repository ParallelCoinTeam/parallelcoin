# How to use these things

Herein are tools that work using Docker and the legacy parallelcoin repository at 
https://github.com/parallelcointeam/parallelcoin

## The Docker

To use the docker, you first need docker installed and the server running, then you can just `source init.sh` in the 
`legacy/` folder and then `halp` will show you all the short commands you can use and the long version they will invoke.

The following sequence will get you into a shell inside the docker where you can run and use the RPC:

```
loki@monolith:~/src/github.com/p9c/pod/docker/legacy$ .build
Sending build context to Docker daemon  197.1MB
Step 1/24 : FROM ubuntu:trusty
 ---> df043b4f0cf1
Step 2/24 : RUN groupadd -r parallelcoin && useradd -r -m -g parallelcoin parallelcoin
 ---> Using cache
 ---> 7c9206c3c1f3
Step 3/24 : RUN apt update
 ---> Using cache
 ---> 92128950c8fb
Step 4/24 : RUN apt -y dist-upgrade
 ---> Using cache
 ---> 65daa82b66aa
Step 5/24 : RUN apt -y install build-essential
 ---> Using cache
 ---> 799f4a53adc5
Step 6/24 : RUN apt -y install libssl-dev
 ---> Using cache
 ---> 09ca5955f00c
Step 7/24 : RUN apt -y install libboost-all-dev
 ---> Using cache
 ---> c693db26c119
Step 8/24 : RUN apt install -y software-properties-common
 ---> Using cache
 ---> 22ad44ae1d34
Step 9/24 : RUN add-apt-repository -y ppa:bitcoin/bitcoin
 ---> Using cache
 ---> c0b9fbd38c66
Step 10/24 : RUN apt-get update
 ---> Using cache
 ---> 86999a2a04cd
Step 11/24 : RUN apt -y install libdb4.8-dev
 ---> Using cache
 ---> 93bd7da7560e
Step 12/24 : RUN apt -y install libdb4.8++-dev
 ---> Using cache
 ---> b8decbf07bf3
Step 13/24 : RUN apt -y install libminiupnpc-dev
 ---> Using cache
 ---> d4fbbeb6ca0f
Step 14/24 : RUN apt -y install build-essential git
 ---> Using cache
 ---> df0a46a733dd
Step 15/24 : RUN apt -y install nano
 ---> Using cache
 ---> bb5ef4676aca
Step 16/24 : RUN apt-get -y install qt4-qmake libqt4-dev build-essential   libboost-dev libboost-system-dev libboost-filesystem-dev   libboost-program-options-dev libboost-thread-dev   libssl-dev libdb++-dev libminiupnpc-dev
 ---> Using cache
 ---> 4e71019e1a69
Step 17/24 : VOLUME /data
 ---> Using cache
 ---> 0e58aa612770
Step 18/24 : WORKDIR /root/.parallelcoin
 ---> Using cache
 ---> 57c46001f962
Step 19/24 : RUN chown parallelcoin /root/.parallelcoin
 ---> Using cache
 ---> c9d2af9bd1e6
Step 20/24 : RUN cd /root/.parallelcoin   && git clone https://github.com/p9c/pod.git
 ---> Using cache
 ---> 84ae497c02ee
Step 21/24 : RUN cd /root/.parallelcoin/pod/legacy/src   && make -f makefile.unix
 ---> Using cache
 ---> efd8d7d05007
Step 22/24 : RUN cd /root/.parallelcoin/pod/legacy/src   && mv parallelcoind /usr/bin/
 ---> Using cache
 ---> 51334e4e48bf
Step 23/24 : EXPOSE 11048 11047 21048 21047
 ---> Using cache
 ---> 6c71911c24d9
Step 24/24 : CMD ["tail", "-f", "/dev/null"]
 ---> Using cache
 ---> 5a172e084ca0
Successfully built 5a172e084ca0
Successfully tagged docker-parallelcoind:latest
loki@monolith:~/src/github.com/p9c/pod/docker/legacy$ .run
WARNING: Published ports are discarded when using host network mode
fe225c0c808815bee68d327ff627eb08b9c131bf3dbc9e5bf1caabcc97a3985f
loki@monolith:~/src/github.com/p9c/pod/docker/legacy$ .enter
root@monolith:~/.parallelcoin# parallelcoind
Parallelcoin server starting
root@monolith:~/.parallelcoin# parallelcoind getinfo
{
    "version" : 1020000,
    "protocolversion" : 80000,
    "walletversion" : 60000,
    "balance" : 0.00000000,
    "blocks" : 250346,
    "timeoffset" : 0,
    "connections" : 0,
    "proxy" : "",
    "pow_algo_id" : 0,
    "pow_algo" : "sha256d",
    "difficulty" : 30768871.39866794,
    "difficulty_sha256d" : 30768871.39866794,
    "difficulty_scrypt" : 1395.89346375,
    "testnet" : false,
    "keypoololdest" : 1613221503,
    "keypoolsize" : 101,
    "paytxfee" : 0.00000000,
    "errors" : ""
}
root@monolith:~/.parallelcoin# 

```

## Building an AppImage

The directories ending in .AppDir contain materials that when combined with the parallelcoin repository that is created 
by the docker in `legacy/` allows you to create an AppImage universal binary. It is built from ubuntu 14.04 so it should 
work on any 64 bit linux from 2014, and could be easily adapted to work with many other servers and wallets based on 
circa 2014-2016 bitcoin codebase.

Then just copy those AppDir folders into a the `src/` subdirectory of the repository linked above, and for the Qt wallet
you need to first run `linuxdeployqt-continuous-x86_64.AppImage` inside (it is a qmake dir, you can reinitialise like 
this using `qmake ../`) and then for the main currently just build, copy the binary in place and if necessary update the 
binaries in the `usr/lib` folder.

With these as a base it should be possible to create universal binaries that run everywhere on the same OS and ABI.