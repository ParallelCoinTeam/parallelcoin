package cl

import "github.com/mitchellh/colorstring"

func ftlTag(color bool) string {
	tag := "FTL"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[red]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "
}
func errTag(color bool) string {
	tag := "ERR"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[yellow]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "
}
func wrnTag(color bool) string {
	tag := "WRN"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[green]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "
}
func infTag(color bool) string {
	tag := "INF"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[cyan]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "
}
func dbgTag(color bool) string {
	tag := "DBG"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[blue]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "
}
func trcTag(color bool) string {
	tag := "TRC"
	if color {
		pre := ""  // colorstring.Color("[light_gray][[dark_gray]")
		post := "" // colorstring.Color("[light_gray]][dark_gray]")
		tag = pre + colorstring.Color("[magenta]"+colorstring.Color("[bold]"+tag)) + post
	} else {
		tag = "[" + tag + "]"
	}
	return " " + tag + " "
}
