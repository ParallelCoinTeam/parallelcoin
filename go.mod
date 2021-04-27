module github.com/p9c/matrjoska

go 1.12

require (
	github.com/VividCortex/ewma v1.1.1
	github.com/aead/siphash v1.0.1
	github.com/atotto/clipboard v0.1.4
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd
	github.com/btcsuite/golangcrypto v0.0.0-20150304025918-53f62d9b43e8
	github.com/btcsuite/goleveldb v1.0.0
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/conformal/fastsha256 v0.0.0-20160815193821-637e65642941
	github.com/davecgh/go-spew v1.1.1
	github.com/enceve/crypto v0.0.0-20160707101852-34d48bb93815
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gookit/color v1.3.8
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/jackpal/gateway v1.0.7
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kkdai/bstream v1.0.0
	github.com/marusama/semaphore v0.0.0-20190110074507-6952cef993b2
	github.com/niubaoshu/gotiny v0.0.3
	github.com/p9c/gel v0.1.16
	github.com/p9c/gio v0.0.5
	github.com/p9c/interrupt v0.0.2
	github.com/p9c/log v0.0.9
	github.com/p9c/opts v0.0.9
	github.com/p9c/qu v0.0.3
	github.com/programmer10110/gostreebog v0.0.0-20170704145444-a3e1d28291b2
	github.com/tstranex/gozmq v0.0.0-20160831212417-0daa84a596ba
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/vivint/infectious v0.0.0-20200605153912-25a574ae18a3
	go.etcd.io/bbolt v1.3.5
	go.uber.org/atomic v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/exp v0.0.0-20210417010653-0739314eea07
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
	lukechampine.com/blake3 v1.1.5
)

//replace gioui.org => github.com/p9c/gio v0.0.3
replace (
	gioui.org => ./pkg/gel/gio
	github.com/p9c/gel => ./pkg/gel
	github.com/p9c/log => ./pkg/log
	github.com/p9c/opts => ./pkg/opts
)
