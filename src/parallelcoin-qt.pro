TEMPLATE = app
TARGET = parallelcoin-qt
macx:TARGET = "Parallelcoin-Qt"
VERSION = 0.9.0
INCLUDEPATH += json qt
QT += core gui network
greaterThan(QT_MAJOR_VERSION, 4): QT += widgets
DEFINES += BOOST_THREAD_USE_LIB BOOST_SPIRIT_THREADSAFE
CONFIG += no_include_pwd
CONFIG += thread

# for boost 1.37, add -mt to the boost libraries
# use: qmake BOOST_LIB_SUFFIX=-mt
# for boost thread win32 with _win32 sufix
# use: BOOST_THREAD_LIB_SUFFIX=_win32-...
# or when linking against a specific BerkelyDB version: BDB_LIB_SUFFIX=-4.8

# Dependency library locations can be customized with:
#    BOOST_INCLUDE_PATH, BOOST_LIB_PATH, BDB_INCLUDE_PATH,
#    BDB_LIB_PATH, OPENSSL_INCLUDE_PATH and OPENSSL_LIB_PATH respectively

# Because ancient version is required, set the path to /usr/local locations
OPENSSL_INCLUDE_PATH=/usr/local/ssl/include
OPENSSL_LIB_PATH=/usr/local/ssl/lib

win32 {
BOOST_LIB_SUFFIX=-mgw49-mt-s-1_57
BOOST_INCLUDE_PATH=C:/deps/boost_1_57_0
BOOST_LIB_PATH=C:/deps/boost_1_57_0/stage/lib
BDB_INCLUDE_PATH=C:/deps/db-4.8.30.NC/build_unix
BDB_LIB_PATH=C:/deps/db-4.8.30.NC/build_unix
OPENSSL_INCLUDE_PATH=C:/deps/openssl-1.0.1m/include
OPENSSL_LIB_PATH=C:/deps/openssl-1.0.1m
MINIUPNPC_INCLUDE_PATH=C:/deps/
MINIUPNPC_LIB_PATH=C:/deps/miniupnpc
QRENCODE_INCLUDE_PATH=C:/deps/qrencode-3.4.4
QRENCODE_LIB_PATH=C:/deps/qrencode-3.4.4/.libs
}

OBJECTS_DIR = build
MOC_DIR = build
UI_DIR = build

# use: qmake "RELEASE=1"
contains(RELEASE, 1) {
    # Mac: compile for maximum compatibility (10.5, 32-bit)
    macx:QMAKE_CXXFLAGS += -mmacosx-version-min=10.5 -arch i386 -isysroot /Developer/SDKs/MacOSX10.5.sdk
    macx:QMAKE_CFLAGS += -mmacosx-version-min=10.5 -arch i386 -isysroot /Developer/SDKs/MacOSX10.5.sdk
    macx:QMAKE_OBJECTIVE_CFLAGS += -mmacosx-version-min=10.5 -arch i386 -isysroot /Developer/SDKs/MacOSX10.5.sdk

    !win32:!macx {
        # Linux: static link and extra security (see: https://wiki.debian.org/Hardening)
        LIBS += -Wl,-Bstatic -Wl,-z,relro -Wl,-z,now
    }
}

# Make all those dumb warnings go away
QMAKE_CXXFLAGS *= -O2 -pthread -Wall -Wextra -Wformat -Wformat-security -Wno-unused-parameter \
    -Wno-deprecated-declarations -Wno-deprecated -Wno-unused-local-typedefs -Wno-class-memaccess \
    -Wno-unused-but-set-parameter -Wno-placement-new -Wno-maybe-uninitialized \
    -Wno-implicit-fallthrough \

!win32 {
    # for extra security against potential buffer overflows: enable GCCs Stack Smashing Protection
    QMAKE_CXXFLAGS *= -fstack-protector-all
    QMAKE_LFLAGS *= -fstack-protector-all
    # Exclude on Windows cross compile with MinGW 4.2.x, as it will result in a non-working executable!
    # This can be enabled for Windows, when we switch to MinGW >= 4.4.x.
}
# for extra security (see: https://wiki.debian.org/Hardening): this flag is GCC compiler-specific
QMAKE_CXXFLAGS *= -D_FORTIFY_SOURCE=2
# for extra security on Windows: enable ASLR and DEP via GCC linker flags
win32:QMAKE_LFLAGS *= -Wl,--dynamicbase -Wl,--nxcompat
# on Windows: enable GCC large address aware linker flag
win32:QMAKE_LFLAGS *= -Wl,--large-address-aware -static

# use: qmake "USE_QRCODE=1"
# libqrencode (http://fukuchi.org/works/qrencode/index.en.html) must be installed for support
contains(USE_QRCODE, 1) {
    message(Building with QRCode support)
    DEFINES += USE_QRCODE
    LIBS += -lqrencode
}

# use: qmake "USE_UPNP=1" ( enabled by default; default)
#  or: qmake "USE_UPNP=0" (disabled by default)
#  or: qmake "USE_UPNP=-" (not supported)
# miniupnpc (http://miniupnp.free.fr/files/) must be installed for support
contains(USE_UPNP, -) {
    message(Building without UPNP support)
} else {
    message(Building with UPNP support)
    count(USE_UPNP, 0) {
        USE_UPNP=1
    }
    DEFINES += USE_UPNP=$$USE_UPNP STATICLIB
    INCLUDEPATH += $$MINIUPNPC_INCLUDE_PATH
    LIBS += $$join(MINIUPNPC_LIB_PATH,,-L,) -lminiupnpc
    win32:LIBS += -liphlpapi
}

# use: qmake "USE_DBUS=1"
contains(USE_DBUS, 1) {
    message(Building with DBUS (Freedesktop notifications) support)
    DEFINES += USE_DBUS
    QT += dbus
}

# use: qmake "USE_IPV6=1" ( enabled by default; default)
#  or: qmake "USE_IPV6=0" (disabled by default)
#  or: qmake "USE_IPV6=-" (not supported)
contains(USE_IPV6, -) {
    message(Building without IPv6 support)
} else {
    count(USE_IPV6, 0) {
        USE_IPV6=1
    }
    DEFINES += USE_IPV6=$$USE_IPV6
}

contains(BITCOIN_NEED_QT_PLUGINS, 1) {
    DEFINES += BITCOIN_NEED_QT_PLUGINS
    QTPLUGIN += qcncodecs qjpcodecs qtwcodecs qkrcodecs qtaccessiblewidgets
}

INCLUDEPATH += leveldb/include leveldb/helpers
LIBS += $$PWD/leveldb/libleveldb.a $$PWD/leveldb/libmemenv.a
!win32 {
    # we use QMAKE_CXXFLAGS_RELEASE even without RELEASE=1 because we use RELEASE to indicate linking preferences not -O preferences
    genleveldb.commands = cd $$PWD/leveldb && CC=$$QMAKE_CC CXX=$$QMAKE_CXX $(MAKE) OPT=\"$$QMAKE_CXXFLAGS $$QMAKE_CXXFLAGS_RELEASE\" libleveldb.a libmemenv.a
} else {
    # make an educated guess about what the ranlib command is called
    isEmpty(QMAKE_RANLIB) {
        QMAKE_RANLIB = $$replace(QMAKE_STRIP, strip, ranlib)
    }
    LIBS += -lshlwapi
    #genleveldb.commands = cd $$PWD/leveldb && CC=$$QMAKE_CC CXX=$$QMAKE_CXX TARGET_OS=OS_WINDOWS_CROSSCOMPILE $(MAKE) OPT=\"$$QMAKE_CXXFLAGS $$QMAKE_CXXFLAGS_RELEASE\" libleveldb.a libmemenv.a && $$QMAKE_RANLIB $$PWD/leveldb/libleveldb.a && $$QMAKE_RANLIB $$PWD/leveldb/libmemenv.a
}
genleveldb.target = $$PWD/leveldb/libleveldb.a
genleveldb.depends = FORCE
PRE_TARGETDEPS += $$PWD/leveldb/libleveldb.a
QMAKE_EXTRA_TARGETS += genleveldb
# Gross ugly hack that depends on qmake internals, unfortunately there is no other way to do it.
QMAKE_CLEAN += $$PWD/leveldb/libleveldb.a; cd $$PWD/leveldb ; $(MAKE) clean

# regenerate build.h
!win32|contains(USE_BUILD_INFO, 1) {
    genbuild.depends = FORCE
    genbuild.commands = cd $$PWD; /bin/sh ../share/genbuild.sh $$OUT_PWD/build/build.h
    genbuild.target = $$OUT_PWD/build/build.h
    PRE_TARGETDEPS += $$OUT_PWD/build/build.h
    QMAKE_EXTRA_TARGETS += genbuild
    DEFINES += HAVE_BUILD_INFO
}

QMAKE_CXXFLAGS_WARN_ON = -fdiagnostics-show-option -Wall -Wextra -Wformat -Wformat-security -Wno-unused-parameter -Wstack-protector

# Input
DEPENDPATH += src json qt
HEADERS += qt/bitcoingui.h \
    qt/transactiontablemodel.h \
    qt/addresstablemodel.h \
    qt/optionsdialog.h \
    qt/sendcoinsdialog.h \
    qt/addressbookpage.h \
    qt/signverifymessagedialog.h \
    qt/aboutdialog.h \
    qt/editaddressdialog.h \
    qt/bitcoinaddressvalidator.h \
    alert.h \
    addrman.h \
    base58.h \
    bignum.h \
    chainparams.h \
    checkpoints.h \
    compat.h \
    sync.h \
    util.h \
    hash.h \
    uint256.h \
    serialize.h \
    core.h \
    main.h \
    net.h \
    key.h \
    db.h \
    walletdb.h \
    script.h \
    init.h \
    bloom.h \
    mruset.h \
    checkqueue.h \
    json/json_spirit_writer_template.h \
    json/json_spirit_writer.h \
    json/json_spirit_value.h \
    json/json_spirit_utils.h \
    json/json_spirit_stream_reader.h \
    json/json_spirit_reader_template.h \
    json/json_spirit_reader.h \
    json/json_spirit_error_position.h \
    json/json_spirit.h \
    qt/clientmodel.h \
    qt/guiutil.h \
    qt/transactionrecord.h \
    qt/guiconstants.h \
    qt/optionsmodel.h \
    qt/monitoreddatamapper.h \
    qt/transactiondesc.h \
    qt/transactiondescdialog.h \
    qt/bitcoinamountfield.h \
    wallet.h \
    keystore.h \
    qt/transactionfilterproxy.h \
    qt/transactionview.h \
    qt/walletmodel.h \
    qt/walletview.h \
    qt/walletstack.h \
    qt/walletframe.h \
    bitcoinrpc.h \
    qt/overviewpage.h \
    qt/csvmodelwriter.h \
    crypter.h \
    qt/sendcoinsentry.h \
    qt/qvalidatedlineedit.h \
    qt/bitcoinunits.h \
    qt/qvaluecombobox.h \
    qt/askpassphrasedialog.h \
    protocol.h \
    qt/notificator.h \
    qt/paymentserver.h \
    allocators.h \
    ui_interface.h \
    qt/rpcconsole.h \
    version.h \
    netbase.h \
    clientversion.h \
    txdb.h \
    leveldb.h \
    threadsafety.h \
    limitedmap.h \
    qt/splashscreen.h \
    qt/intro.h \
    scrypt.h

SOURCES += qt/bitcoin.cpp \
    qt/bitcoingui.cpp \
    qt/transactiontablemodel.cpp \
    qt/addresstablemodel.cpp \
    qt/optionsdialog.cpp \
    qt/sendcoinsdialog.cpp \
    qt/addressbookpage.cpp \
    qt/signverifymessagedialog.cpp \
    qt/aboutdialog.cpp \
    qt/editaddressdialog.cpp \
    qt/bitcoinaddressvalidator.cpp \
    alert.cpp \
    chainparams.cpp \
    version.cpp \
    sync.cpp \
    util.cpp \
    hash.cpp \
    netbase.cpp \
    key.cpp \
    script.cpp \
    core.cpp \
    main.cpp \
    init.cpp \
    net.cpp \
    bloom.cpp \
    checkpoints.cpp \
    addrman.cpp \
    db.cpp \
    walletdb.cpp \
    qt/clientmodel.cpp \
    qt/guiutil.cpp \
    qt/transactionrecord.cpp \
    qt/optionsmodel.cpp \
    qt/monitoreddatamapper.cpp \
    qt/transactiondesc.cpp \
    qt/transactiondescdialog.cpp \
    qt/bitcoinstrings.cpp \
    qt/bitcoinamountfield.cpp \
    wallet.cpp \
    keystore.cpp \
    qt/transactionfilterproxy.cpp \
    qt/transactionview.cpp \
    qt/walletmodel.cpp \
    qt/walletview.cpp \
    qt/walletstack.cpp \
    qt/walletframe.cpp \
    bitcoinrpc.cpp \
    rpcdump.cpp \
    rpcnet.cpp \
    rpcmining.cpp \
    rpcwallet.cpp \
    rpcblockchain.cpp \
    rpcrawtransaction.cpp \
    qt/overviewpage.cpp \
    qt/csvmodelwriter.cpp \
    crypter.cpp \
    qt/sendcoinsentry.cpp \
    qt/qvalidatedlineedit.cpp \
    qt/bitcoinunits.cpp \
    qt/qvaluecombobox.cpp \
    qt/askpassphrasedialog.cpp \
    protocol.cpp \
    qt/notificator.cpp \
    qt/paymentserver.cpp \
    qt/rpcconsole.cpp \
    noui.cpp \
    leveldb.cpp \
    txdb.cpp \
    qt/splashscreen.cpp \
    qt/intro.cpp \
    scrypt.cpp

RESOURCES += qt/bitcoin.qrc

FORMS += qt/forms/sendcoinsdialog.ui \
    qt/forms/addressbookpage.ui \
    qt/forms/signverifymessagedialog.ui \
    qt/forms/aboutdialog.ui \
    qt/forms/editaddressdialog.ui \
    qt/forms/transactiondescdialog.ui \
    qt/forms/overviewpage.ui \
    qt/forms/sendcoinsentry.ui \
    qt/forms/askpassphrasedialog.ui \
    qt/forms/rpcconsole.ui \
    qt/forms/optionsdialog.ui \
    qt/forms/intro.ui

contains(USE_QRCODE, 1) {
HEADERS += qt/qrcodedialog.h
SOURCES += qt/qrcodedialog.cpp
FORMS += qt/forms/qrcodedialog.ui
}

contains(BITCOIN_QT_TEST, 1) {
SOURCES += qt/test/test_main.cpp \
    qt/test/uritests.cpp
HEADERS += qt/test/uritests.h
DEPENDPATH += qt/test
QT += testlib
TARGET = parallelcoin-qt_test
DEFINES += BITCOIN_QT_TEST
  macx: CONFIG -= app_bundle
}

# Todo: Remove this line when switching to Qt5, as that option was removed
CODECFORTR = UTF-8

# for lrelease/lupdate
# also add new translations to qt/bitcoin.qrc under translations/
TRANSLATIONS = $$files(qt/locale/bitcoin_*.ts)

isEmpty(QMAKE_LRELEASE) {
    win32:QMAKE_LRELEASE = $$[QT_INSTALL_BINS]\\lrelease.exe
    else:QMAKE_LRELEASE = $$[QT_INSTALL_BINS]/lrelease
}
isEmpty(QM_DIR):QM_DIR = $$PWD/qt/locale
# automatically build translations, so they can be included in resource file
TSQM.name = lrelease ${QMAKE_FILE_IN}
TSQM.input = TRANSLATIONS
TSQM.output = $$QM_DIR/${QMAKE_FILE_BASE}.qm
TSQM.commands = $$QMAKE_LRELEASE ${QMAKE_FILE_IN} -qm ${QMAKE_FILE_OUT}
TSQM.CONFIG = no_link
QMAKE_EXTRA_COMPILERS += TSQM

# "Other files" to show in Qt Creator
OTHER_FILES += README.md \
    doc/*.txt \
    doc/*.md \
    bitcoind.cpp \
    qt/res/bitcoin-qt.rc \
    test/*.cpp \
    test/*.h \
    qt/test/*.cpp \
    qt/test/*.h

# platform specific defaults, if not overridden on command line
isEmpty(BOOST_LIB_SUFFIX) {
    macx:BOOST_LIB_SUFFIX = -mt
    win32:BOOST_LIB_SUFFIX = -mgw44-mt-s-1_50
}

isEmpty(BOOST_THREAD_LIB_SUFFIX) {
    BOOST_THREAD_LIB_SUFFIX = $$BOOST_LIB_SUFFIX
}

isEmpty(BDB_LIB_PATH) {
    macx:BDB_LIB_PATH = /opt/local/lib/db48
}

isEmpty(BDB_LIB_SUFFIX) {
    macx:BDB_LIB_SUFFIX = -4.8
}

isEmpty(BDB_INCLUDE_PATH) {
    macx:BDB_INCLUDE_PATH = /opt/local/include/db48
}

isEmpty(BOOST_LIB_PATH) {
    macx:BOOST_LIB_PATH = /opt/local/lib
}

isEmpty(BOOST_INCLUDE_PATH) {
    macx:BOOST_INCLUDE_PATH = /opt/local/include
}

win32:DEFINES += WIN32
win32:RC_FILE = qt/res/bitcoin-qt.rc

win32:!contains(MINGW_THREAD_BUGFIX, 0) {
    # At least qmake's win32-g++-cross profile is missing the -lmingwthrd
    # thread-safety flag. GCC has -mthreads to enable this, but it doesn't
    # work with static linking. -lmingwthrd must come BEFORE -lmingw, so
    # it is prepended to QMAKE_LIBS_QT_ENTRY.
    # It can be turned off with MINGW_THREAD_BUGFIX=0, just in case it causes
    # any problems on some untested qmake profile now or in the future.
    DEFINES += _MT
    QMAKE_LIBS_QT_ENTRY = -lmingwthrd $$QMAKE_LIBS_QT_ENTRY
}

!win32:!macx {
    DEFINES += LINUX
    LIBS += -lrt
    # _FILE_OFFSET_BITS=64 lets 32-bit fopen transparently support large files.
    DEFINES += _FILE_OFFSET_BITS=64
}

macx:HEADERS += qt/macdockiconhandler.h qt/macnotificationhandler.h
macx:OBJECTIVE_SOURCES += qt/macdockiconhandler.mm qt/macnotificationhandler.mm
macx:LIBS += -framework Foundation -framework ApplicationServices -framework AppKit -framework CoreServices
macx:DEFINES += MAC_OSX MSG_NOSIGNAL=0
macx:ICON = qt/res/icons/bitcoin.icns
macx:QMAKE_CFLAGS_THREAD += -pthread
macx:QMAKE_LFLAGS_THREAD += -pthread
macx:QMAKE_CXXFLAGS_THREAD += -pthread
macx:QMAKE_INFO_PLIST = ../share/qt/Info.plist

# Set libraries and includes at end, to use platform-defined defaults if not overridden
INCLUDEPATH += $$BOOST_INCLUDE_PATH $$BDB_INCLUDE_PATH $$OPENSSL_INCLUDE_PATH $$QRENCODE_INCLUDE_PATH
LIBS += $$join(BOOST_LIB_PATH,,-L,) $$join(BDB_LIB_PATH,,-L,) $$join(OPENSSL_LIB_PATH,,-L,) $$join(QRENCODE_LIB_PATH,,-L,)
LIBS += -lssl -lcrypto -ldb_cxx$$BDB_LIB_SUFFIX
# -lgdi32 has to happen after -lcrypto (see  #681)
win32:LIBS += -lws2_32 -lshlwapi -lmswsock -lole32 -loleaut32 -luuid -lgdi32
LIBS += -lboost_system$$BOOST_LIB_SUFFIX -lboost_filesystem$$BOOST_LIB_SUFFIX -lboost_program_options$$BOOST_LIB_SUFFIX -lboost_thread$$BOOST_THREAD_LIB_SUFFIX -ldl
win32:LIBS += -lboost_chrono$$BOOST_LIB_SUFFIX
macx:LIBS += -lboost_chrono$$BOOST_LIB_SUFFIX

contains(RELEASE, 1) {
    !win32:!macx {
        # Linux: turn dynamic linking back on for c/c++ runtime libraries
        LIBS += -Wl,-Bdynamic
    }
}

system($$QMAKE_LRELEASE -silent $$TRANSLATIONS)
