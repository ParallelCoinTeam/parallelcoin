package duoui

import (
	"errors"
	"fmt"
	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/pages"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/clipboard"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/gui/gelook/ico"
	"github.com/p9c/pod/pkg/util/interrupt"
	"image"
	"image/color"
	"sync"
)

type DuoUI struct {
	ly *model.DuoUI
	rc *rcd.RcVar
}

var (
	clipboardStarted bool
	clipboardMu      sync.Mutex
	logoButton       = new(gel.Button)
	headerList       = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
	buttonToastOK = new(gel.Button)
	listToasts    = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: false,
		Alignment:   0,
		Position:    layout.Position{},
	}
	passPhrase        string
	confirmPassPhrase string
	passEditor        = &gel.Editor{
		SingleLine: true,
		// Submit:     true,
	}
	confirmPassEditor = &gel.Editor{
		SingleLine: true,
		// Submit:     true,
	}
	listWallet = &layout.List{
		Axis: layout.Vertical,
	}
	encryption         = new(gel.CheckBox)
	seed               = new(gel.CheckBox)
	testnet            = new(gel.CheckBox)
	buttonCreateWallet = new(gel.Button)
)

func (ui *DuoUI) DuoUIbody() func() {
	return func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(ui.ly.Context,
			layout.Rigid(ui.DuoUIsidebar()),
			layout.Flexed(1, ui.DuoUIcontent()),
		)
	}
}

func (ui *DuoUI) DuoUIcontent() func() {
	return func() {
		ui.rc.CurrentPage.Layout(ui.ly.Context)
	}
}

func (ui *DuoUI) DuoUIfooter() func() {
	ctx := ui.ly.Context
	th := ui.ly.Theme
	return func() {
		footer := th.DuoUIcontainer(0, th.Colors["Dark"])
		footer.FullWidth = true
		footer.Layout(ctx, layout.N, func() {
			layout.Flex{Spacing: layout.SpaceBetween}.Layout(ctx,
				layout.Rigid(
					component.FooterLeftMenu(ui.rc, ctx, th, ui.ly.Pages)),
				layout.Flexed(1, func() {}),
				layout.Rigid(
					component.FooterRightMenu(ui.rc, ctx, th, ui.ly.Pages)),
			)
		})
	}
}

func (ui *DuoUI) DuoUIheader() func() {
	th := ui.ly.Theme
	ctx := ui.ly.Context
	return func() {
		iSize := 32
		iWidth := 48
		iHeight := 48
		iPadV := 3
		iPadH := 3
		if ui.ly.Viewport > 740 {
			iSize = 64
			iWidth = 96
			iHeight = 96
			iPadV = 6
			iPadH = 6
		}
		th.DuoUIcontainer(0, th.Colors["Dark"]).Layout(ctx, layout.NW, func() {
			layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(ctx,
				layout.Rigid(func() {
					var logoMeniItem gelook.DuoUIbutton
					logoMeniItem = th.DuoUIbutton("", "", "",
						th.Colors["Dark"], "", "", "logo",
						th.Colors["Light"], 0, iSize, iWidth,
						iHeight, iPadV, iPadH, iPadV, iPadH)
					for logoButton.Clicked(ctx) {
						th.ChangeLightDark()
					}
					logoMeniItem.IconLayout(ctx, logoButton)
				}),
				layout.Flexed(1,
					component.HeaderMenu(ui.rc, ctx, th, ui.ly.Pages),
				),
				layout.Rigid(
					component.Label(ctx, th, th.Fonts["Primary"], 12,
						th.Colors["Light"],
						ui.rc.Status.Wallet.Balance.Load()+" "+ui.rc.
							Settings.Abbrevation),
				),
				layout.Rigid(
					component.Label(ctx, th, th.Fonts["Primary"], 12,
						th.Colors["Light"], fmt.Sprint(ui.ly.Viewport)),
				),
			)
		})
	}
}

func (ui *DuoUI) DuoUIloaderCreateWallet() {
	cs := ui.ly.Context.Constraints
	th := ui.ly.Theme
	ctx := ui.ly.Context
	gelook.DuoUIdrawRectangle(ctx, cs.Width.Max, cs.Height.Max,
		th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Center.Layout(ctx, func() {
		controllers := []func(){
			func() {
				bal := th.H5("Enter the private passphrase for your new" +
					" wallet:")
				bal.Font.Typeface = th.Fonts["Primary"]
				bal.Color = th.Colors["Dark"]
				bal.Layout(ctx)
			},
			func() {
				layout.UniformInset(unit.Dp(8)).Layout(ctx, func() {
					e := th.DuoUIeditor("Enter Passphrase", "Dark", "Light", 32)
					e.Font.Typeface = th.Fonts["Primary"]
					e.Font.Style = text.Regular
					e.Layout(ctx, passEditor)
					for _, e := range passEditor.Events(ctx) {
						switch e.(type) {
						case gel.ChangeEvent:
							passPhrase = passEditor.Text()
						}
					}
				})
			},
			func() {
				layout.UniformInset(unit.Dp(8)).Layout(ctx, func() {
					e := th.DuoUIeditor("Repeat Passphrase", "Dark", "Light", 32)
					e.Font.Typeface = th.Fonts["Primary"]
					e.Font.Style = text.Regular
					e.Layout(ctx, confirmPassEditor)
					for _, e := range confirmPassEditor.Events(ctx) {
						switch e.(type) {
						case gel.ChangeEvent:
							confirmPassPhrase = confirmPassEditor.Text()
						}
					}
				})
			},
			func() {
				encryptionCheckBox := th.DuoUIcheckBox(
					"Do you want to add an additional layer of encryption"+
						" for public data?", th.Colors["Dark"],
					th.Colors["Dark"])
				encryptionCheckBox.Font.Typeface = th.Fonts["Primary"]
				encryptionCheckBox.Color = gelook.HexARGB(th.Colors["Dark"])
				encryptionCheckBox.Layout(ctx, encryption)
			},
			func() {
				// TODO: needs input box for seed
				seedCheckBox := th.DuoUIcheckBox(
					"Do you have an existing wallet seed you want to use?",
					th.Colors["Dark"], th.Colors["Dark"])
				seedCheckBox.Font.Typeface = th.Fonts["Primary"]
				seedCheckBox.Color = gelook.HexARGB(th.Colors["Dark"])
				seedCheckBox.Layout(ctx, seed)
			},
			func() {
				testnetCheckBox := th.DuoUIcheckBox(
					"Use testnet?", th.Colors["Dark"], th.Colors["Dark"])
				testnetCheckBox.Font.Typeface = th.Fonts["Primary"]
				testnetCheckBox.Color = gelook.HexARGB(th.Colors["Dark"])
				testnetCheckBox.Layout(ctx, testnet)
			},
			func() {
				var createWalletbuttonComp gelook.DuoUIbutton
				createWalletbuttonComp = th.DuoUIbutton(th.
					Fonts["Secondary"], "CREATE WALLET", th.Colors["Dark"],
					th.Colors["Light"], th.Colors["Light"],
					th.Colors["Dark"], "", th.Colors["Dark"], 16, 0,
					125, 32, 4, 4, 4, 4)
				for buttonCreateWallet.Clicked(ctx) {
					if passPhrase != "" && passPhrase == confirmPassPhrase {
						if testnet.Checked(ctx) {
							ui.rc.UseTestnet()
						}
						ui.rc.CreateWallet(passPhrase, "", "", "")
						if testnet.Checked(ctx) {
							interrupt.RequestRestart()
						}
					}
				}
				createWalletbuttonComp.Layout(ctx, buttonCreateWallet)
			},
		}
		listWallet.Layout(ctx, len(controllers), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(ctx, controllers[i])
		})
	})
}

func (ui *DuoUI) DuoUImainScreen() {
	ctx := ui.ly.Context
	th := ui.ly.Theme
	th.DuoUIcontainer(0, th.Colors["Dark"]).Layout(ctx,
		layout.Center, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(ctx,
				layout.Rigid(ui.DuoUIheader()),
				layout.Flexed(1, ui.DuoUIbody()),
				layout.Rigid(ui.DuoUIfooter()),
			)
		})
}

func (ui *DuoUI) DuoUImenu() func() {
	nav := ui.ly.Navigation
	return func() {
		nav.Width = 48
		nav.Height = 48
		nav.TextSize = 0
		nav.IconSize = 24
		nav.PaddingVertical = 4
		nav.PaddingHorizontal = 0
		if ui.ly.Viewport > 740 {
			nav.Width = 96
			nav.Height = 72
			nav.TextSize = 48
			nav.IconSize = 36
			nav.PaddingVertical = 8
			nav.PaddingHorizontal = 0
		}
		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly}.
			Layout(ui.ly.Context, layout.Rigid(
				component.MainNavigation(ui.rc, ui.ly.Context,
					ui.ly.Theme, ui.ly.Pages, nav)),
			)
	}
}

// splash screen
func (ui *DuoUI) DuoUIsidebar() func() {
	return func() {
		ui.ly.Theme.DuoUIcontainer(0,
			ui.ly.Theme.Colors["Dark"]).Layout(ui.ly.Context,
			layout.NW, func() {
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(ui.ly.Context,
					layout.Rigid(ui.DuoUImenu()),
				)
			})
	}
}

func (ui *DuoUI) DuoUIsplashScreen() {
	ctx := ui.ly.Context
	th := ui.ly.Theme
	th.DuoUIcontainer(0, th.Colors["Dark"]).
		Layout(ctx, layout.Center, func() {
			logo, _ := gelook.NewDuoUIicon(ico.ParallelCoin)
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(ctx,
				layout.Rigid(func() {
					layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(ctx,
						layout.Rigid(func() {
							layout.UniformInset(unit.Dp(8)).
								Layout(ctx, func() {
									size := ctx.Px(unit.Dp(256)) - 2*ctx.
										Px(unit.Dp(8))
									if logo != nil {
										logo.Color = gelook.HexARGB(th.
											Colors["Light"])
										logo.Layout(ctx, unit.Px(float32(size)))
									}
									ctx.Dimensions = layout.Dimensions{
										Size: image.Point{X: size, Y: size},
									}
								})
						}),
						layout.Flexed(1, func() {
							layout.UniformInset(unit.Dp(60)).Layout(ctx, func() {
								txt := th.H1("PLAN NINE FROM FAR, " +
									"FAR AWAY SPACE")
								txt.Font.Typeface = th.Fonts["Secondary"]
								txt.Color = th.Colors["Light"]
								txt.Layout(ctx)
							})
						}),
					)
				}),
				layout.Flexed(1, component.DuoUIlogger(ui.rc, ctx, th)),
			)
		})
}

// Main wallet screen
func (ui *DuoUI) toastButton(text, txtColor, bgColor, txtHoverColor, bgHoverColor, icon, iconColor string, button *gel.Button) func() {
	var b gelook.DuoUIbutton
	return func() {
		layout.Inset{
			Top: unit.Dp(8), Bottom: unit.Dp(8),
			Left: unit.Dp(8), Right: unit.Dp(8),
		}.Layout(ui.ly.Context, func() {
			b = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Fonts["Primary"], text,
				txtColor, bgColor, txtHoverColor, bgHoverColor,
				icon, iconColor, 16, 24, 120, 60,
				0, 0, 0, 0)
			for button.Clicked(ui.ly.Context) {
				// ui.rc.ShowToast = false
			}
			b.Layout(ui.ly.Context, button)
		})
	}
}

func DuOuI(rc *rcd.RcVar) (duo *model.DuoUI, err error) {
	duo = &model.DuoUI{
		Window: app.NewWindow(
			app.Size(unit.Dp(1024), unit.Dp(640)),
			app.Title("ParallelCoin"),
		),
	}
	duo.Context = layout.NewContext(duo.Window.Queue())
	// rc.StartLogger()
	// sys.Components["logger"].View()
	// d.sys.Components["logger"].View
	duo.Navigation = &model.DuoUInav{
		Items: make(map[string]*gelook.DuoUIthemeNav),
	}
	// navigations["mainMenu"] = mainMenu()
	// Icons
	// rc.Settings.Daemon = rcd.GetCoreSettings()
	duo.Theme = gelook.NewDuoUItheme()
	// duo.Pages = components.LoadPages(duo.Context, duo.Theme, rc)
	duo.Pages = &model.DuoUIpages{
		Controller: nil,
		Theme:      pages.LoadPages(rc, duo.Context, duo.Theme),
	}
	component.SetPage(rc, duo.Pages.Theme["OVERVIEW"])
	clipboardMu.Lock()
	if !clipboardStarted {
		clipboardStarted = true
		clipboard.Start()
	}
	clipboardMu.Unlock()
	return
}

func DuoUImainLoop(d *model.DuoUI, r *rcd.RcVar) error {
	Debug("starting up duo ui main loop")
	ui := &DuoUI{ly: d, rc: r}
	ctx := ui.ly.Context
	for {
		select {
		case <-ui.rc.Ready:
			updateTrigger := make(chan struct{}, 1)
			go func() {
			quitTrigger:
				for {
					select {
					case <-updateTrigger:
						Trace("repaint forced")
						ui.ly.Window.Invalidate()
					case <-ui.rc.Quit:
						break quitTrigger
					}
				}
			}()
			go ui.rc.ListenInit(updateTrigger)
			ui.rc.IsReady = true
			r.Boot.IsBoot = false
		case <-ui.rc.Quit:
			Debug("quit signal received")
			if !interrupt.Requested() {
				interrupt.Request()
			}
			// This case is for handling when some external application is
			//controlling the GUI and to gracefully handle the back-end
			//servers being shut down by the interrupt library receiving an
			//interrupt signal  Probably nothing needs to be run between
			//starting it and shutting down
			<-interrupt.HandlersDone
			Debug("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
			// TODO events of gui
		case e := <-ui.rc.Commands.Events:
			switch e := e.(type) {
			case rcd.DuoUIcommandEvent:
				ui.rc.Commands.History = append(ui.rc.Commands.History, e.Command)
				ui.ly.Window.Invalidate()
			}
		case e := <-ui.ly.Window.Events():
			ui.ly.Viewport = ctx.Constraints.Width.Max
			switch e := e.(type) {
			case system.DestroyEvent:
				Debug("destroy event received")
				interrupt.Request()
				// Here do cleanup like are you sure (
				//optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				ctx.Reset(e.Config, e.Size)
				if ui.rc.Boot.IsBoot {
					if ui.rc.Boot.IsFirstRun {
						ui.DuoUIloaderCreateWallet()
					} else {
						ui.DuoUIsplashScreen()
					}
					e.Frame(ctx.Ops)
				} else {
					ui.DuoUImainScreen()
					if ui.rc.Dialog.Show {
						component.DuoUIdialog(ui.rc, ctx, ui.ly.Theme)
						// ui.DuoUItoastSys()
					}
					e.Frame(ctx.Ops)
				}
				//ui.ly.Window.Invalidate()
			}
		}
	}
}

func renderIcon(gtx *layout.Context, icon *gelook.DuoUIicon) func() {
	return func() {
		icon.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0x55, B: 0x30}
		icon.Layout(gtx, unit.Dp(float32(48)))
		pointer.Rect(image.Rectangle{Max: image.Point{
			X: 64,
			Y: 64,
		}}).Add(gtx.Ops)
	}
}
