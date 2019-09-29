package css

func AMPgrid() string {
	return `
.grid-container { height: 100%; margin: 0; }
.grid-container * {
border: 1px solid red;
position: relative;
}
.grid-container *:after {
content:attr(class);
position: absolute;
top: 0;
left: 0;
}
.grid-container {
display: grid;
grid-template-columns: 60px 1fr;
grid-template-rows: 60px 1fr;
grid-template-areas: "Logo Header" "Sidebar Main";
}
.Logo { grid-area: Logo; }
.Header {
display: grid;
grid-template-columns: 1fr 1fr 1fr 60px 180px 120px 120px 60px;
grid-template-rows: 1fr;
grid-template-areas: "h1 h2 h3 h4 h5 h6 h7 h8";
grid-gap:5px;
grid-area: Header;
}
.h1 { grid-area: h1; }
.h2 { grid-area: h2; }
.h3 { grid-area: h3; }
.h4 { grid-area: h4; }
.h5 { grid-area: h5; }
.h6 { grid-area: h6; }
.h7 { grid-area: h7; }
.h8 { grid-area: h8; }
.Sidebar {
display: grid;
grid-template-columns: 1fr;
grid-template-rows: 600px 1fr 60px;
grid-template-areas: "Nav" "Side" "Open";
grid-area: Sidebar;
}
.Open { grid-area: Open; }
.Nav { grid-area: Nav; }
.Side { grid-area: Side; }
.Main {
padding:15px;
display: grid;
grid-template-columns: 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr 1fr;
grid-template-rows: 1fr 1fr 1fr 1fr 1fr 1fr;
grid-template-areas: "Balance Balance Balance Balance Address Address Address Address Address" "Txs Txs Txs Txs Amount Amount Amount Amount Amount" "Txs Txs Txs Txs Mid Mid Mid NetHR NetHR" "Txs Txs Txs Txs Mid Mid Mid LocalHR LocalHR" "BottomLeft BottomLeft Bottom Bottom Mid Mid Mid Log Log" "BottomLeft BottomLeft Bottom Bottom Status Status Status Log Log";
grid-gap:15px;
grid-area: Main;
}
.Balance { grid-area: Balance; }
.Address { grid-area: Address; }
.Amount { grid-area: Amount; }
.NetHR { grid-area: NetHR; }
.LocalHR { grid-area: LocalHR; }
.BottomLeft {
display: grid;
grid-template-columns: 1fr 1fr;
grid-template-rows: 1fr 1fr 1fr;
grid-template-areas: "b1 b2" "b4 b3" "b5 b6";
grid-area: BottomLeft;
}
.b1 { grid-area: b1; }
.b2 { grid-area: b2; }
.b3 { grid-area: b3; }
.b4 { grid-area: b4; }
.b5 { grid-area: b5; }
.b6 { grid-area: b6; }
.Bottom { grid-area: Bottom; }
.Log { grid-area: Log; }
.Status { grid-area: Status; }
.Mid { grid-area: Mid; }
.Txs { grid-area: Txs; }
`
}
