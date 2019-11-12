package utils

import (
	"github.com/therecipe/qt/webengine"
)

func GetBiosMessage(view *webengine.QWebEngineView, msg string) {
	view.Page().RunJavaScript(`var para = document.createElement("p");
var node = document.createTextNode("` + msg + `");
para.appendChild(node);
var element = document.getElementById("biosMessages");
element.appendChild(para);`)
}
