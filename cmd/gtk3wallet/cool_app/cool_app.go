// This file includes wrapers for symbols included since GTK 3.12, and
// and should not be included in a build intended to target any older GTK
// versions.  To target an older build, such as 3.10, use
// 'go build -tags gtk_3_10'.  Otherwise, if no build tags are used, GTK 3.12
// is assumed and this file is built.
// +build !gtk_3_6,!gtk_3_8,!gtk_3_10

package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/d2r2/gotk3/gdk"
	"github.com/d2r2/gotk3/glib"
	"github.com/d2r2/gotk3/gtk"
	"github.com/d2r2/gotk3/pango"
	"github.com/davecgh/go-spew/spew"
)

// ========================================================================================
// ************************* GTK GUI UTILITIES SECTION START ******************************
// ========================================================================================
//	In real application copy this section to separate file as utilities functions to simplify
//	creation of GLIB/GTK+ components and widgets, including menus, dialog boxes, messages,
//	application settings and so on...

// SetupHeader construct Header widget with standard initialization.
func SetupHeader(title, subtitle string, showCloseButton bool) (*gtk.HeaderBar, error) {
	hdr, err := gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	hdr.SetShowCloseButton(showCloseButton)
	hdr.SetTitle(title)
	if subtitle != "" {
		hdr.SetSubtitle(subtitle)
	}
	return hdr, nil
}

// SetupMenuItemWithIcon construct MenuItem widget with icon image.
func SetupMenuItemWithIcon(label, detailedAction string, icon *glib.Icon) (*glib.MenuItem, error) {
	mi, err := glib.MenuItemNew(label, detailedAction)
	if err != nil {
		return nil, err
	}
	//mi.SetAttributeValue("verb-icon", iconNameVar)
	mi.SetIcon(icon)
	return mi, nil
}

// SetupMenuItemWithThemedIcon construct MenuItem widget with image
// taken by iconName from themed icons image lib.
func SetupMenuItemWithThemedIcon(label, detailedAction, iconName string) (*glib.MenuItem, error) {
	iconNameVar, err := glib.VariantStringNew(iconName)
	if err != nil {
		return nil, err
	}
	mi, err := glib.MenuItemNew(label, detailedAction)
	if err != nil {
		return nil, err
	}
	mi.SetAttributeValue("verb-icon", iconNameVar)
	return mi, nil
}

// SetupToolButton construct ToolButton widget with standart initialization.
func SetupToolButton(themedIconName, label string) (*gtk.ToolButton, error) {
	var btn *gtk.ToolButton
	var img *gtk.Image
	var err error
	if themedIconName != "" {
		img, err = gtk.ImageNewFromIconName(themedIconName, gtk.ICON_SIZE_BUTTON)
		if err != nil {
			return nil, err
		}
	}

	btn, err = gtk.ToolButtonNew(img, label)
	if err != nil {
		return nil, err
	}
	return btn, nil
}

// SetupButtonWithThemedImage construct Button widget with image
// taken by themedIconName from themed icons image lib.
func SetupButtonWithThemedImage(themedIconName string) (*gtk.Button, error) {
	img, err := gtk.ImageNewFromIconName(themedIconName, gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	btn, err := gtk.ButtonNew()
	if err != nil {
		return nil, err
	}

	btn.Add(img)

	return btn, nil
}

func getPixbufFromBytes(bytes []byte) (*gdk.Pixbuf, error) {
	b2, err := glib.BytesNew(bytes)
	if err != nil {
		return nil, err
	}
	ms, err := glib.MemoryInputStreamFromBytesNew(b2)
	if err != nil {
		return nil, err
	}
	pb, err := gdk.PixbufNewFromStream(&ms.InputStream, nil)
	if err != nil {
		return nil, err
	}
	return pb, nil
}

func getPixbufAnimationFromBytes(bytes []byte) (*gdk.PixbufAnimation, error) {
	b2, err := glib.BytesNew(bytes)
	if err != nil {
		return nil, err
	}
	ms, err := glib.MemoryInputStreamFromBytesNew(b2)
	if err != nil {
		return nil, err
	}
	pba, err := gdk.PixbufAnimationNewFromStream(&ms.InputStream, nil)
	if err != nil {
		return nil, err
	}
	return pba, nil
}

// SetupMenuButtonWithThemedImage construct MenuButton widget with image
// taken by themedIconName from themed icons image lib.
func SetupMenuButtonWithThemedImage(themedIconName string) (*gtk.MenuButton, error) {
	img, err := gtk.ImageNewFromIconName(themedIconName, gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}

	btn, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}

	btn.Add(img)

	return btn, nil
}

// AppendSectionAsHorzButtons used for Popover widget menu
// as a hint to display items as a horizontal buttons.
func AppendSectionAsHorzButtons(main, section *glib.Menu) error {
	val1, err := glib.VariantStringNew("horizontal-buttons")
	if err != nil {
		return err
	}
	mi1, err := glib.MenuItemNew("", "")
	if err != nil {
		return err
	}
	mi1.SetSection(section)
	mi1.SetAttributeValue("display-hint", val1)
	main.AppendItem(mi1)
	//section.AppendItem(mi1)
	return nil
}

// DialogButton simplify Dialog window initialization.
// Keep all necessary information about how attached
// dialog button should look and act.
type DialogButton struct {
	Text      string
	Response  gtk.ResponseType
	Default   bool
	Customize func(button *gtk.Button) error
}

// GetActiveWindow find real active window in application running.
func GetActiveWindow(win *gtk.Window) (*gtk.Window, error) {
	app, err := win.GetApplication()
	if err != nil {
		return nil, err
	}
	return app.GetActiveWindow(), nil
}

// IsResponseYes gives true if dialog window
// responded with gtk.RESPONSE_YES.
func IsResponseYes(response gtk.ResponseType) bool {
	return response == gtk.RESPONSE_YES
}

// IsResponseNo gives true if dialog window
// responded with gtk.RESPONSE_NO.
func IsResponseNo(response gtk.ResponseType) bool {
	return response == gtk.RESPONSE_NO
}

// IsResponseNone gives true if dialog window
// responded with gtk.RESPONSE_NONE.
func IsResponseNone(response gtk.ResponseType) bool {
	return response == gtk.RESPONSE_NONE
}

// IsResponseOk gives true if dialog window
// responded with gtk.RESPONSE_OK.
func IsResponseOk(response gtk.ResponseType) bool {
	return response == gtk.RESPONSE_OK
}

// IsResponseCancel gives true if dialog window
// responded with gtk.RESPONSE_CANCEL.
func IsResponseCancel(response gtk.ResponseType) bool {
	return response == gtk.RESPONSE_CANCEL
}

// IsResponseReject gives true if dialog window
// responded with gtk.RESPONSE_REJECT.
func IsResponseReject(response gtk.ResponseType) bool {
	return response == gtk.RESPONSE_REJECT
}

// IsResponseClose gives true if dialog window
// responded with gtk.RESPONSE_CLOSE.
func IsResponseClose(response gtk.ResponseType) bool {
	return response == gtk.RESPONSE_CLOSE
}

// IsResponseDeleteEvent gives true if dialog window
// responded with gtk.RESPONSE_DELETE_EVENT.
func IsResponseDeleteEvent(response gtk.ResponseType) bool {
	return response == gtk.RESPONSE_DELETE_EVENT
}

// PrintDialogResponse print and debug dialog responce.
func PrintDialogResponse(response gtk.ResponseType) {
	if IsResponseNo(response) {
		log.Println("Dialog result = NO")
	} else if IsResponseYes(response) {
		log.Println("Dialog result = YES")
	} else if IsResponseNone(response) {
		log.Println("Dialog result = NONE")
	} else if IsResponseOk(response) {
		log.Println("Dialog result = OK")
	} else if IsResponseReject(response) {
		log.Println("Dialog result = REJECT")
	} else if IsResponseCancel(response) {
		log.Println("Dialog result = CANCEL")
	} else if IsResponseClose(response) {
		log.Println("Dialog result = CLOSE")
	} else if IsResponseDeleteEvent(response) {
		log.Println("Dialog result = DELETE_EVENT")
	}
}

// DialogParagraph is an object which keep text paragraph added
// to dialog window, complemented with all necessary format options.
type DialogParagraph struct {
	Text          string
	Markup        bool
	HorizAlign    gtk.Align
	Justify       gtk.Justification
	Ellipsize     pango.EllipsizeMode
	MaxWidthChars int
}

func NewDialogParagraph(text string) *DialogParagraph {
	v := &DialogParagraph{Text: text, HorizAlign: gtk.Align(-1), Justify: gtk.Justification(-1),
		Ellipsize: pango.EllipsizeMode(-1), MaxWidthChars: -1}
	return v
}

func NewMarkupDialogParagraph(text string) *DialogParagraph {
	v := &DialogParagraph{Text: text, Markup: true, HorizAlign: gtk.Align(-1), Justify: gtk.Justification(-1),
		Ellipsize: pango.EllipsizeMode(-1), MaxWidthChars: -1}
	return v
}

func (v *DialogParagraph) SetHorizAlign(align gtk.Align) *DialogParagraph {
	v.HorizAlign = align
	return v
}

func (v *DialogParagraph) SetJustify(justify gtk.Justification) *DialogParagraph {
	v.Justify = justify
	return v
}

func (v *DialogParagraph) SetEllipsize(ellipsize pango.EllipsizeMode) *DialogParagraph {
	v.Ellipsize = ellipsize
	return v
}

func (v *DialogParagraph) SetMaxWidthChars(maxWidthChars int) *DialogParagraph {
	v.MaxWidthChars = maxWidthChars
	return v
}

func (v *DialogParagraph) createLabel() (*gtk.Label, error) {
	lbl, err := gtk.LabelNew("")
	if err != nil {
		return nil, err
	}
	if v.Markup {
		lbl.SetMarkup(v.Text)
	} else {
		lbl.SetText(v.Text)
	}
	if v.HorizAlign != gtk.Align(-1) {
		lbl.SetHAlign(v.HorizAlign)
	}
	if v.Justify != gtk.Justification(-1) {
		lbl.SetJustify(v.Justify)
	}
	if v.Ellipsize != pango.EllipsizeMode(-1) {
		lbl.SetEllipsize(v.Ellipsize)
	}
	if v.MaxWidthChars != -1 {
		lbl.SetMaxWidthChars(v.MaxWidthChars)
	}
	return lbl, nil
}

func TextToDialogParagraphs(lines []string) []*DialogParagraph {
	var msgs []*DialogParagraph
	for _, line := range lines {
		msgs = append(msgs, NewDialogParagraph(line))
	}
	return msgs
}

func TextToMarkupDialogParagraphs(lines []string) []*DialogParagraph {
	var msgs []*DialogParagraph
	for _, line := range lines {
		msgs = append(msgs, NewMarkupDialogParagraph(line))
	}
	return msgs
}

// SetupMessageDialog construct MessageDialog widget with customized settings.
func SetupMessageDialog(parent *gtk.Window, markupTitle string, secondaryMarkupTitle string,
	paragraphs []*DialogParagraph, addButtons []DialogButton,
	addExtraControls func(area *gtk.Box) error) (*gtk.MessageDialog, error) {

	var active *gtk.Window
	var err error

	if parent != nil {
		active, err = GetActiveWindow(parent)
		if err != nil {
			return nil, err
		}
	}

	dlg, err := gtk.MessageDialogNew(active, /*gtk.DIALOG_MODAL|*/
		gtk.DIALOG_USE_HEADER_BAR, gtk.MESSAGE_WARNING, gtk.BUTTONS_NONE, nil, nil)
	if err != nil {
		return nil, err
	}
	if active != nil {
		dlg.SetTransientFor(active)
	}
	dlg.SetMarkup(markupTitle)
	if secondaryMarkupTitle != "" {
		dlg.FormatSecondaryMarkup(secondaryMarkupTitle)
	}

	for _, button := range addButtons {
		btn, err := dlg.AddButton(button.Text, button.Response)
		if err != nil {
			return nil, err
		}

		if button.Default {
			dlg.SetDefaultResponse(button.Response)
		}

		if button.Customize != nil {
			err := button.Customize(btn)
			if err != nil {
				return nil, err
			}
		}
	}

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}

	grid.SetRowSpacing(6)
	grid.SetHAlign(gtk.ALIGN_CENTER)

	box, err := dlg.GetMessageArea()
	if err != nil {
		return nil, err
	}
	box.Add(grid)

	col := 1
	row := 0

	// add empty line after title
	paragraphs = append([]*DialogParagraph{NewDialogParagraph("")}, paragraphs...)

	for _, paragraph := range paragraphs {
		lbl, err := paragraph.createLabel()
		if err != nil {
			return nil, err
		}
		grid.Attach(lbl, col, row, 1, 1)
		row++
	}

	if addExtraControls != nil {
		box1, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
		if err != nil {
			return nil, err
		}
		grid.Attach(box1, col, row, 1, 1)

		err = addExtraControls(box1)
		if err != nil {
			return nil, err
		}
	}

	box.ShowAll()

	return dlg, nil
}

// RunMessageDialog construct and run MessageDialog widget with customized settings.
func RunMessageDialog(parent *gtk.Window, markupTitle string, secondaryMarkupTitle string,
	paragraphs []*DialogParagraph, ignoreCloseBox bool, addButtons []DialogButton,
	addExtraControls func(area *gtk.Box) error) (gtk.ResponseType, error) {

	dlg, err := SetupMessageDialog(parent, markupTitle, secondaryMarkupTitle,
		paragraphs, addButtons, addExtraControls)
	if err != nil {
		return 0, err
	}
	defer dlg.Destroy()

	dlg.ShowAll()
	res := dlg.Run()
	for gtk.ResponseType(res) == gtk.RESPONSE_NONE || gtk.ResponseType(res) == gtk.RESPONSE_DELETE_EVENT && ignoreCloseBox {
		res = dlg.Run()
	}
	return gtk.ResponseType(res), nil
}

// SetupDialog construct Dialog widget with customized settings.
func SetupDialog(parent *gtk.Window, messageType gtk.MessageType, userHeaderbar bool,
	title string, paragraphs []*DialogParagraph, addButtons []DialogButton,
	addExtraControls func(area *gtk.Box) error) (*gtk.Dialog, error) {

	var active *gtk.Window
	var err error

	if parent != nil {
		active, err = GetActiveWindow(parent)
		if err != nil {
			return nil, err
		}
	}

	flags := gtk.DIALOG_MODAL
	if userHeaderbar {
		flags |= gtk.DIALOG_USE_HEADER_BAR
	}
	dlg, err := gtk.DialogWithFlagsNew(title, active, flags)
	if err != nil {
		return nil, err
	}

	dlg.SetDefaultSize(100, 100)
	dlg.SetTransientFor(active)
	dlg.SetDeletable(false)

	var img *gtk.Image
	size := gtk.ICON_SIZE_DIALOG
	if userHeaderbar {
		size = gtk.ICON_SIZE_LARGE_TOOLBAR
	}
	var iconName string
	switch messageType {
	case gtk.MESSAGE_WARNING:
		iconName = "dialog-warning"
	case gtk.MESSAGE_ERROR:
		iconName = "dialog-error"
	case gtk.MESSAGE_INFO:
		iconName = "dialog-information"
	case gtk.MESSAGE_QUESTION:
		iconName = "dialog-question"
	}

	if iconName != "" {
		img, err = gtk.ImageNewFromIconName(iconName, size)
		if err != nil {
			return nil, err
		}
	}

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}

	grid.SetBorderWidth(10)
	grid.SetColumnSpacing(10)
	grid.SetRowSpacing(6)
	grid.SetHAlign(gtk.ALIGN_CENTER)

	box, err := dlg.GetContentArea()
	if err != nil {
		return nil, err
	}
	box.Add(grid)

	if img != nil {
		if userHeaderbar {
			hdr, err := dlg.GetHeaderBar()
			if err != nil {
				return nil, err
			}

			hdr.PackStart(img)
		} else {
			grid.Attach(img, 0, 0, 1, 1)
		}
	}

	for _, button := range addButtons {
		btn, err := dlg.AddButton(button.Text, button.Response)
		if err != nil {
			return nil, err
		}

		if button.Default {
			dlg.SetDefaultResponse(button.Response)
		}

		if button.Customize != nil {
			err := button.Customize(btn)
			if err != nil {
				return nil, err
			}
		}
	}

	col := 1
	row := 0

	for _, paragraph := range paragraphs {
		lbl, err := paragraph.createLabel()
		if err != nil {
			return nil, err
		}
		grid.Attach(lbl, col, row, 1, 1)
		row++
	}

	if addExtraControls != nil {
		box1, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
		if err != nil {
			return nil, err
		}
		grid.Attach(box1, col, row, 1, 1)

		err = addExtraControls(box1)
		if err != nil {
			return nil, err
		}
	}

	_, w := dlg.GetPreferredWidth()
	_, h := dlg.GetPreferredHeight()
	dlg.Resize(w, h)

	return dlg, nil
}

// RunDialog construct and run Dialog widget with customized settings.
func RunDialog(parent *gtk.Window, messageType gtk.MessageType, userHeaderbar bool,
	title string, paragraphs []*DialogParagraph, ignoreCloseBox bool, addButtons []DialogButton,
	addExtraControls func(area *gtk.Box) error) (gtk.ResponseType, error) {

	dlg, err := SetupDialog(parent, messageType, userHeaderbar, title,
		paragraphs, addButtons, addExtraControls)
	if err != nil {
		return 0, err
	}
	defer dlg.Destroy()

	//dlg.ShowAll()
	dlg.ShowAll()
	res := dlg.Run()
	for gtk.ResponseType(res) == gtk.RESPONSE_NONE || gtk.ResponseType(res) == gtk.RESPONSE_DELETE_EVENT && ignoreCloseBox {
		res = dlg.Run()
	}
	return gtk.ResponseType(res), nil
}

func ErrorMessage(parent *gtk.Window, titleMarkup string, text []*DialogParagraph) error {
	buttons := []DialogButton{
		{"_OK", gtk.RESPONSE_OK, false, nil},
	}
	_, err := RunMessageDialog(parent, titleMarkup, "", text, false, buttons, nil)
	if err != nil {
		return err
	}
	return nil
}

func QuestionDialog(parent *gtk.Window, title string,
	messages []*DialogParagraph, defaultYes bool) (bool, error) {

	title2 := spew.Sprintf("%s", title)
	buttons := []DialogButton{
		{"_YES", gtk.RESPONSE_YES, defaultYes, nil},
		{"_NO", gtk.RESPONSE_NO, !defaultYes, nil},
	}
	response, err := RunDialog(parent, gtk.MESSAGE_QUESTION, true, title2,
		messages, false, buttons, nil)
	if err != nil {
		return false, err
	}
	PrintDialogResponse(response)
	return IsResponseYes(response), nil
}

// GetActionNameAndState display status of action-with-state, which used in
// menu-with-state behavior. Convenient for debug purpose.
func GetActionNameAndState(act *glib.SimpleAction) (string, *glib.Variant, error) {
	name, err := act.GetName()
	if err != nil {
		return "", nil, err
	}
	state := act.GetState()
	return name, state, nil
}

// SetMargins set margins of a widget to the passed values,
// replacing 4 calls with only one.
func SetMargins(widget gtk.IWidget, left int, top int, right int, bottom int) {
	w := widget.GetWidget()
	w.SetMarginStart(left)
	w.SetMarginTop(top)
	w.SetMarginEnd(right)
	w.SetMarginBottom(bottom)
}

// SetAllMargins set all margins of a widget to the same value.
func SetAllMargins(widget gtk.IWidget, margin int) {
	SetMargins(widget, margin, margin, margin, margin)
}

// AppendValues append multiple values to a row in a list store.
func AppendValues(ls *gtk.ListStore, values ...interface{}) (*gtk.TreeIter, error) {
	iter := ls.Append()
	for i := 0; i < len(values); i++ {
		err := ls.SetValue(iter, i, values[i])
		if err != nil {
			return nil, err
		}
	}
	return iter, nil
}

// CreateNameValueCombo create a GtkComboBox that holds
// a set of name/value pairs where the name is displayed.
func CreateNameValueCombo(keyValues []struct{ value, key string }) (*gtk.ComboBox, error) {
	ls, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		return nil, err
	}

	for _, item := range keyValues {
		_, err = AppendValues(ls, item.value, item.key)
		if err != nil {
			return nil, err
		}
	}

	cb, err := gtk.ComboBoxNew()
	if err != nil {
		return nil, err
	}
	err = UpdateNameValueCombo(cb, keyValues)
	if err != nil {
		return nil, err
	}
	cb.SetFocusOnClick(false)
	cb.SetIDColumn(1)
	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}
	cell.SetAlignment(0, 0)
	cb.PackStart(cell, false)
	cb.AddAttribute(cell, "text", 0)

	return cb, nil
}

// UpdateNameValueCombo update GtkComboBox list of name/value pairs.
func UpdateNameValueCombo(cb *gtk.ComboBox, keyValues []struct{ value, key string }) error {
	ls, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		return err
	}

	for _, item := range keyValues {
		_, err = AppendValues(ls, item.value, item.key)
		if err != nil {
			return err
		}
	}

	cb.SetModel(ls)
	return nil
}

// GetComboValue return GtkComboBox selected value from specific column.
func GetComboValue(cb *gtk.ComboBox, columnID int) (*glib.Value, error) {
	ti, err := cb.GetActiveIter()
	if err != nil {
		return nil, err
	}
	tm, err := cb.GetModel()
	if err != nil {
		return nil, err
	}
	val, err := tm.GetValue(ti, 0)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// GetGtkVersion return actually installed GTK+ version.
func GetGtkVersion() (magor, minor, micro uint) {
	magor = gtk.GetMajorVersion()
	minor = gtk.GetMinorVersion()
	micro = gtk.GetMicroVersion()
	return
}

// GetGlibVersion return actually installed GLIB version.
func GetGlibVersion() (magor, minor, micro uint) {
	magor = glib.GetMajorVersion()
	minor = glib.GetMinorVersion()
	micro = glib.GetMicroVersion()
	return
}

// GetSchema obtains glib.SettingsSchema from glib.Settings.
func GetSchema(v *glib.Settings) (*glib.SettingsSchema, error) {
	val, err := v.GetProperty("settings-schema")
	if err != nil {
		return nil, err
	}
	if schema, ok := val.(*glib.SettingsSchema); ok {
		return schema, nil
	} else {
		return nil, errors.New("GLib settings-schema property is not convertible to SettingsSchema")
	}
}

// FixProgressBarCSS eliminate issue with default GtkProgressBar control formating.
func applyStyleCSS(widget *gtk.Widget, css string) error {
	//	provider, err := gtk.CssProviderNew()
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}
	err = provider.LoadFromData(css)
	if err != nil {
		return err
	}
	sc, err := widget.GetStyleContext()
	if err != nil {
		return err
	}
	sc.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_USER)
	//sc.AddClass("osd")
	return nil
}

// Binding cache link between Key string identifier and GLIB object property.
// Code taken from https://github.com/gnunn1/tilix project.
type Binding struct {
	Key      string
	Object   glib.IObject
	Property string
	Flags    glib.SettingsBindFlags
}

// BindingHelper is a bookkeeping class that keeps track of objects which are
// binded to a GSettings object so they can be unbinded later. it
// also supports the concept of deferred bindings where a binding
// can be added but is not actually attached to a Settings object
// until one is set.
type BindingHelper struct {
	bindings []Binding
	settings *glib.Settings
}

// BindingHelperNew creates new BindingHelper object.
func BindingHelperNew(settings *glib.Settings) *BindingHelper {
	bh := &BindingHelper{settings: settings}
	return bh
}

// SetSettings will replace underlying GLIB Settings object to unbind
// previously set bindings and re-bind to the new settings automatically.
func (v *BindingHelper) SetSettings(value *glib.Settings) {
	if value != v.settings {
		if v.settings != nil {
			v.Unbind()
		}
		v.settings = value
		if v.settings != nil {
			v.bindAll()
		}
	}
}

func (v *BindingHelper) bindAll() {
	if v.settings != nil {
		for _, b := range v.bindings {
			v.settings.Bind(b.Key, b.Object, b.Property, b.Flags)
		}
	}
}

// addBind add a binding to the list
func (v *BindingHelper) addBind(key string, object glib.IObject, property string, flags glib.SettingsBindFlags) {
	v.bindings = append(v.bindings, Binding{key, object, property, flags})
}

// Bind add a binding to list and binds to Settings if it is set.
func (v *BindingHelper) Bind(key string, object glib.IObject, property string, flags glib.SettingsBindFlags) {
	v.addBind(key, object, property, flags)
	if v.settings != nil {
		v.settings.Bind(key, object, property, flags)
	}
}

// Unbind all added binds from settings object.
func (v *BindingHelper) Unbind() {
	for _, b := range v.bindings {
		v.settings.Unbind(b.Object, b.Property)
	}
}

// Clear unbind all bindings and clears list of bindings.
func (v *BindingHelper) Clear() {
	v.Unbind()
	v.bindings = nil
}

// SettingsArray is a way how to create multiple (indexed) GLib setting's group.
// For instance, multiple backup profiles with identical
// settings inside of each profile. Either each backup profile may
// contain more than one data source for backup.
type SettingsArray struct {
	settings *glib.Settings
	arrayID  string
}

func NewSettingsArray(settings *glib.Settings, arrayID string) *SettingsArray {
	v := &SettingsArray{settings: settings, arrayID: arrayID}
	return v
}

func (v *SettingsArray) DeleteNode(childSettings *glib.Settings, nodeID string) error {
	schema, err := GetSchema(childSettings)
	if err != nil {
		return err
	}
	keys := schema.ListKeys()
	for _, key := range keys {
		childSettings.Reset(key)
	}

	sources := v.settings.GetStrv(v.arrayID)
	var newSources []string
	for _, id := range sources {
		if id != nodeID {
			newSources = append(newSources, id)
		}
	}
	v.settings.SetStrv(v.arrayID, newSources)
	return nil
}

func (v *SettingsArray) AddNode() (nodeID string, err error) {
	sources := v.settings.GetStrv(v.arrayID)
	var ni int
	if len(sources) > 0 {
		ni, err = strconv.Atoi(sources[len(sources)-1])
		if err != nil {
			return "", err
		}
		ni++
	}
	//lg.Println(spew.Sprintf("New node id: %+v", ni))
	sources = append(sources, strconv.Itoa(ni))
	v.settings.SetStrv(v.arrayID, sources)
	return sources[len(sources)-1], nil
}

func (v *SettingsArray) GetArrayIDs() []string {
	sources := v.settings.GetStrv(v.arrayID)
	return sources
}

// ========================================================================================
// ************************* GTK GUI UTILITIES SECTION END ********************************
// ========================================================================================

// String constants used as titles/identifiers
var (
	APP_ID string = "org.d2r2.gotk3.cool_app_1"

	APP_TITLE string = "Cool App (GTK+ 3 UI adaptation for golang based on imporved GOTK3)"

	//Preference Constants
	SETTINGS_ID string = APP_ID + ".Settings"

	SETTINGS_AUTO_HIDE_MOUSE_KEY                   string = "auto-hide-mouse"
	SETTINGS_DONT_SHOW_ABOUT_ON_STARTUP_KEY        string = "dont-show-about-dialog-on-startup"
	SETTINGS_PROMPT_ON_NEW_SESSION_KEY             string = "prompt-on-new-session"
	SETTINGS_ENABLE_TRANSPARENCY_KEY               string = "enable-transparency"
	SETTINGS_CLOSE_WITH_LAST_SESSION_KEY           string = "close-with-last-session"
	SETTINGS_APP_TITLE_KEY                         string = "app-title"
	SETTINGS_CONTROL_CLICK_TITLE_KEY               string = "control-click-titlebar"
	SETTINGS_INHERIT_WINDOW_STATE_KEY              string = "new-window-inherit-state"
	SETTINGS_USE_OVERLAY_SCROLLBAR_KEY             string = "use-overlay-scrollbar"
	SETTINGS_TERMINAL_FOCUS_FOLLOWS_MOUSE_KEY      string = "focus-follow-mouse"
	SETTINGS_MIDDLE_CLICK_CLOSE_KEY                string = "middle-click-close"
	SETTINGS_CONTROL_SCROLL_ZOOM_KEY               string = "control-scroll-zoom"
	SETTINGS_WINDOW_SAVE_STATE_KEY                 string = "window-save-state"
	SETTINGS_NOTIFY_ON_PROCESS_COMPLETE_KEY        string = "notify-on-process-complete"
	SETTINGS_PASTE_ADVANCED_DEFAULT_KEY            string = "paste-advanced-default"
	SETTINGS_UNSAFE_PASTE_ALERT_KEY                string = "unsafe-paste-alert"
	SETTINGS_STRIP_FIRST_COMMENT_CHAR_ON_PASTE_KEY string = "paste-strip-first-char"
	SETTINGS_COPY_ON_SELECT_KEY                    string = "copy-on-select"
	SETTINGS_WINDOW_STYLE_KEY                      string = "window-style"

	SETTINGS_TERMINAL_TITLE_STYLE_KEY          string = "terminal-title-style"
	SETTINGS_TERMINAL_TITLE_STYLE_VALUE_NORMAL string = "normal"
	SETTINGS_TERMINAL_TITLE_STYLE_VALUE_SMALL  string = "small"
	SETTINGS_TERMINAL_TITLE_STYLE_VALUE_NONE   string = "none"

	SETTINGS_TAB_POSITION_KEY string = "tab-position"

	// Theme Settings
	SETTINGS_THEME_VARIANT_KEY          string = "theme-variant"
	SETTINGS_THEME_VARIANT_SYSTEM_VALUE string = "system"
	SETTINGS_THEME_VARIANT_LIGHT_VALUE  string = "light"
	SETTINGS_THEME_VARIANT_DARK_VALUE   string = "dark"

	SETTINGS_BACKGROUND_IMAGE_KEY                string = "background-image"
	SETTINGS_BACKGROUND_IMAGE_MODE_KEY           string = "background-image-mode"
	SETTINGS_BACKGROUND_IMAGE_MODE_SCALE_VALUE   string = "scale"
	SETTINGS_BACKGROUND_IMAGE_MODE_TILE_VALUE    string = "tile"
	SETTINGS_BACKGROUND_IMAGE_MODE_CENTER_VALUE  string = "center"
	SETTINGS_BACKGROUND_IMAGE_MODE_STRETCH_VALUE string = "stretch"

	SETTINGS_SIDEBAR_RIGHT string = "sidebar-on-right"

	SETTINGS_TERMINAL_TITLE_SHOW_WHEN_SINGLE_KEY string = "terminal-title-show-when-single"

	SETTINGS_USE_TABS_KEY string = "use-tabs"
)

// Keeps here global fullscreen data, which available throughout the application.
// TODO: add thread-safe protection.
type FullscreenGlobalData struct {
	EventBox             *gtk.EventBox
	Revealer             *gtk.Revealer
	FullscreeAction      *glib.Action
	RightTopMenuButton   *gtk.MenuButton
	InFullscreenEventBox bool
	Timer                *time.Timer
}

// Demonstration of Popover menu creation with different kind of menus:
//	1) Regular menus (with text)
//	2) Button menus with themed icons (native GTK+ icons) and loaded file-base icons
//	3) Submenus
func createMenuModelForPopover() (glib.IMenuModel, error) {

	main, err := glib.MenuNew()
	if err != nil {
		return nil, err
	}

	var section *glib.Menu
	var item *glib.MenuItem

	// New menu section (with buttons)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	item, err = SetupMenuItemWithThemedIcon("About", "win.AboutAction", "help-about-symbolic")
	if err != nil {
		return nil, err
	}
	section.AppendItem(item)
	item, err = SetupMenuItemWithThemedIcon("Preference", "win.PreferenceAction", "preferences-other-symbolic")
	if err != nil {
		return nil, err
	}
	section.AppendItem(item)
	item, err = SetupMenuItemWithThemedIcon("Fullscreen", "win.FullscreenAction", "view-fullscreen-symbolic")
	if err != nil {
		return nil, err
	}
	section.AppendItem(item)
	err = AppendSectionAsHorzButtons(main, section)
	if err != nil {
		return nil, err
	}

	// New menu section (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	section.Append("Enable/Disable option", "win.CheckBoxAction")
	main.AppendSection("", section)

	// New menu section (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	section.Append("Red", "win.ChooseColor('red')")
	section.Append("Green", "win.ChooseColor('green')")
	section.Append("Blue", "win.ChooseColor('blue')")
	main.AppendSection("", section)

	// New menu section (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	section.Append("Dialog demo 1", "win.DialogAction1")
	section.Append("Dialog demo 2", "win.DialogAction2")
	section.Append("Dialog demo 3", "win.DialogAction3")
	section.Append("Dialog demo 4", "win.DialogAction4")
	section.Append("Dialog demo 5", "win.DialogAction5")
	main.AppendSection("", section)

	// New menu section with submenu (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	subMenu, err := glib.MenuNew()
	if err != nil {
		return nil, err
	}
	subMenu.Append("Fullscreen mode", "win.FullscreenAction")
	section.AppendSubmenu("Submenu", subMenu)
	main.AppendSection("", section)

	// New menu section with submenu (with buttons)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	item, err = SetupMenuItemWithThemedIcon("Quit application", "win.QuitAction", "application-exit-symbolic")
	if err != nil {
		return nil, err
	}
	section.AppendItem(item)
	err = AppendSectionAsHorzButtons(main, section)
	if err != nil {
		return nil, err
	}

	return main, nil
}

// Creates GLIB's MenuModel interface for using in GTK MenuBar.
func createMenuModelForMenuBar() (glib.IMenuModel, error) {
	rootMenu, err := glib.MenuNew()
	if err != nil {
		return nil, err
	}
	main, err := glib.MenuNew()
	if err != nil {
		return nil, err
	}
	rootMenu.AppendSubmenu("Main menu", main)

	var section *glib.Menu
	var item *glib.MenuItem

	// New menu section (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	main.AppendSection("", section)
	section.Append("Fullscreen mode", "win.FullscreenAction")

	// New menu section (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	section.Append("Enable/Disable option", "win.CheckBoxAction")
	main.AppendSection("", section)

	// New menu section (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	main.AppendSection("", section)
	section.Append("Red", "win.ChooseColor('red')")
	section.Append("Green", "win.ChooseColor('green')")
	section.Append("Blue", "win.ChooseColor('blue')")

	// New menu section with submeny (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	main.AppendSection("", section)
	subMenu, err := glib.MenuNew()
	if err != nil {
		return nil, err
	}
	subMenu.Append("Fullscreen", "win.FullscreenAction")
	section.AppendSubmenu("Submenu", subMenu)

	// New menu section (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	section.Append("Dialog demo 1", "win.DialogAction1")
	section.Append("Dialog demo 2", "win.DialogAction2")
	section.Append("Dialog demo 3", "win.DialogAction3")
	section.Append("Dialog demo 4", "win.DialogAction4")
	section.Append("Dialog demo 5", "win.DialogAction5")
	main.AppendSection("", section)

	// New menu section with submeny (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	file, err := glib.FileForPathNew("./icons/ajax-loader-gears_32x32.gif")
	//file, err := glib.FileForPathNew("./icons/com.gexperts.Tilix.png")
	if err != nil {
		return nil, err
	}
	icon, err := glib.FileIconNew(file)
	if err != nil {
		return nil, err
	}
	item, err = SetupMenuItemWithIcon("Preference dialog demo", "win.PreferenceAction", &icon.Icon)
	if err != nil {
		return nil, err
	}
	section.AppendItem(item)
	icon2, err := glib.ThemedIconNew("help-about-symbolic")
	if err != nil {
		return nil, err
	}
	item, err = SetupMenuItemWithIcon("About demo", "win.AboutAction", &icon2.Icon)
	if err != nil {
		return nil, err
	}
	section.AppendItem(item)
	main.AppendSection("", section)

	// New menu section with submeny (automatically separeted by dividers)
	section, err = glib.MenuNew()
	if err != nil {
		return nil, err
	}
	icon3, err := glib.ThemedIconNew("application-exit-symbolic")
	if err != nil {
		return nil, err
	}
	item, err = SetupMenuItemWithIcon("Quit application", "win.QuitAction", &icon3.Icon)
	if err != nil {
		return nil, err
	}
	section.AppendItem(item)
	main.AppendSection("", section)

	return rootMenu, nil
}

// Creates GTK Toolbar actions attached to buttons.
func createToolbar() (*gtk.Toolbar, error) {
	tbx, err := gtk.ToolbarNew()
	if err != nil {
		return nil, err
	}
	tbx.SetStyle(gtk.TOOLBAR_BOTH_HORIZ)

	var tbtn *gtk.ToolButton
	var tdvd *gtk.SeparatorToolItem
	var img *gtk.Image

	img, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	img.SetFromIconName("application-exit-symbolic", gtk.ICON_SIZE_BUTTON)
	tbtn, err = gtk.ToolButtonNew(img, "")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.QuitAction")
	tbx.Add(tbtn)

	img, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	img.SetFromIconName("view-fullscreen-symbolic", gtk.ICON_SIZE_BUTTON)
	tbtn, err = gtk.ToolButtonNew(img, "")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.FullscreenAction")
	tbx.Add(tbtn)

	tdvd, err = gtk.SeparatorToolItemNew()
	if err != nil {
		return nil, err
	}
	tbx.Add(tdvd)

	tbtn, err = gtk.ToolButtonNew(nil, "Demo 1")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.DialogAction1")
	tbx.Add(tbtn)

	tbtn, err = gtk.ToolButtonNew(nil, "Demo 2")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.DialogAction2")
	tbx.Add(tbtn)

	tbtn, err = gtk.ToolButtonNew(nil, "Demo 3")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.DialogAction3")
	tbx.Add(tbtn)

	tbtn, err = gtk.ToolButtonNew(nil, "Demo 4")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.DialogAction4")
	tbx.Add(tbtn)

	tbtn, err = gtk.ToolButtonNew(nil, "Demo 5")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.DialogAction5")
	tbx.Add(tbtn)

	tdvd, err = gtk.SeparatorToolItemNew()
	if err != nil {
		return nil, err
	}
	tbx.Add(tdvd)

	img, err = gtk.ImageNewFromFile("./icons/ajax-loader-gears_32x32.gif")
	if err != nil {
		return nil, err
	}
	tbtn, err = gtk.ToolButtonNew(img, "")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.PreferenceAction")
	tbx.Add(tbtn)

	img, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	img.SetFromIconName("help-about-symbolic", gtk.ICON_SIZE_BUTTON)
	tbtn, err = gtk.ToolButtonNew(img, "")
	if err != nil {
		return nil, err
	}
	tbtn.SetActionName("win.AboutAction")
	tbx.Add(tbtn)

	return tbx, nil
}

// Create demonstration action-with-state with boolean logic (on/of).
// Action trigger is included.
func createCheckBoxAction() (glib.IAction, error) {
	v, err := glib.VariantBooleanNew(true)
	if err != nil {
		return nil, err
	}
	act, err := glib.SimpleActionStatefullNew("CheckBoxAction", nil, v)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, errors.New("error")
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		if state != nil && state.IsOfType(glib.VARIANT_TYPE_BOOLEAN) {
			state, err = glib.VariantBooleanNew(!state.GetBoolean())
			if err != nil {
				log.Fatal(err)
			}
			action.ChangeState(state)
		}
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create demonstration action-with-state with list logic (opt1/opt2/.../optN).
// Action trigger is included.
func createChooseAction() (glib.IAction, error) {
	v, err := glib.VariantStringNew("green")
	if err != nil {
		return nil, err
	}
	act, err := glib.SimpleActionStatefullNew("ChooseColor", glib.VARIANT_TYPE_STRING, v)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, errors.New("error")
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		act.ChangeState(param)
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create fullscreen action with boolean in the base (fullscreen on/off)
// Action trigger is included.
func createFullscreenAction(win *gtk.Window, data *FullscreenGlobalData) (glib.IAction, error) {
	v, err := glib.VariantBooleanNew(false)
	if err != nil {
		return nil, err
	}
	act, err := glib.SimpleActionStatefullNew("FullscreenAction", nil, v)
	if err != nil {
		return nil, err
	}
	data.FullscreeAction = &act.Action

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		if state != nil && state.IsOfType(glib.VARIANT_TYPE_BOOLEAN) {
			state, err = glib.VariantBooleanNew(!state.GetBoolean())
			if err != nil {
				log.Fatal(err)
			}
			action.ChangeState(state)

			if state.GetBoolean() {
				win.Fullscreen()
				data.EventBox.ShowAll()
				//data.Revealer.SetRevealChild(false)
				showFullscreenHeader(data, true)
				manageFullscreenHeaderHideTimer(data, time.Millisecond*2500, true)
			} else {
				data.EventBox.Hide()
				//data.Revealer.SetRevealChild(false)
				win.Unfullscreen()
			}
		}
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create regular exit app action.
// Action trigger is included.
func createQuitAction(win *gtk.Window) (glib.IAction, error) {
	act, err := glib.SimpleActionNew("QuitAction", nil)
	if err != nil {
		return nil, err
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		application, err := win.GetApplication()
		if err != nil {
			log.Fatal(err)
		}
		application.Quit()
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create regular "about dialog" action.
// Action trigger is included.
func createAboutAction(win *gtk.Window, gsSettings *glib.Settings) (glib.IAction, error) {
	act, err := glib.SimpleActionNew("AboutAction", nil)
	if err != nil {
		return nil, err
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		dlg, err := gtk.AboutDialogNew()
		if err != nil {
			log.Fatal(err)
		}

		dlg.SetAuthors([]string{"Written by Denis Dyakov <denis.dyakov@gmail.com>"})
		dlg.SetProgramName("Cool App")
		dlg.SetLogoIconName("face-cool-symbolic")
		dlg.SetVersion("v0.1")

		bh := BindingHelperNew(gsSettings)
		// Show about dialog on application startup
		cbAboutInfo, err := gtk.CheckButtonNewWithLabel("Do not show about information on app startup")
		if err != nil {
			log.Fatal(err)
		}
		bh.Bind(SETTINGS_DONT_SHOW_ABOUT_ON_STARTUP_KEY, cbAboutInfo, "active", glib.SETTINGS_BIND_DEFAULT)

		content, err := dlg.GetContentArea()
		if err != nil {
			log.Fatal(err)
		}
		content.Add(cbAboutInfo)
		content.ShowAll()

		var buf bytes.Buffer
		glibMajor, glibMinor, glibMicro := GetGlibVersion()
		glibBuildVersion := glib.GetBuildVersion()
		gtkMajor, gtkMinor, gtkMicro := GetGtkVersion()
		gtkBuildVersion := gtk.GetBuildVersion()
		buf.WriteString(fmt.Sprintln("This application built for education purpose and compose"))
		buf.WriteString(fmt.Sprintln("practices to write GTK+3 user interface in Go language."))
		buf.WriteString(fmt.Sprintln())
		buf.WriteString(fmt.Sprintln("Environment:"))
		buf.WriteString(fmt.Sprintln(fmt.Sprintf("GLIB compiled version %s, detected version %d.%d.%d",
			glibBuildVersion, glibMajor, glibMinor, glibMicro)))
		buf.WriteString(fmt.Sprintln(fmt.Sprintf("GTK+ compiled version %s, detected version %d.%d.%d",
			gtkBuildVersion, gtkMajor, gtkMinor, gtkMicro)))
		buf.WriteString(fmt.Sprintln(fmt.Sprintf("Application compiled with %s %s",
			runtime.Version(), runtime.GOARCH)))
		buf.WriteString(fmt.Sprintln())
		buf.WriteString(fmt.Sprintln("Features:"))
		buf.WriteString(fmt.Sprintln("- Actions as code entry points with states and stateless."))
		buf.WriteString(fmt.Sprintln("- Fullscreen mode code pattern out-of-the-box."))
		buf.WriteString(fmt.Sprintln("- Preference dialog demo with save/restore functionality out-of-the-box."))
		buf.WriteString(fmt.Sprintln("- Modern popover menu functionality (right upper corner button)."))
		buf.WriteString(fmt.Sprintln("- Various dialog's windows demonstations."))
		buf.WriteString(fmt.Sprintln())
		buf.WriteString(fmt.Sprint("Follow my golang projects on GitHub:"))
		dlg.SetComments(buf.String())

		dlg.SetWebsite("https://github.com/d2r2/")

		dlg.SetTransientFor(win)
		dlg.SetModal(true)
		dlg.ShowNow()
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create demonstration GTK MessageDialog action.
// Action trigger is included.
func createDialogAction1(win *gtk.Window) (glib.IAction, error) {
	act, err := glib.SimpleActionNew("DialogAction1", nil)
	if err != nil {
		return nil, err
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		title := "<span weight='bold' size='larger'>Configuration Issue Detected (demonstration)</span>"
		paragraphs := []*DialogParagraph{NewDialogParagraph("Message dialog based on GTK+ MessageDialog.").
			SetJustify(gtk.JUSTIFY_LEFT).SetHorizAlign(gtk.ALIGN_START)}
		paragraphs = append(paragraphs, NewDialogParagraph("This message based on GtkMessageDialog functionality.\n"+
			"GtkMessageDialog doesn't support setting of image\n"+
			"via SetImage as message type specification since 3.22.\n"+
			"So, that's why there is not warning icon either something similar here.").
			SetJustify(gtk.JUSTIFY_LEFT).SetHorizAlign(gtk.ALIGN_START))
		paragraphs = append(paragraphs, NewDialogParagraph("There appears to be an issue with the configuration of the application.\n"+
			"This issue is not serious, but correcting it will improve your experience.").
			SetJustify(gtk.JUSTIFY_LEFT).SetHorizAlign(gtk.ALIGN_START))

		buttons := []DialogButton{
			{"_OK", gtk.RESPONSE_OK, false, nil},
		}

		response, err := RunMessageDialog(win, title, "", paragraphs, false, buttons,
			func(area *gtk.Box) error {
				lbl, err := gtk.LabelNew("Click the link below for more information, if still not clear:")
				if err != nil {
					return err
				}
				lbl.SetHAlign(gtk.ALIGN_START)
				area.Add(lbl)

				link, err := gtk.LinkButtonNew("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
				if err != nil {
					return err
				}
				area.Add(link)

				cb, err := gtk.CheckButtonNewWithLabel("Do not show this message again")
				if err != nil {
					return err
				}
				area.Add(cb)

				return nil
			})
		if err != nil {
			log.Fatal(err)
		}

		PrintDialogResponse(response)
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create demonstration GTK Dialog action.
// Action trigger is included.
func createDialogAction2(win *gtk.Window) (glib.IAction, error) {
	act, err := glib.SimpleActionNew("DialogAction2", nil)
	if err != nil {
		return nil, err
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		title := "Configuration Issue Detected (demonstration)"
		paragraphs := []*DialogParagraph{NewDialogParagraph("Message dialog based on GTK+ Dialog.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER)}
		paragraphs = append(paragraphs, NewDialogParagraph("This message based on GtkDialog functionality.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER))
		paragraphs = append(paragraphs, NewDialogParagraph("There appears to be an issue with the configuration of the application.\n"+
			"This issue is not serious, but correcting it will improve your experience.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER))
		buttons := []DialogButton{
			{"_OK", gtk.RESPONSE_OK, false, nil},
		}

		response, err := RunDialog(win, gtk.MESSAGE_WARNING, true, title, paragraphs, false, buttons,
			func(area *gtk.Box) error {
				lbl, err := gtk.LabelNew("Click the link below for more information, if still not clear:")
				if err != nil {
					return err
				}
				lbl.SetHAlign(gtk.ALIGN_START)
				area.Add(lbl)

				link, err := gtk.LinkButtonNew("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
				if err != nil {
					return err
				}
				area.Add(link)

				cb, err := gtk.CheckButtonNewWithLabel("Do not show this message again")
				if err != nil {
					return err
				}
				area.Add(cb)

				return nil
			})
		if err != nil {
			log.Fatal(err)
		}

		PrintDialogResponse(response)
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create demonstration GTK MessageDialog action.
// Action trigger is included.
func createDialogAction3(win *gtk.Window) (glib.IAction, error) {
	act, err := glib.SimpleActionNew("DialogAction3", nil)
	if err != nil {
		return nil, err
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		title := "<span weight='bold' size='larger'>Configuration Issue Detected (demonstration)</span>"
		paragraphs := []*DialogParagraph{NewDialogParagraph("Message dialog based on GTK+ MessageDialog.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER)}
		paragraphs = append(paragraphs, NewDialogParagraph("This message based on GtkMessageDialog functionality.\n"+
			"GtkMessageDialog doesn't support setting of image\n"+
			"via SetImage as message type specification since 3.22.\n"+
			"So, that's why there is not warning icon either something similar here.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER))
		paragraphs = append(paragraphs, NewDialogParagraph("There appears to be an issue with the configuration of the application.\n"+
			"This issue is not serious, but correcting it will improve your experience.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER))

		buttons := []DialogButton{
			{"_Yes", gtk.RESPONSE_YES, false, func(btn *gtk.Button) error {
				style, err := btn.GetStyleContext()
				if err != nil {
					return err
				}
				style.AddClass("suggested-action")
				//style.AddClass("destructive-action")
				return nil
			}},
			{"_No", gtk.RESPONSE_NO, true, nil},
		}

		response, err := RunMessageDialog(win, title, "", paragraphs, true, buttons,
			func(area *gtk.Box) error {
				lbl, err := gtk.LabelNew("Click the link below for more information, if still not clear:")
				if err != nil {
					return err
				}
				lbl.SetHAlign(gtk.ALIGN_START)
				area.Add(lbl)

				link, err := gtk.LinkButtonNew("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
				if err != nil {
					return err
				}
				area.Add(link)

				return nil
			})
		if err != nil {
			log.Fatal(err)
		}

		PrintDialogResponse(response)
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create demonstration GTK Dialog action.
// Action trigger is included.
func createDialogAction4(win *gtk.Window) (glib.IAction, error) {
	act, err := glib.SimpleActionNew("DialogAction4", nil)
	if err != nil {
		return nil, err
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		title := "Choose option (demonstration)"
		paragraphs := []*DialogParagraph{NewDialogParagraph("Press Yes to start processing.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER)}
		paragraphs = append(paragraphs, NewDialogParagraph("Note: processing might takes significant amount of time.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER))

		buttons := []DialogButton{
			{"_Yes", gtk.RESPONSE_YES, true, func(btn *gtk.Button) error {
				/*
					style, err := btn.GetStyleContext()
					if err != nil {
						return err
					}
					//style.AddClass("suggested-action")
					style.RemoveClass("suggested-action")

					style.AddClass("destructive-action")
				*/
				return nil
			}},
			{"_No", gtk.RESPONSE_NO, false, nil},
		}

		response, err := RunDialog(win, gtk.MESSAGE_QUESTION, true, title, paragraphs, true, buttons,
			nil)
		if err != nil {
			log.Fatal(err)
		}

		PrintDialogResponse(response)
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create demonstration GTK Dialog action.
// Action trigger is included.
func createDialogAction5(win *gtk.Window) (glib.IAction, error) {
	act, err := glib.SimpleActionNew("DialogAction5", nil)
	if err != nil {
		return nil, err
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		title := "Choose option (demonstration)"
		paragraphs := []*DialogParagraph{NewDialogParagraph("Press Yes to start processing.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER)}
		paragraphs = append(paragraphs, NewDialogParagraph("Note: processing might takes significant amount of time.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER))

		buttons := []DialogButton{
			{"_Yes", gtk.RESPONSE_YES, false, func(btn *gtk.Button) error {
				style, err := btn.GetStyleContext()
				if err != nil {
					return err
				}
				//style.AddClass("suggested-action")
				style.AddClass("destructive-action")
				return nil
			}},
			{"_No", gtk.RESPONSE_NO, true, nil},
		}

		response, err := RunDialog(win, gtk.MESSAGE_QUESTION, false, title, paragraphs, true, buttons,
			nil)
		if err != nil {
			log.Fatal(err)
		}

		PrintDialogResponse(response)
	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Create sophisticated multi-page preference dialog
// with save/restore functionallity to/from the GLib Setting object.
// Action activation require to have GLib Setting Schema
// prliminary installed, otherwise will not work raising message.
// Installation bash script from app folder must be performed in advance.
func createPreferenceAction(win *gtk.Window) (glib.IAction, error) {
	act, err := glib.SimpleActionNew("PreferenceAction", nil)
	if err != nil {
		return nil, err
	}

	_, err = act.Connect("activate", func(action *glib.SimpleAction, param *glib.Variant) {
		name, state, err := GetActionNameAndState(action)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(spew.Sprintf("%v action activated with current state %v and args %v",
			name, state, param))

		app, err := win.GetApplication()
		if err != nil {
			log.Fatal(err)
		}

		found, err := checkSchemaSettingsIsInstalled(app)
		if err != nil {
			log.Fatal(err)
		}

		if found {

			win, err := createPreferenceDialog(app)
			if err != nil {
				log.Fatal(err)
			}

			win.ShowAll()
			win.Show()

			_, err = win.Connect("destroy", func(window *gtk.ApplicationWindow) {
				window.Destroy()
				log.Println("Destroy window")
			})
			if err != nil {
				log.Fatal(err)
			}
		}

	})
	if err != nil {
		return nil, err
	}

	return act, nil
}

// Automatically hide/unhide panel used when fullscreen mode is on.
func createRevealer(data *FullscreenGlobalData) (*gtk.Revealer, error) {
	rev, err := gtk.RevealerNew()
	if err != nil {
		return nil, err
	}
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	rev.Add(box)

	menu, err := createMenuModelForPopover()
	if err != nil {
		return nil, err
	}

	btn, err := SetupMenuButtonWithThemedImage("open-menu-symbolic")
	if err != nil {
		return nil, err
	}
	data.RightTopMenuButton = btn

	btn.SetUsePopover(true)
	btn.SetMenuModel(menu)

	btn2, err := gtk.SeparatorToolItemNew()
	if err != nil {
		return nil, err
	}

	btn3, err := SetupToolButton("view-restore-symbolic", "Leave ssdfsdf")
	if err != nil {
		return nil, err
	}
	//	btn2.SetTitle("asdasdas")
	btn3.SetActionName("win.FullscreenAction")

	hdr, err := SetupHeader(APP_TITLE, "(fullscreen mode)", false)
	if err != nil {
		log.Fatal(err)
	}

	hdr.PackEnd(btn3)
	hdr.PackEnd(btn2)
	hdr.PackEnd(btn)

	box.PackStart(hdr, false, false, 0)

	return rev, nil

}

// Define timer timeouts to hide fullscreen header in case of inactivity.
func manageFullscreenHeaderHideTimer(data *FullscreenGlobalData, dur time.Duration, start bool) {
	if start {
		if data.Timer != nil {
			if !data.Timer.Stop() {
			}
			data.Timer.Reset(dur)
		} else {
			data.Timer = time.AfterFunc(dur, func() {
				showFullscreenHeader(data, false)
			})
		}
	} else {
		if data.Timer != nil {
			if !data.Timer.Stop() {
			}
		}
	}
}

// Show/hide fullscreen header.
func showFullscreenHeader(data *FullscreenGlobalData, show bool) {
	if show {
		data.InFullscreenEventBox = true

		v := data.FullscreeAction.GetState()
		if v != nil && v.GetBoolean() {
			data.Revealer.SetRevealChild(true)
		}
	} else {
		data.InFullscreenEventBox = false
		v := data.FullscreeAction.GetState()
		if v != nil && v.GetBoolean() && !data.RightTopMenuButton.GetActive() {
			data.Revealer.SetRevealChild(false)
		}
	}
}

// Create EventBox widget to catch mouse move events over hidden fullscreen panel.
func createEventBox(data *FullscreenGlobalData) (*gtk.EventBox, error) {
	eb, err := gtk.EventBoxNew()
	if err != nil {
		return nil, err
	}
	dur := time.Millisecond * 1500
	_, err = eb.Connect("enter-notify-event", func() {
		log.Println("EventBox enter notify")

		manageFullscreenHeaderHideTimer(data, dur, false)
		showFullscreenHeader(data, true)
	})
	if err != nil {
		return nil, err
	}
	_, err = eb.Connect("leave-notify-event", func() {
		_, err2 := glib.IdleAdd(func() {
			log.Println("EventBox leave notify")

			manageFullscreenHeaderHideTimer(data, dur, true)
		})
		if err2 != nil {
			log.Fatal(err2)
		}
	})
	if err != nil {
		return nil, err
	}

	return eb, nil
}

// Keeps here extra data for each page of multi-page preference dialog.
type PreferenceRow struct {
	Name  string
	Title string
}

func PreferenceRowNew(name, title string) (*PreferenceRow, *gtk.ListBoxRow, error) {
	pr := &PreferenceRow{Name: name, Title: title}
	lbl, err := gtk.LabelNew(name)
	if err != nil {
		return nil, nil, err
	}

	lbr, err := gtk.ListBoxRowNew()
	if err != nil {
		return nil, nil, err
	}
	lbl.SetHAlign(gtk.ALIGN_START)
	SetAllMargins(lbl, 6)
	lbr.Add(lbl)

	return pr, lbr, nil
}

// Create preference dialog with "Global" page, where controls
// being bound to GLib Setting object to save/restore functionality.
func GlobalPreferencesNew(gsSettings *glib.Settings) (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		return nil, err
	}

	SetAllMargins(box, 18)

	bh := BindingHelperNew(gsSettings)

	lblBehavior, err := gtk.LabelNew(fmt.Sprintf("<b>%s</b>", "Behavior"))
	if err != nil {
		return nil, err
	}
	lblBehavior.SetUseMarkup(true)
	lblBehavior.SetHAlign(gtk.ALIGN_START)
	box.Add(lblBehavior)

	// Show about dialog on application startup
	cbAboutInfo, err := gtk.CheckButtonNewWithLabel("Do not show about information on app startup")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_DONT_SHOW_ABOUT_ON_STARTUP_KEY, cbAboutInfo, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbAboutInfo)

	//Prompt on new session
	cbPrompt, err := gtk.CheckButtonNewWithLabel("Prompt when creating a new session")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_PROMPT_ON_NEW_SESSION_KEY, cbPrompt, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbPrompt)

	//Focus follows the mouse
	cbFocusMouse, err := gtk.CheckButtonNewWithLabel("Focus a terminal when the mouse moves over it")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_TERMINAL_FOCUS_FOLLOWS_MOUSE_KEY, cbFocusMouse, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbFocusMouse)

	//Auto hide the mouse
	cbAutoHideMouse, err := gtk.CheckButtonNewWithLabel("Autohide the mouse pointer when typing")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_AUTO_HIDE_MOUSE_KEY, cbAutoHideMouse, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbAutoHideMouse)

	//middle click closes the terminal
	cbMiddleClickClose, err := gtk.CheckButtonNewWithLabel("Close terminal by clicking middle mouse button on title")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_MIDDLE_CLICK_CLOSE_KEY, cbMiddleClickClose, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbMiddleClickClose)

	//zoom in/out terminal with scroll wheel
	cbControlScrollZoom, err := gtk.CheckButtonNewWithLabel("Zoom the terminal using <Control> and scroll wheel")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_CONTROL_SCROLL_ZOOM_KEY, cbControlScrollZoom, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbControlScrollZoom)

	//require control modifier when clicking title
	cbControlClickTitle, err := gtk.CheckButtonNewWithLabel("Require the <Control> modifier to edit title on click")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_CONTROL_CLICK_TITLE_KEY, cbControlClickTitle, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbControlClickTitle)

	//Closing of last session closes window
	cbCloseWithLastSession, err := gtk.CheckButtonNewWithLabel("Close window when last session is closed")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_CLOSE_WITH_LAST_SESSION_KEY, cbCloseWithLastSession, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbCloseWithLastSession)

	cbNewWindowInheritState, err := gtk.CheckButtonNewWithLabel("New window inherits directory and profile from active terminal")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_INHERIT_WINDOW_STATE_KEY, cbNewWindowInheritState, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbNewWindowInheritState)

	// Save window state (maximized, minimized, fullscreen) between invocations
	cbWindowSaveState, err := gtk.CheckButtonNewWithLabel("Save and restore window state")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_WINDOW_SAVE_STATE_KEY, cbWindowSaveState, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbWindowSaveState)

	// *********** Clipboard Options
	lblClipboard, err := gtk.LabelNew(fmt.Sprintf("<b>%s</b>", "Clipboard"))
	if err != nil {
		return nil, err
	}
	lblClipboard.SetUseMarkup(true)
	lblClipboard.SetHAlign(gtk.ALIGN_START)
	box.Add(lblClipboard)

	//Advacned paste is default
	cbAdvDefault, err := gtk.CheckButtonNewWithLabel("Always use advanced paste dialog")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_PASTE_ADVANCED_DEFAULT_KEY, cbAdvDefault, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbAdvDefault)

	//Unsafe Paste Warning
	cbUnsafe, err := gtk.CheckButtonNewWithLabel("Warn when attempting unsafe paste")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_UNSAFE_PASTE_ALERT_KEY, cbUnsafe, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbUnsafe)

	//Strip Paste
	cbStrip, err := gtk.CheckButtonNewWithLabel("Strip first character of paste if comment or variable declaration")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_STRIP_FIRST_COMMENT_CHAR_ON_PASTE_KEY, cbStrip, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbStrip)

	//Copy on Select
	cbCopyOnSelect, err := gtk.CheckButtonNewWithLabel("Automatically copy text to clipboard when selecting")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_COPY_ON_SELECT_KEY, cbCopyOnSelect, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbCopyOnSelect)

	// Disclaimer
	wdg, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	box.Add(wdg)
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintln("******************************************************"))
	buf.WriteString(fmt.Sprintln("Design of this preference page taken from <a href=\"https://github.com/gnunn1/tilix\">Tilix</a> project."))
	buf.WriteString(fmt.Sprintln("Settings here mainly for demonstration purpose and have"))
	buf.WriteString(fmt.Sprintln("minimal impact to application, or none at all."))
	buf.WriteString(fmt.Sprintln("******************************************************"))
	lbl, err := gtk.LabelNew("")
	if err != nil {
		return nil, err
	}
	lbl.SetMarkup(buf.String())
	lbl.SetJustify(gtk.JUSTIFY_CENTER)
	box.PackEnd(lbl, true, true, 0)

	_, err = box.Connect("destroy", func(b *gtk.Box) {
		bh.Unbind()
		log.Println("Destroy box")
	})
	if err != nil {
		return nil, err
	}

	return box, nil
}

// Create preference dialog with "Appearance" page, where controls
// being bound to GLib Setting object to save/restore functionality.
func AppearancePreferencesNew(gsSettings *glib.Settings) (*gtk.Box, error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		return nil, err
	}

	SetAllMargins(box, 18)

	bh := BindingHelperNew(gsSettings)

	grid, err := gtk.GridNew()
	grid.SetColumnSpacing(12)
	grid.SetRowSpacing(6)
	row := 0

	//Window style
	lbl, err := gtk.LabelNew("Window style")
	if err != nil {
		return nil, err
	}
	lbl.SetHAlign(gtk.ALIGN_END)
	grid.Attach(lbl, 0, row, 1, 1)
	bWindowStyle, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	values := []struct{ value, key string }{
		{"Normal", "normal"},
		{"Disable CSD", "disable-csd"},
		{"Disable CSD, hide toolbar", "disable-csd-hide-toolbar"},
		{"Borderless", "borderless"},
	}
	cbWindowStyle, err := CreateNameValueCombo(values)
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_WINDOW_STYLE_KEY, cbWindowStyle, "active-id", glib.SETTINGS_BIND_DEFAULT)
	bWindowStyle.Add(cbWindowStyle)

	lblRestart, err := gtk.LabelNew("Window restart required")
	if err != nil {
		return nil, err
	}
	lblRestart.SetHAlign(gtk.ALIGN_START)
	lblRestart.SetSensitive(false)
	bWindowStyle.Add(lblRestart)

	grid.Attach(bWindowStyle, 1, row, 1, 1)
	row++

	//Render terminal titlebars smaller then default
	lbl, err = gtk.LabelNew("Terminal title style")
	if err != nil {
		return nil, err
	}
	lbl.SetHAlign(gtk.ALIGN_END)
	grid.Attach(lbl, 0, row, 1, 1)
	values = []struct{ value, key string }{
		{"Normal", SETTINGS_TERMINAL_TITLE_STYLE_VALUE_NORMAL},
		{"Small", SETTINGS_TERMINAL_TITLE_STYLE_VALUE_SMALL},
		{"None", SETTINGS_TERMINAL_TITLE_STYLE_VALUE_NONE},
	}
	cbTitleStyle, err := CreateNameValueCombo(values)
	bh.Bind(SETTINGS_TERMINAL_TITLE_STYLE_KEY, cbTitleStyle, "active-id", glib.SETTINGS_BIND_DEFAULT)
	grid.Attach(cbTitleStyle, 1, row, 1, 1)
	row++

	lbl, err = gtk.LabelNew("Tab position")
	if err != nil {
		return nil, err
	}
	lbl.SetHAlign(gtk.ALIGN_END)
	grid.Attach(lbl, 0, row, 1, 1)
	values = []struct{ value, key string }{
		{"Left", "left"},
		{"Right", "right"},
		{"Top", "top"},
		{"Bottom", "bottom"},
	}
	cbTabPosition, err := CreateNameValueCombo(values)
	bh.Bind(SETTINGS_TAB_POSITION_KEY, cbTabPosition, "active-id", glib.SETTINGS_BIND_DEFAULT)
	grid.Attach(cbTabPosition, 1, row, 1, 1)
	row++

	//Dark Theme
	lbl, err = gtk.LabelNew("Theme variant")
	if err != nil {
		return nil, err
	}
	lbl.SetHAlign(gtk.ALIGN_END)
	grid.Attach(lbl, 0, row, 1, 1)
	values = []struct{ value, key string }{
		{"Default", SETTINGS_THEME_VARIANT_SYSTEM_VALUE},
		{"Light", SETTINGS_THEME_VARIANT_LIGHT_VALUE},
		{"Dark", SETTINGS_THEME_VARIANT_DARK_VALUE},
	}
	cbThemeVariant, err := CreateNameValueCombo(values)
	bh.Bind(SETTINGS_THEME_VARIANT_KEY, cbThemeVariant, "active-id", glib.SETTINGS_BIND_DEFAULT)
	grid.Attach(cbThemeVariant, 1, row, 1, 1)
	row++

	//Background Image
	lbl, err = gtk.LabelNew("Background image")
	if err != nil {
		return nil, err
	}
	grid.Attach(lbl, 0, row, 1, 1)

	fcbImage, err := gtk.FileChooserButtonNew("Select Image", gtk.FILE_CHOOSER_ACTION_OPEN)
	if err != nil {
		return nil, err
	}
	fcbImage.SetHExpand(true)
	ff, err := gtk.FileFilterNew()
	if err != nil {
		return nil, err
	}
	ff.SetName("All Image Files")
	ff.AddMimeType("image/jpeg")
	ff.AddMimeType("image/png")
	ff.AddMimeType("image/bmp")
	fcbImage.AddFilter(ff)
	ff, err = gtk.FileFilterNew()
	if err != nil {
		return nil, err
	}
	ff.AddPattern("*")
	ff.SetName("All Files")
	fcbImage.AddFilter(ff)
	filename := gsSettings.GetString(SETTINGS_BACKGROUND_IMAGE_KEY)
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		// log.Println(spew.Sprintf("File %q found", filename))
		fcbImage.SetFilename(filename)
	}
	_, err = fcbImage.Connect("file-set", func(fcb *gtk.FileChooserButton) {
		selectedFilename := fcb.GetFilename()
		if _, err := os.Stat(selectedFilename); !os.IsNotExist(err) {
			gsSettings.SetString(SETTINGS_BACKGROUND_IMAGE_KEY, selectedFilename)
		}
	})
	if err != nil {
		return nil, err
	}

	btnReset, err := SetupButtonWithThemedImage("edit-delete-symbolic")
	if err != nil {
		return nil, err
	}
	btnReset.SetTooltipText("Reset background image")
	_, err = btnReset.Connect("clicked", func(btn *gtk.Button) {
		fcbImage.UnselectAll()
		gsSettings.Reset(SETTINGS_BACKGROUND_IMAGE_KEY)
	})
	if err != nil {
		return nil, err
	}

	values = []struct{ value, key string }{
		{"Scale", SETTINGS_BACKGROUND_IMAGE_MODE_SCALE_VALUE},
		{"Tile", SETTINGS_BACKGROUND_IMAGE_MODE_TILE_VALUE},
		{"Center", SETTINGS_BACKGROUND_IMAGE_MODE_CENTER_VALUE},
		{"Stretch", SETTINGS_BACKGROUND_IMAGE_MODE_STRETCH_VALUE},
	}
	cbImageMode, err := CreateNameValueCombo(values)
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_BACKGROUND_IMAGE_MODE_KEY, cbImageMode, "active-id", glib.SETTINGS_BIND_DEFAULT)

	// Background image settings only enabled if transparency is enabled
	bh.Bind(SETTINGS_ENABLE_TRANSPARENCY_KEY, fcbImage, "sensitive", glib.SETTINGS_BIND_DEFAULT)
	bh.Bind(SETTINGS_ENABLE_TRANSPARENCY_KEY, btnReset, "sensitive", glib.SETTINGS_BIND_DEFAULT)
	bh.Bind(SETTINGS_ENABLE_TRANSPARENCY_KEY, cbImageMode, "sensitive", glib.SETTINGS_BIND_DEFAULT)

	bChooser, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
	bChooser.Add(fcbImage)
	bChooser.Add(btnReset)

	bImage, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	bImage.Add(bChooser)
	bImage.Add(cbImageMode)
	grid.Attach(bImage, 1, row, 1, 1)
	row++

	box.Add(grid)

	cbRightSidebar, err := gtk.CheckButtonNewWithLabel("Place the sidebar on the right")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_SIDEBAR_RIGHT, cbRightSidebar, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbRightSidebar)

	cbTitleShowWhenSingle, err := gtk.CheckButtonNewWithLabel("Show the terminal title even if it's the only terminal")
	if err != nil {
		return nil, err
	}
	bh.Bind(SETTINGS_TERMINAL_TITLE_SHOW_WHEN_SINGLE_KEY, cbTitleShowWhenSingle, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbTitleShowWhenSingle)

	cbUseTabs, err := gtk.CheckButtonNewWithLabel("Use tabs instead of sidebar (Application restart required)")
	bh.Bind(SETTINGS_USE_TABS_KEY, cbUseTabs, "active", glib.SETTINGS_BIND_DEFAULT)
	box.Add(cbUseTabs)

	// Disclaimer
	wdg, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	box.Add(wdg)
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintln("******************************************************"))
	buf.WriteString(fmt.Sprintln("Design of this preference page taken from <a href=\"https://github.com/gnunn1/tilix\">Tilix</a> project."))
	buf.WriteString(fmt.Sprintln("Settings here mainly for demonstration purpose and have"))
	buf.WriteString(fmt.Sprintln("minimal impact to application, or none at all."))
	buf.WriteString(fmt.Sprintln("******************************************************"))
	lbl, err = gtk.LabelNew("")
	if err != nil {
		return nil, err
	}
	lbl.SetMarkup(buf.String())
	lbl.SetJustify(gtk.JUSTIFY_CENTER)
	box.PackEnd(lbl, true, true, 0)

	_, err = box.Connect("destroy", func(b *gtk.Box) {
		bh.Unbind()
	})
	if err != nil {
		return nil, err
	}

	return box, nil
}

// Verify, that GLib Setting's schema is installed, otherwise return false.
func checkSchemaSettingsIsInstalled(app *gtk.Application) (bool, error) {
	parentWin := app.GetActiveWindow()
	// Verify that GSettingsSchema is installed
	schemaSource := glib.SettingsSchemaSourceGetDefault()
	if schemaSource == nil {
		title := "<span weight='bold' size='larger'>Schema settings configuration error</span>"
		paragraphs := []*DialogParagraph{NewDialogParagraph("No one GTK+ schema settings is found.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER)}
		paragraphs = append(paragraphs, NewDialogParagraph("Please install xml schema and repeat operation.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER))

		buttons := []DialogButton{
			{"_OK", gtk.RESPONSE_OK, false, nil},
		}

		_, err := RunMessageDialog(parentWin, title, "", paragraphs, false, buttons, nil)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	schema := schemaSource.Lookup(SETTINGS_ID, false)
	if schema == nil {
		title := "<span weight='bold' size='larger'>Schema settings configuration error</span>"
		paragraphs := []*DialogParagraph{NewDialogParagraph(fmt.Sprintf("GTK+ schema %q is not found.", SETTINGS_ID)).
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER)}
		paragraphs = append(paragraphs, NewDialogParagraph("Please install xml schema and repeat operation.").
			SetJustify(gtk.JUSTIFY_CENTER).SetHorizAlign(gtk.ALIGN_CENTER))

		buttons := []DialogButton{
			{"_OK", gtk.RESPONSE_OK, false, nil},
		}

		_, err := RunMessageDialog(parentWin, title, "", paragraphs, false, buttons, nil)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// Create sophisticated multi-page preference dialog
// with save/restore functionallity to/from the GLib Setting object.
func createPreferenceDialog(app *gtk.Application) (*gtk.ApplicationWindow, error) {
	parentWin := app.GetActiveWindow()
	win, err := gtk.ApplicationWindowNew(app)
	if err != nil {
		return nil, err
	}

	// Settings
	win.SetTitle("Preferences")
	win.SetTransientFor(parentWin)
	win.SetDestroyWithParent(true)
	win.SetShowMenubar(false)
	gsSettings, err := glib.SettingsNew(SETTINGS_ID)
	if err != nil {
		return nil, err
	}

	// Create window header
	hbMain, err := SetupHeader("", "", true)
	if err != nil {
		return nil, err
	}
	hbMain.SetHExpand(true)

	hbSide, err := SetupHeader("Preferences", "", false)
	if err != nil {
		return nil, err
	}
	hbSide.SetHExpand(false)

	bTitle, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	bTitle.Add(hbSide)
	sTitle, err := gtk.SeparatorNew(gtk.ORIENTATION_VERTICAL)
	bTitle.Add(sTitle)
	bTitle.Add(hbMain)

	win.SetTitlebar(bTitle)

	var list = make(map[uintptr]interface{})

	// Create Stack and boxes
	pages, err := gtk.StackNew()
	if err != nil {
		return nil, err
	}
	pages.SetHExpand(true)
	pages.SetVExpand(true)

	// Create ListBox
	lbSide, err := gtk.ListBoxNew()
	if err != nil {
		return nil, err
	}
	lbSide.SetCanFocus(true)
	lbSide.SetSelectionMode(gtk.SELECTION_BROWSE)
	lbSide.SetVExpand(true)
	_, err = lbSide.Connect("row-selected", func(lb *gtk.ListBox, row *gtk.ListBoxRow) {
		r := list[row.Native()]
		if r1, ok := r.(*PreferenceRow); ok {
			//log.Println(spew.Sprintf("%+v", r1))
			pages.SetVisibleChildName(r1.Name)
			hbMain.SetTitle(r1.Title)
		}
	})
	if err != nil {
		return nil, err
	}

	gp, err := GlobalPreferencesNew(gsSettings)
	if err != nil {
		return nil, err
	}
	pages.AddTitled(gp, "Global", "Global")
	pr, row, err := PreferenceRowNew("Global", "Global")
	if err != nil {
		return nil, err
	}
	list[row.Native()] = pr
	lbSide.Add(row)

	ap, err := AppearancePreferencesNew(gsSettings)
	if err != nil {
		return nil, err
	}
	pages.AddTitled(ap, "Appearance", "Appearance")
	pr, row, err = PreferenceRowNew("Appearance", "Appearance")
	if err != nil {
		return nil, err
	}
	list[row.Native()] = pr
	lbSide.Add(row)

	sw, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	sw.Add(lbSide)
	sw.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	sw.SetShadowType(gtk.SHADOW_NONE)
	sw.SetSizeRequest(220, -1)

	bSide, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	bSide.Add(sw)

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	box.Add(bSide)
	sep, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	box.Add(sep)
	box.Add(pages)

	win.Add(box)

	sgSide, err := gtk.SizeGroupNew(gtk.SIZE_GROUP_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	sgSide.AddWidget(hbSide)
	sgSide.AddWidget(bSide)

	sgMain, err := gtk.SizeGroupNew(gtk.SIZE_GROUP_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	sgMain.AddWidget(hbMain)
	sgMain.AddWidget(pages)

	// Set initial title
	hbMain.SetTitle("Global")

	return win, nil
}

func main() {

	data := &FullscreenGlobalData{}

	gtk.Init(nil)
	app, err := gtk.ApplicationNew(APP_ID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal(err)
	}

	_, err = app.Application.Connect("activate", func(application *gtk.Application) {
		win, err := gtk.ApplicationWindowNew(application)
		if err != nil {
			log.Fatal(err)
		}
		win.SetTitle("Example")
		win.SetDefaultSize(900, 600)

		_, err2 := win.Connect("destroy", func(window *gtk.ApplicationWindow) {
			application, err := window.GetApplication()
			if err != nil {
				log.Fatal(err)
			}
			application.Quit()
		})
		if err2 != nil {
			log.Fatal(err2)
		}

		var act glib.IAction

		act, err = createQuitAction(&win.Window)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		gsSettings, err := glib.SettingsNew(SETTINGS_ID)
		if err != nil {
			log.Fatal(err)
		}
		act, err = createAboutAction(&win.Window, gsSettings)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		act, err = createDialogAction1(&win.Window)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		act, err = createDialogAction2(&win.Window)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		act, err = createDialogAction3(&win.Window)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		act, err = createDialogAction4(&win.Window)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		act, err = createDialogAction5(&win.Window)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		act, err = createCheckBoxAction()
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		act, err = createChooseAction()
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		act, err = createPreferenceAction(&win.Window)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		menu, err := createMenuModelForPopover()
		if err != nil {
			log.Fatal(err)
		}
		btn, err := SetupMenuButtonWithThemedImage("open-menu-symbolic")
		if err != nil {
			log.Fatal(err)
		}
		btn.SetUsePopover(true)
		btn.SetMenuModel(menu)

		hdr, err := SetupHeader(APP_TITLE, "", true)
		if err != nil {
			log.Fatal(err)
		}
		hdr.PackEnd(btn)

		win.SetTitlebar(hdr)

		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		if err != nil {
			log.Fatal(err)
		}

		_, err = box.Connect("destroy", func() {
			log.Println("Destroy box")
		})
		if err != nil {
			log.Fatal(err)
		}

		eb, err := createEventBox(data)
		if err != nil {
			log.Fatal(err)
		}
		data.EventBox = eb
		rev, err := createRevealer(data)
		if err != nil {
			log.Fatal(err)
		}
		data.Revealer = rev
		eb.Add(rev)
		eb.SetSizeRequest(-1, 1)
		eb.Hide()

		box.PackStart(eb, false, false, 0)

		menu2, err := createMenuModelForMenuBar()
		if err != nil {
			log.Fatal(err)
		}

		menuBar, err := gtk.MenuBarFromModelNew(menu2)
		if err != nil {
			log.Fatal(err)
		}

		box.PackStart(menuBar, false, false, 0)

		tbx, err := createToolbar()
		if err != nil {
			log.Fatal(err)
		}

		box.PackStart(tbx, false, false, 0)

		div, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		if err != nil {
			log.Fatal(err)
		}
		box.PackStart(div, false, false, 0)

		sw, err := gtk.ScrolledWindowNew(nil, nil)
		// sw, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		// sw, err := gtk.FrameNew("sdfsdfsdf")
		// sw, err := gtk.ButtonNew()
		if err != nil {
			log.Fatal(err)
		}

		vp, err := gtk.ViewportNew(nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		box2, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		if err != nil {
			log.Fatal(err)
		}
		_, err = box2.Connect("destroy", func() {
			log.Println("Destroy box2")
		})
		if err != nil {
			log.Fatal(err)
		}

		vp.Add(box2)
		sw.Add(vp)
		_, err = sw.Connect("destroy", func() {
			log.Println("Destroy sw")
			// chd, err := sw.GetChild()
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// log.Println(spew.Sprintf("%+v", chd))
			// log.Println(spew.Sprintf("%+v", box2))
			// sw.Remove(chd)
			// chd.Unref()
			// chd.Destroy()

			// parent, err := chd.GetParent()
			// if err != nil {
			// log.Fatal(err)
			// }

		})
		if err != nil {
			log.Fatal(err)
		}

		box.Add(sw)

		win.Add(box)

		act, err = createFullscreenAction(&win.Window, data)
		if err != nil {
			log.Fatal(err)
		}
		win.AddAction(act)

		win.ShowAll()

		// Run code, when app message queue becomes empty.
		if !gsSettings.GetBoolean(SETTINGS_DONT_SHOW_ABOUT_ON_STARTUP_KEY) {
			_, err2 := glib.IdleAdd(func() {
				action := win.LookupAction("AboutAction")
				if action != nil {
					action.Activate(nil)
				}
			})
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	app.Run([]string{})
}
