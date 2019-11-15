package main

/*
#cgo CFLAGS: -pipe -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -fno-plt -Wall -W -D_REENTRANT -fPIC -DQT_NO_DEBUG -DQT_QUICKCONTROLS2_LIB -DQT_QUICK_LIB -DQT_SVG_LIB -DQT_WIDGETS_LIB -DQT_GUI_LIB -DQT_QML_LIB -DQT_NETWORK_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -pipe -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -fno-plt -Wall -W -D_REENTRANT -fPIC -DQT_NO_DEBUG -DQT_QUICKCONTROLS2_LIB -DQT_QUICK_LIB -DQT_SVG_LIB -DQT_WIDGETS_LIB -DQT_GUI_LIB -DQT_QML_LIB -DQT_NETWORK_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -I. -isystem /usr/include/qt -isystem /usr/include/qt/QtQuickControls2 -isystem /usr/include/qt/QtQuick -isystem /usr/include/qt/QtSvg -isystem /usr/include/qt/QtWidgets -isystem /usr/include/qt/QtGui -isystem /usr/include/qt/QtQml -isystem /usr/include/qt/QtNetwork -isystem /usr/include/qt/QtCore -I. -I/usr/lib/qt/mkspecs/linux-g++
#cgo LDFLAGS: -Wl,-O1 -Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now -Wl,-rpath-link,/usr/lib
#cgo LDFLAGS: /usr/lib/libQt5QuickControls2.so /usr/lib/libQt5Quick.so /usr/lib/libQt5Svg.so /usr/lib/libQt5Widgets.so /usr/lib/libQt5Gui.so /usr/lib/libQt5Qml.so /usr/lib/libQt5Network.so /usr/lib/libQt5Core.so /usr/lib/libGL.so -lpthread
#cgo CFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
#cgo CXXFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
*/
import "C"
