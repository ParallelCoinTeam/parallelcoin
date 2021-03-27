module github.com/p9c/pod

go 1.16

require (
	gioui.org v0.0.0-20201229000053-33103593a1b4
	github.com/BurntSushi/xgb v0.0.0-20200324125942-20f126ea2843
	github.com/VividCortex/ewma v1.1.1
	github.com/aead/siphash v1.0.1
	github.com/atotto/clipboard v0.1.2
	github.com/bitbandi/go-x11 v0.0.0-20171024232457-5fddbc9b2b09
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd
	github.com/btcsuite/golangcrypto v0.0.0-20150304025918-53f62d9b43e8
	github.com/btcsuite/goleveldb v1.0.0
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/conformal/fastsha256 v0.0.0-20160815193821-637e65642941
	github.com/coreos/bbolt v1.3.3
	github.com/davecgh/go-spew v1.1.1
	github.com/enceve/crypto v0.0.0-20160707101852-34d48bb93815
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/gookit/color v1.3.8
	github.com/jackpal/gateway v1.0.7
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kkdai/bstream v1.0.0
	github.com/kr/text v0.2.0 // indirect
	github.com/marusama/semaphore v0.0.0-20190110074507-6952cef993b2
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/niubaoshu/gotiny v0.0.3
	github.com/onsi/ginkgo v1.12.1 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/programmer10110/gostreebog v0.0.0-20170704145444-a3e1d28291b2
	github.com/spacemonkeygo/errors v0.0.0-20201030155909-2f5f890dbc62 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/templexxx/reedsolomon v1.1.3
	github.com/tstranex/gozmq v0.0.0-20160831212417-0daa84a596ba
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/urfave/cli v1.22.4
	github.com/vivint/infectious v0.0.0-20190108171102-2455b059135b
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/atomic v1.6.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/exp v0.0.0-20200924195034-c827fd4f18b9
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/sys v0.0.0-20210326220804-49726bf1d181 // indirect
	golang.org/x/tools v0.1.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	lukechampine.com/blake3 v1.0.0

)

//replace gioui.org => ./gio
