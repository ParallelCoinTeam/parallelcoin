import QtQuick 2.3
import QtQuick.Window 2.3
import QtWebEngine 1.7

Window {
    visible: true
    width: 1024
    height: 750

    WebEngineView {
        anchors.fill: parent
        url: "index.html"
    }
}
