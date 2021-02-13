FROM ubuntu:trusty

RUN groupadd -r parallelcoin && useradd -r -m -g parallelcoin parallelcoin

RUN apt update
RUN apt -y dist-upgrade
RUN apt -y install build-essential
RUN apt -y install libssl-dev
RUN apt -y install libboost-all-dev
RUN apt install -y software-properties-common
RUN add-apt-repository -y ppa:bitcoin/bitcoin
RUN apt-get update
RUN apt -y install libdb4.8-dev
RUN apt -y install libdb4.8++-dev
RUN apt -y install libminiupnpc-dev
RUN apt -y install build-essential git
RUN apt -y install nano
RUN apt-get -y install qt4-qmake libqt4-dev build-essential \
  libboost-dev libboost-system-dev libboost-filesystem-dev \
  libboost-program-options-dev libboost-thread-dev \
  libssl-dev libdb++-dev libminiupnpc-dev 
VOLUME /data
WORKDIR /root/.parallelcoin
RUN chown parallelcoin /root/.parallelcoin
RUN cd /root/.parallelcoin \
  && git clone https://github.com/p9c/pod.git
RUN cd /root/.parallelcoin/pod/legacy/src \
  && make -f makefile.unix
RUN cd /root/.parallelcoin/pod/legacy/src \
  && mv parallelcoind /usr/bin/
#RUN cd /root/.parallelcoin/pod/legacy/src \
#  && qmake \
#  && make
EXPOSE 11048 11047 21048 21047
CMD ["tail", "-f", "/dev/null"]
