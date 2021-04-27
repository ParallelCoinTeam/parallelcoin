package btcjson


// UnusableFlags are the command usage flags which this utility are not able to use. In particular it doesn't support
// websockets and consequently notifications.
const UnusableFlags = UFWebsocketOnly | UFNotification
