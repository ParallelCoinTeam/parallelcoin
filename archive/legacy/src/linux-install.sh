#!/bin/bash
echo "Building and installing parallelcoind and parallelcoin-qt"
echo "Building openssh 1.0.1..."
wget -c https://github.com/openssl/openssl/archive/OpenSSL_1_0_1u.tar.gz
rm -rf openssl-OpenSSL_1_0_1u
tar zxvf OpenSSL_1_0_1u.tar.gz
cd openssl-OpenSSL_1_0_1u
./config
make -j$(nproc)
sudo make install
cd ..
echo "getting boost 1.58"
wget -c http://sourceforge.net/projects/boost/files/boost/1.58.0/boost_1_58_0.tar.gz
rm -rf boost_1_58_0
tar zxvf boost_1_58_0.tar.gz

echo "Building parallelcoind..."
make -j$(nproc) -f makefile.unix		# Headless bitcoin
echo "Building parallelcoin-qt..."
qmake
make -j$(nproc)
echo "Installing (you will need to enter your password for sudo)"
sudo cp parallelcoind parallelcoin-qt /usr/local/bin/
sudo cp parallelcoind.service /etc/systemd/system/
sed "s/####/`whoami`/g"
cp parallelcoin-qt.desktop $HOME/.local/share/applications/
cp qt/res/images/Wallet_Logo.png $HOME/.local/share/icons/parallelcoin.png
echo "cleaning up"
make cleanboo
make distclean
make -f makefile.unix clean
cd openssl-1.0.1u
make clean
echo "All done"
