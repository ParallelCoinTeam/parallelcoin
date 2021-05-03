// +build !linux,!freebsd,!openbsd

package clipboard

func Start() {
}

func Get() string {
	return ""
}

func GetPrimary() string {
	return ""
}

func Set(text string) (e error) {
	return nil
}

func SetPrimary(text string) (e error) {
	return nil
}
