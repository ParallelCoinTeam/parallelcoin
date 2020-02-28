package controller

type DuoUIpage struct {
	Name    string
	Command func(interface{})
	Data    interface{}
}

type DuoUIpages *map[string]*DuoUIpage
