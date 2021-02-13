
FROM jonathonf/manjaro
RUN mkdir -p /root/.gnupg /root/.gnupg/dirmngr-cache.d 
RUN touch /root/.gnupg/dirmngr_ldapservers.conf 
RUN pacman -Sy 
RUN gpg -k 
RUN pacman-key --init 
RUN pacman-key --populate archlinux 
RUN pacman-key --populate manjaro 
RUN dirmngr < /dev/null 
RUN pacman -S --noconfirm pacman-mirrors 
RUN pacman -S --noconfirm libusb
# RUN pacman -S --noconfirm --force ca-certificates 
RUN pacman -Su --noconfirm 
RUN pacman -Sc --noconfirm
RUN pacman -S --noconfirm pacman-contrib
RUN rankmirrors -f 15  
RUN pacman -S --noconfirm base-devel 
RUN pacman -S --noconfirm wget 
RUN pacman -S --noconfirm git
RUN pacman -S --noconfirm go
RUN wget https://github.com/Jguer/yay/archive/v8.1157.0.tar.gz; tar xvf v8.1157.0.tar.gz;cd yay-8.1157.0; make; make install;cd; ln -s /usr/bin/yay /usr/bin/pacaur
RUN set -o pipefail -e; AUR_USER=user; useradd -m $AUR_USER; echo "${AUR_USER}:" | chpasswd -e ; \
    pkgs=$(pacman -S base-devel --print-format '%n ');pkgs=${pkgs//systemd/};pkgs=${pkgs//$'\n'/}; \
    pacman -S --needed --noprogressbar --noconfirm $pkgs vi;\
    echo "$AUR_USER      ALL = NOPASSWD: ALL" >> /etc/sudoers;\
    sed -i 's,#MAKEFLAGS="-j2",MAKEFLAGS="-j$(nproc)",g' /etc/makepkg.conf;\
    sed -i "s,PKGEXT='.pkg.tar.xz',PKGEXT='.pkg.tar',g" /etc/makepkg.conf;
RUN su user -c 'yay -S --noconfirm base-devel boost boost-libs miniupnpc db4.8 openssl-1.0'
RUN git clone https://github.com/parallelcointeam/parallelcoin.git
RUN cd parallelcoin/src; \
    echo "Building parallelcoind...";\
    make -j$(nproc) -f makefile.unix;\
    sudo cp parallelcoind /usr/local/bin/;\
    echo "cleaning up";\
    make -f makefile.unix clean;\
    echo "All done"
RUN su user -c 'yay -S --noconfirm nano net-tools'
EXPOSE 11047 
EXPOSE 11048
EXPOSE 21047
EXPOSE 21048
WORKDIR /root
VOLUME /data
RUN mkdir /root/.parallelcoin
# RUN ls -la
CMD /usr/local/bin/parallelcoind -txindex -printtoconsole -debug -debugnet -rpcuser=user -rpcpassword=pa55word -connect=127.0.0.1:11047 -connect=seed1.parallelcoin.io -bind=127.0.0.1 -port=11147 -rpcport=11148
#CMD tail -f /dev/null
