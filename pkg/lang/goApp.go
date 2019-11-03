package lang

func goAppDict() Com {
	return Com{
		Component: "goApp",
		Languages: []Language{
			{
				Code: "en",
				Definitions: []Text{
					{
						ID: "DESCRIPTION",
						Definition: "Parallelcoin Pod Suite -- All-in-one everything" +
							" for Parallelcoin!",
					},
					{
						ID: "COPYRIGHT",
						Definition: "Legacy portions derived from btcsuite/btcd under" +
							" ISC licence. The remainder is already in your" +
							" possession. Use it wisely.",
					},
					{
						ID:         "NOSUBCMDREQ",
						Definition: "no subcommand requested",
					},
					{
						ID:         "STARTINGNODE",
						Definition: "starting node",
					},
				},
			},
			{
				Code: "rs",
				Definitions: []Text{
					{
						ID: "DESCRIPTION",
						Definition: "ParalelniNovcic Pod Programski Paket -- Sve-u-jednom garnitura" +
							" za ParalelniNovcic!",
					},
					{
						ID: "COPYRIGHT",
						Definition: "Izvedeni deovi iz btcsuite/btcd su pod" +
							" ISC licencom. Sve ostalo je u vasem" +
							" vlasnistvu. Koristite se DUO-m mudro!",
					},
					{
						ID:         "NOSUBCMDREQ",
						Definition: "nije pozvana pod komanda",
					},
					{
						ID:         "STARTINGNODE",
						Definition: "pokrece se ",
					},
				},
			},
		},
	}
}
