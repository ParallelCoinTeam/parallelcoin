package sap

import (
	"fmt"
)

//
//func NewSap(d mod.Dap) mod.Sap{
//	return mod.Sap{
//		UI: d.UI,
//	}
//}

//func (s *sap) Sap()*mod.Sap {
//	return &mod.Sap{
//		//Nav:   g.getMenuItems(ui),
//	}
//}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}
