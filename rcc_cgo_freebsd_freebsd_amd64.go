package main

/*
#cgo CFLAGS: -pipe -O2 -Wall -W -pthread -fPIC -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -pipe -O2 -Wall -W -pthread -fPIC -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -I../../p9c -I. -I/usr/local/include/qt5 -I/usr/local/include/qt5/QtGui -I/usr/local/include/qt5/QtCore -I. -I/usr/local/include/libdrm -I/usr/local/include -I/usr/local/lib/qt5/mkspecs/freebsd-clang
#cgo LDFLAGS: -pthread -Wl,-rpath,/usr/local/lib/qt5
#cgo LDFLAGS:  -L/usr/local/lib /usr/local/lib/qt5/libQt5Gui.so /usr/local/lib/qt5/libQt5Core.so /usr/local/lib/libGL.so
#cgo CFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
#cgo CXXFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
*/
import "C"
