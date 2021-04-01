module github.com/p9c/pod

go 1.16

require (
	gioui.org v0.0.0-20201229000053-33103593a1b4
	github.com/VividCortex/ewma v1.1.1
	github.com/aead/siphash v1.0.1
	github.com/atotto/clipboard v0.1.4
	github.com/bitbandi/go-x11 v0.0.0-20171024232457-5fddbc9b2b09
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd
	github.com/btcsuite/golangcrypto v0.0.0-20150304025918-53f62d9b43e8
	github.com/btcsuite/goleveldb v1.0.0
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/conformal/fastsha256 v0.0.0-20160815193821-637e65642941
	github.com/coreos/bbolt v1.3.3
	github.com/davecgh/go-spew v1.1.1
	github.com/enceve/crypto v0.0.0-20160707101852-34d48bb93815
	github.com/gookit/color v1.3.8
	github.com/jackpal/gateway v1.0.7
	github.com/jessevdk/go-flags v1.4.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kkdai/bstream v1.0.0
	github.com/marusama/semaphore v0.0.0-20190110074507-6952cef993b2
	github.com/niubaoshu/gotiny v0.0.3
	github.com/p9c/gel v0.1.5
	github.com/p9c/log v0.0.6
	github.com/p9c/opts v0.0.5
	github.com/p9c/qu v0.0.3
	github.com/programmer10110/gostreebog v0.0.0-20170704145444-a3e1d28291b2
	github.com/templexxx/reedsolomon v1.1.3
	github.com/tstranex/gozmq v0.0.0-20160831212417-0daa84a596ba
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/urfave/cli v1.22.5
	github.com/vivint/infectious v0.0.0-20190108171102-2455b059135b
	go.uber.org/atomic v1.7.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/exp v0.0.0-20200924195034-c827fd4f18b9
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/sys v0.0.0-20210326220804-49726bf1d181 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.2.7 // indirect
	lukechampine.com/blake3 v1.0.0

)

//replace gioui.org => ./gio
