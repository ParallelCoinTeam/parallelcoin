Feature-rich GTK+3 app demo written in go
=================================================

Aim of this application is to interconnect GLIB/GTK+ components and widgets together
with use of helpfull code patterns, meet the real needs of developers.

It's obligatory to have GTK+ 3.12 and high installed, otherwise app will not compile
(GtkPopover, GtkHeaderBar require more recent GTK+3 installation).

Features of this demonstration:
1) All actions in application implemented via GLIB's GAction component working as entry point
for any activity in application. Advanced use of GAction utilize
"[action with states and parameters](https://developer.gnome.org/GAction/)"
with seamless integration of such actions to UI menus.
2) New widgets such as GtkHeaderBar, GtkPopover and others are used.
3) Good example of preference dialog with save/restore functionality.
4) Helpfull code pattern are present: dialogs and message boxes demo
working with save/restore settings (via GSettings),
fullscreen mode on/off, actions with states or stateless and others.

Screenshots
-----------

Main form with about dialog:
![image](https://raw.github.com/d2r2/gotk3/master/examples/cool_app/docs/cool_app_screen_1.png)

Main form and modern popover menu:
![image](https://raw.github.com/d2r2/gotk3/master/examples/cool_app/docs/cool_app_screen_2.png)

Main form and classic main menu:
![image](https://raw.github.com/d2r2/gotk3/master/examples/cool_app/docs/cool_app_screen_3.png)

One of few dialogs demo:
![image](https://raw.github.com/d2r2/gotk3/master/examples/cool_app/docs/cool_app_screen_4.png)

Preference dialog with save/restore settings functionality:
![image](https://raw.github.com/d2r2/gotk3/master/examples/cool_app/docs/cool_app_screen_5.png)


Installation
------------

Almost no action needed, the main requirement is to have the GOTK3 library preliminary installed.
Still, to make a "preference dialog" function properly scripts `install_schema.sh`/`uninstall_schema.sh`
must be used, to install/compile [GLIB setting's schema](https://developer.gnome.org/GSettings/).

To run application type in console `go run ./cool_app.go`.

Additional recommendations
-------------------------
- Use GNOME application `gtk3-icon-browser` to find themed icons available at your linux desktop.

Contact
-------

Please use [Github issue tracker](https://github.com/d2r2/gotk3/issues)
for filing bugs or feature requests.
