package gel

type DuoUIpage struct {
	Name    string
	Command func()
	Data    interface{}
}

type DuoUIpages *map[string]*DuoUIpage
