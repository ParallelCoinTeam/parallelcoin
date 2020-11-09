package dap

func (d *dap) StartServices() (err error) {
	Debug("starting up services")
	// Start Node
	Debug("starting up SERVICES")

	// err = d.NodeService()
	// if err != nil {
	//	Error(err)
	// }
	// d.boot.Rc.Cx.RPCServer = <-d.boot.Rc.Cx.NodeChan
	// d.boot.Rc.Cx.Node.Store(true)

	// Start wallet
	// err = d.WalletService()
	// if err != nil {
	//	Error(err)
	// }
	// d.boot.Rc.Cx.WalletServer = <-d.boot.Rc.Cx.WalletChan
	// d.boot.Rc.Cx.Wallet.Store(true)
	// d.boot.Rc.Cx.WalletServer.Rescan(nil, nil)
	// d.boot.Rc.Ready <- struct{}{}
	return
}

// func (d *dap) WalletService() error {
//	d.boot.Rc.Cx.WalletKill = make(chan struct{})
//	d.boot.Rc.Cx.Wallet.Store(false)
//	var err error
//	if !*d.boot.Rc.Cx.Config.WalletOff {
//		go func() {
//			Info("starting wallet")
//			//utils.GetBiosMessage(view, "starting wallet")
//			err = walletmain.Main(d.boot.Rc.Cx)
//			if err != nil {
//				fmt.Println("error running wallet:", err)
//				os.Exit(1)
//			}
//		}()
//	}
//	interrupt.AddHandler(func() {
//		Warn("interrupt received, " +
//			"shutting down shell modules")
//		close(d.boot.Rc.Cx.WalletKill)
//	})
//	return err
// }

// func (d *dap) NodeService() error {
//	d.boot.Rc.Cx.NodeKill = make(chan struct{})
//	d.boot.Rc.Cx.Node.Store(false)
//	var err error
//	if !*d.boot.Rc.Cx.Config.NodeOff {
//		go func() {
//			Info(d.boot.Rc.Cx.Language.RenderText("goApp_STARTINGNODE"))
//			//utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))
//			err = node.Main(d.boot.Rc.Cx, nil)
//			if err != nil {
//				Info("error running node:", err)
//				os.Exit(1)
//			}
//		}()
//	}
//	interrupt.AddHandler(func() {
//		Warn("interrupt received, " +
//			"shutting down node")
//		close(d.boot.Rc.Cx.NodeKill)
//	})
//	return err
// }
