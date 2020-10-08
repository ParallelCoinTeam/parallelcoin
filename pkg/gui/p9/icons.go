// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"image/color"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/ico"
)

func (th *Theme) NewIcons() (i map[string]*Ico) {
	i = make(map[string]*Ico)
	i["Checked"] = th.Icon().Src(icons.ToggleCheckBox)
	i["Unchecked"] = th.Icon().Src(icons.ToggleCheckBoxOutlineBlank)
	i["RadioChecked"] = th.Icon().Src(icons.ToggleRadioButtonChecked)
	i["RadioUnchecked"] = th.Icon().Src(icons.ToggleRadioButtonUnchecked)
	i["iconCancel"] = th.Icon().Src(icons.NavigationCancel)
	i["iconOK"] = th.Icon().Src(icons.NavigationCheck)
	i["iconClose"] = th.Icon().Src(icons.NavigationClose)
	i["foldIn"] = th.Icon().Src(icons.ContentRemove)
	i["minimize"] = th.Icon().Src(icons.NavigationExpandMore)
	i["zoom"] = th.Icon().Src(icons.NavigationExpandLess)
	i["logo"] = th.Icon().Src(ico.ParallelCoin)
	i["overviewIcon"] = th.Icon().Src(icons.ActionHome)
	i["sendIcon"] = th.Icon().Src(icons.ActionStarRate)
	i["receiveIcon"] = th.Icon().Src(icons.NavigationArrowDropDown)
	i["addressBookIcon"] = th.Icon().Src(icons.ActionBook)
	i["historyIcon"] = th.Icon().Src(icons.ActionHistory)
	i["closeIcon"] = th.Icon().Src(icons.NavigationClose)
	i["settingsIcon"] = th.Icon().Src(icons.ActionSettings)
	i["blocksIcon"] = th.Icon().Src(icons.ActionExplore)
	i["networkIcon"] = th.Icon().Src(icons.ActionFingerprint)
	i["traceIcon"] = th.Icon().Src(icons.ActionTrackChanges)
	// i["consoleIcon"] = Icon().Src(icons.ActionInput)
	i["helpIcon"] = th.Icon().Src(icons.NavigationArrowDropDown)
	i["counterPlusIcon"] = th.Icon().Src(icons.ImageExposurePlus1)
	i["counterMinusIcon"] = th.Icon().Src(icons.ImageExposureNeg1)
	i["CommunicationImportExport"] = th.Icon().Src(icons.CommunicationImportExport)
	i["NotificationNetworkCheck"] = th.Icon().Src(icons.NotificationNetworkCheck)
	i["NotificationSync"] = th.Icon().Src(icons.NotificationSync)
	i["NotificationSyncDisabled"] = th.Icon().Src(icons.NotificationSyncDisabled)
	i["NotificationSyncProblem"] = th.Icon().Src(icons.NotificationSyncProblem)
	i["NotificationVPNLock"] = th.Icon().Src(icons.NotificationVPNLock)
	i["network"] = th.Icon().Src(icons.NotificationWiFi)
	i["MapsLayers"] = th.Icon().Src(icons.MapsLayers)
	i["MapsLayersClear"] = th.Icon().Src(icons.MapsLayersClear)
	i["ImageTimer"] = th.Icon().Src(icons.ImageTimer)
	i["ImageRemoveRedEye"] = th.Icon().Src(icons.ImageRemoveRedEye)
	i["DeviceSignalCellular0Bar"] = th.Icon().Src(icons.DeviceSignalCellular0Bar)
	i["DeviceWidgets"] = th.Icon().Src(icons.DeviceWidgets)
	i["ActionTimeline"] = th.Icon().Src(icons.ActionTimeline)
	i["HardwareWatch"] = th.Icon().Src(icons.HardwareWatch)
	i["consoleIcon"] = th.Icon().Src(icons.HardwareKeyboardHide)
	i["DeviceSignalCellular0Bar"] = th.Icon().Src(icons.DeviceSignalCellular0Bar)
	i["HardwareWatch"] = th.Icon().Src(icons.HardwareWatch)
	i["EditorMonetizationOn"] = th.Icon().Src(icons.EditorMonetizationOn)
	i["Run"] = th.Icon().Src(icons.AVPlayArrow)
	i["Stop"] = th.Icon().Src(icons.AVStop)
	i["Pause"] = th.Icon().Src(icons.AVPause)
	i["Kill"] = th.Icon().Src(icons.NavigationCancel)
	i["Restart"] = th.Icon().Src(icons.NavigationRefresh)
	i["Grab"] = th.Icon().Src(icons.NavigationMenu)
	i["Up"] = th.Icon().Src(icons.NavigationArrowDropUp)
	i["Down"] = th.Icon().Src(icons.NavigationArrowDropDown)
	i["iconGrab"] = th.Icon().Src(icons.NavigationMenu)
	i["iconUp"] = th.Icon().Src(icons.NavigationArrowDropUp)
	i["iconDown"] = th.Icon().Src(icons.NavigationArrowDropDown)
	i["Copy"] = th.Icon().Src(icons.ContentContentCopy)
	i["Paste"] = th.Icon().Src(icons.ContentContentPaste)
	i["Sidebar"] = th.Icon().Src(icons.ActionChromeReaderMode)
	i["Filter"] = th.Icon().Src(icons.ContentFilterList)
	i["FilterAll"] = th.Icon().Src(icons.ActionDoneAll)
	i["FilterNone"] = th.Icon().Src(icons.ContentBlock)
	i["Build"] = th.Icon().Src(icons.ActionBuild)
	i["Folded"] = th.Icon().Src(icons.NavigationChevronRight)
	i["Unfolded"] = th.Icon().Src(icons.NavigationExpandMore)
	i["HideAll"] = th.Icon().Src(icons.NavigationUnfoldLess)
	i["ShowAll"] = th.Icon().Src(icons.NavigationUnfoldMore)
	i["HideItem"] = th.Icon().Src(icons.ActionVisibilityOff)
	i["ShowItem"] = th.Icon().Src(icons.ActionVisibility)
	i["TRC"] = th.Icon().Src(icons.ActionSearch)
	i["DBG"] = th.Icon().Src(icons.ActionBugReport)
	i["INF"] = th.Icon().Src(icons.ActionInfo)
	i["WRN"] = th.Icon().Src(icons.ActionHelp)
	i["CHK"] = th.Icon().Src(icons.AlertWarning)
	i["ERR"] = th.Icon().Src(icons.AlertError)
	i["FTL"] = th.Icon().Src(icons.ImageFlashOn)
	i["Delete"] = th.Icon().Src(icons.ActionDelete)
	i["Send"] = th.Icon().Src(icons.ContentSend)
	return i
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
