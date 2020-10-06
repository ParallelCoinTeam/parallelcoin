// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"image/color"

	"github.com/p9c/pod/pkg/gui/ico"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func (th *Theme) NewIcons() (i map[string]*_icon) {
	i = make(map[string]*_icon)
	i["Checked"] = mustIcon(th.Icon(icons.ToggleCheckBox))
	i["Unchecked"] = mustIcon(th.Icon(icons.ToggleCheckBoxOutlineBlank))
	i["RadioChecked"] = mustIcon(th.Icon(icons.ToggleRadioButtonChecked))
	i["RadioUnchecked"] = mustIcon(th.Icon(icons.ToggleRadioButtonUnchecked))
	i["iconCancel"] = mustIcon(th.Icon(icons.NavigationCancel))
	i["iconOK"] = mustIcon(th.Icon(icons.NavigationCheck))
	i["iconClose"] = mustIcon(th.Icon(icons.NavigationClose))
	i["foldIn"] = mustIcon(th.Icon(icons.ContentRemove))
	i["minimize"] = mustIcon(th.Icon(icons.NavigationExpandMore))
	i["zoom"] = mustIcon(th.Icon(icons.NavigationExpandLess))
	i["logo"] = mustIcon(th.Icon(ico.ParallelCoin))
	i["overviewIcon"] = mustIcon(th.Icon(icons.ActionHome))
	i["sendIcon"] = mustIcon(th.Icon(icons.ActionStarRate))
	i["receiveIcon"] = mustIcon(th.Icon(icons.NavigationArrowDropDown))
	i["addressBookIcon"] = mustIcon(th.Icon(icons.ActionBook))
	i["historyIcon"] = mustIcon(th.Icon(icons.ActionHistory))
	i["closeIcon"] = mustIcon(th.Icon(icons.NavigationClose))
	i["settingsIcon"] = mustIcon(th.Icon(icons.ActionSettings))
	i["blocksIcon"] = mustIcon(th.Icon(icons.ActionExplore))
	i["networkIcon"] = mustIcon(th.Icon(icons.ActionFingerprint))
	i["traceIcon"] = mustIcon(th.Icon(icons.ActionTrackChanges))
	// i["consoleIcon"] = mustIcon(Icon(icons.ActionInput))
	i["helpIcon"] = mustIcon(th.Icon(icons.NavigationArrowDropDown))
	i["counterPlusIcon"] = mustIcon(th.Icon(icons.ImageExposurePlus1))
	i["counterMinusIcon"] = mustIcon(th.Icon(icons.ImageExposureNeg1))
	i["CommunicationImportExport"] = mustIcon(th.Icon(icons.CommunicationImportExport))
	i["NotificationNetworkCheck"] = mustIcon(th.Icon(icons.NotificationNetworkCheck))
	i["NotificationSync"] = mustIcon(th.Icon(icons.NotificationSync))
	i["NotificationSyncDisabled"] = mustIcon(th.Icon(icons.NotificationSyncDisabled))
	i["NotificationSyncProblem"] = mustIcon(th.Icon(icons.NotificationSyncProblem))
	i["NotificationVPNLock"] = mustIcon(th.Icon(icons.NotificationVPNLock))
	i["network"] = mustIcon(th.Icon(icons.NotificationWiFi))
	i["MapsLayers"] = mustIcon(th.Icon(icons.MapsLayers))
	i["MapsLayersClear"] = mustIcon(th.Icon(icons.MapsLayersClear))
	i["ImageTimer"] = mustIcon(th.Icon(icons.ImageTimer))
	i["ImageRemoveRedEye"] = mustIcon(th.Icon(icons.ImageRemoveRedEye))
	i["DeviceSignalCellular0Bar"] = mustIcon(th.Icon(icons.DeviceSignalCellular0Bar))
	i["DeviceWidgets"] = mustIcon(th.Icon(icons.DeviceWidgets))
	i["ActionTimeline"] = mustIcon(th.Icon(icons.ActionTimeline))
	i["HardwareWatch"] = mustIcon(th.Icon(icons.HardwareWatch))
	i["consoleIcon"] = mustIcon(th.Icon(icons.HardwareKeyboardHide))
	i["DeviceSignalCellular0Bar"] = mustIcon(th.Icon(icons.DeviceSignalCellular0Bar))
	i["HardwareWatch"] = mustIcon(th.Icon(icons.HardwareWatch))
	i["EditorMonetizationOn"] = mustIcon(th.Icon(icons.EditorMonetizationOn))
	i["Run"] = mustIcon(th.Icon(icons.AVPlayArrow))
	i["Stop"] = mustIcon(th.Icon(icons.AVStop))
	i["Pause"] = mustIcon(th.Icon(icons.AVPause))
	i["Kill"] = mustIcon(th.Icon(icons.NavigationCancel))
	i["Restart"] = mustIcon(th.Icon(icons.NavigationRefresh))
	i["Grab"] = mustIcon(th.Icon(icons.NavigationMenu))
	i["Up"] = mustIcon(th.Icon(icons.NavigationArrowDropUp))
	i["Down"] = mustIcon(th.Icon(icons.NavigationArrowDropDown))
	i["iconGrab"] = mustIcon(th.Icon(icons.NavigationMenu))
	i["iconUp"] = mustIcon(th.Icon(icons.NavigationArrowDropUp))
	i["iconDown"] = mustIcon(th.Icon(icons.NavigationArrowDropDown))
	i["Copy"] = mustIcon(th.Icon(icons.ContentContentCopy))
	i["Paste"] = mustIcon(th.Icon(icons.ContentContentPaste))
	i["Sidebar"] = mustIcon(th.Icon(icons.ActionChromeReaderMode))
	i["Filter"] = mustIcon(th.Icon(icons.ContentFilterList))
	i["FilterAll"] = mustIcon(th.Icon(icons.ActionDoneAll))
	i["FilterNone"] = mustIcon(th.Icon(icons.ContentBlock))
	i["Build"] = mustIcon(th.Icon(icons.ActionBuild))
	i["Folded"] = mustIcon(th.Icon(icons.NavigationChevronRight))
	i["Unfolded"] = mustIcon(th.Icon(icons.NavigationExpandMore))
	i["HideAll"] = mustIcon(th.Icon(icons.NavigationUnfoldLess))
	i["ShowAll"] = mustIcon(th.Icon(icons.NavigationUnfoldMore))
	i["HideItem"] = mustIcon(th.Icon(icons.ActionVisibilityOff))
	i["ShowItem"] = mustIcon(th.Icon(icons.ActionVisibility))
	i["TRC"] = mustIcon(th.Icon(icons.ActionSearch))
	i["DBG"] = mustIcon(th.Icon(icons.ActionBugReport))
	i["INF"] = mustIcon(th.Icon(icons.ActionInfo))
	i["WRN"] = mustIcon(th.Icon(icons.ActionHelp))
	i["CHK"] = mustIcon(th.Icon(icons.AlertWarning))
	i["ERR"] = mustIcon(th.Icon(icons.AlertError))
	i["FTL"] = mustIcon(th.Icon(icons.ImageFlashOn))
	i["Delete"] = mustIcon(th.Icon(icons.ActionDelete))
	i["Send"] = mustIcon(th.Icon(icons.ContentSend))
	return i
}

func mustIcon(ic *_icon, err error) *_icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
//
// func fill(gtx *layout.Context, col color.RGBA) {
// 	cs := gtx.Constraints
// 	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
// 	dr := f32.Rectangle{
// 		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
// 	}
// 	paint.ColorOp{Color: col}.Add(gtx.Ops)
// 	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
// 	gtx.Dimensions = layout.Dimensions{Size: d}
// }
