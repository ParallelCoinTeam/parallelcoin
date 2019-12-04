Vue.config.devtools = true;
Vue.prototype.$eventHub = new Vue(); 
 


var rcvar = {
	osStatus:{},
	osBalance:0,
	osUnconfirmed:0,
	osTxsNumber:0,
	osHashes:0,
	osNetHash:0,
	osHeight:0,
	osBestBlock:0,
	osBlockCount:0,
	osNetlastBlock:0,
	osConnections:0,
	osTransactions:[],
	osTxs:[],
	osLastTxs:{
		"txs":[
			{
				"abandoned": false,
				"account": "default",
				"address": "8dKZQqzdcgxcAxpB86vUHN2oRtoL4cQaQr",
				"amount": 0.03444004,
				"blockhash": "9aed6d7ff84a362bdacabe0fd49b84c368354404205ce5de2e2d62c1abb70381",
				"blocktime": 1574262671,
				"category": "generate",
				"confirmations": 1,
				"generated": true,
				"time": 1574262672,
				"timereceived": 1574262672,
				"trusted": false,
				"txid": "33e3da64ef6d2ef7315683b5a0fa5fe7d4f4d77d564f6cad3214508f70d32d23",
				"vout": 0,
				"walletconflicts": []
			  },
			  {
				"abandoned": false,
				"account": "default",
				"address": "8czZuVuBS13zbZsPKXvGGHx34VnauxHREs",
				"amount": 0.03445372,
				"blockhash": "2f1894d67c5010deefb4cd423c4faf4d37e5b70eceac09f49450df3b5ee32a44",
				"blocktime": 1574262662,
				"category": "generate",
				"confirmations": 3,
				"generated": true,
				"time": 1574262662,
				"timereceived": 1574262662,
				"trusted": false,
				"txid": "85764f758baabbde69a267154b2663d39a311eb5a4205fd59797fc5cc010a314",
				"vout": 0,
				"walletconflicts": []
			  }			
		],
		"txsnumber":0
	},
	osLocalHost:[],
	isScreen:'PageOverview',
};

var duOStransactions = {
	"txs":[
		{
			"abandoned": false,
			"account": "default",
			"address": "8dKZQqzdcgxcAxpB86vUHN2oRtoL4cQaQr",
			"amount": 0.03444004,
			"blockhash": "9aed6d7ff84a362bdacabe0fd49b84c368354404205ce5de2e2d62c1abb70381",
			"blocktime": 1574262671,
			"category": "generate",
			"confirmations": 1,
			"generated": true,
			"time": 1574262672,
			"timereceived": 1574262672,
			"trusted": false,
			"txid": "33e3da64ef6d2ef7315683b5a0fa5fe7d4f4d77d564f6cad3214508f70d32d23",
			"vout": 0,
			"walletconflicts": []
		  },
		  {
			"abandoned": false,
			"account": "default",
			"address": "8czZuVuBS13zbZsPKXvGGHx34VnauxHREs",
			"amount": 0.03445372,
			"blockhash": "2f1894d67c5010deefb4cd423c4faf4d37e5b70eceac09f49450df3b5ee32a44",
			"blocktime": 1574262662,
			"category": "generate",
			"confirmations": 3,
			"generated": true,
			"time": 1574262662,
			"timereceived": 1574262662,
			"trusted": false,
			"txid": "85764f758baabbde69a267154b2663d39a311eb5a4205fd59797fc5cc010a314",
			"vout": 0,
			"walletconflicts": []
		  }			
	],
	"txsnumber":0
};


var duOSaddressbook = {
	theme:false,
};

var duOSblocks = {
	theme:false,
};



var duoSystem = {
	theme:false,
	isBoot:false,
	isLoading:false,
	isDev:true,
	isScreen:'PageOverview',
	timer: '',
	status: {
	"ver":"0.0.1",
	"walletver":{
		"podjsonrpcapi":{
			"versionstring":"1.3.0",
			"major":1,
			"minor":3,
			"patch":0,
			"prerelease":"",
			"buildmetadata":""
		}
	},
	"uptime":1573942261,
	"cpu":[{
		"cpu":0,
		"vendorId":"GenuineIntel",
		"family":"0x6 ",
		"model":"0x55 ",
		"stepping":4,
		"physicalId":"",
		"coreId":"",
		"cores":20,
		"modelName":"Intel(R) Core(TM) i9-7900X CPU @ 3.30GHz",
		"mhz":3312,
		"cacheSize":0,
		"flags":["fpu","vme","de","pse","tsc","msr","pae","mce","cx8","apic","sep","mtrr","pge","mca","cmov","pat","pse36","clflush","dts","acpi","mmx","fxsr","sse","sse2","ss","htt","tm","pbe","sse3","pclmulqdq","dtes64","mon","ds_cpl","vmx","est","tm2","ssse3","sdbg","fma","cx16","xtpr","pdcm","pcid","dca","sse4.1","sse4.2","x2apic","movbe","popcnt","tscdlt","aesni","xsave","osxsave","avx","f16c","rdrand","syscall","nx","page1gb","rdtscp","lm","lahf","abm","prefetch","fsgsbase","tscadj","bmi1","hle","avx2","fdpexc","smep","bmi2","erms","invpcid","rtm","pqm","nfpusg","mpx","pqe","avx512f","avx512dq","rdseed","adx","smap","clflushopt","clwb","proctrace","avx512cd","avx512bw","avx512vl","xsaveopt","xsavec","xinuse","xsaves"],
		"microcode":""}
	],
	"cpupercent":[0,0,0,0,100,0,100,0,0,0,100,0,100,0,0,0,0,0,0,0],
	"mem":{
		"total":33953574912,"available":16904683520,"used":17048891392,
		"usedPercent":50.21236036613781,"free":744972288,"active":12869365760,
		"inactive":8563982336,"wired":3648053248,
		"laundry":7595728896,"buffers":1481977856,"cached":0,
		"writeback":0,"dirty":0,"writebacktmp":0,
		"shared":0,"slab":0,"sreclaimable":0,"sunreclaim":0,"pagetables":0,"swapcached":0,"commitlimit":0,"committedas":0,"hightotal":0,"highfree":0,
		"lowtotal":0,"lowfree":0,"swaptotal":0,"swapfree":0,"mapped":0,"vmalloctotal":0,"vmallocused":0,"vmallocchunk":0,"hugepagestotal":0,
		"hugepagesfree":0,
		"hugepagesize":0
	},
	"disk":{
		"path":"/","fstype":"ufs","total":460514775040,"free":23478919168,"used":400194674688,"usedPercent":94.45825288418138,
		"inodesTotal":58185598,"inodesUsed":2644615,"inodesFree":55540983,"inodesUsedPercent":4.545136753600092
	},
	"net":"TestNet3",
	"chain":"testnet",
	"hashrate":0,
	"height":0,
	"bestblockhash":"2c63958400f9e1b5cf00a19fe8e75bec459ff2c38636c69b79c06d7989d3c723",
	"networkhashrate":0,
	"diff":1e-7,
	"blockcount":0,
	"connectioncount":0,
	"networklastblock":0,
	balance: {
		"balance":"0",
		"unconfirmed":"0"
		},
	},
	transactions:{
		"txs":[
			{
				"abandoned": false,
				"account": "default",
				"address": "8dKZQqzdcgxcAxpB86vUHN2oRtoL4cQaQr",
				"amount": 0.03444004,
				"blockhash": "9aed6d7ff84a362bdacabe0fd49b84c368354404205ce5de2e2d62c1abb70381",
				"blocktime": 1574262671,
				"category": "generate",
				"confirmations": 1,
				"generated": true,
				"time": 1574262672,
				"timereceived": 1574262672,
				"trusted": false,
				"txid": "33e3da64ef6d2ef7315683b5a0fa5fe7d4f4d77d564f6cad3214508f70d32d23",
				"vout": 0,
				"walletconflicts": []
			  },
			  {
				"abandoned": false,
				"account": "default",			
				"address": "8dw3dFzvd",		
				"amount": 0.03444688,
				"blockhash": "77f710efcb0969c8109cc08dc2531ddc73c2ab1adf26e650c9afc4c4daf88823",
				"blocktime": 1574262662,
				"category": "generate",
				"confirmations": 2,
				"generated": true,
				"time": 1574262662,
				"timereceived": 1574262662,
				"trusted": false,
				"txid": "7eb39549bc2acc23e571746a9bf253f653f5b3bb9c80846773fefc42b6989361",
				"vout": 0,
				"walletconflicts": []
			  },
			  {
				"abandoned": false,
				"account": "default",
				"address": "8czZuVuBS13zbZsPKXvGGHx34VnauxHREs",
				"amount": 0.03445372,
				"blockhash": "2f1894d67c5010deefb4cd423c4faf4d37e5b70eceac09f49450df3b5ee32a44",
				"blocktime": 1574262662,
				"category": "generate",
				"confirmations": 3,
				"generated": true,
				"time": 1574262662,
				"timereceived": 1574262662,
				"trusted": false,
				"txid": "85764f758baabbde69a267154b2663d39a311eb5a4205fd59797fc5cc010a314",
				"vout": 0,
				"walletconflicts": []
			  }			
		],
		"txsnumber":0
	},
};

var	duoBlocks = {};
var	duoAaddressBook = {};