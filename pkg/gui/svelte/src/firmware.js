import { writable } from "svelte/store";

export const right = writable(false);
export const persistent = writable(true);
export const elevation = writable(false);
export const showNav = writable(true);
export const showNavMobile = writable(false);
export const breakpoint = writable("");

export const duoSystem = {
	theme:false,
	isBoot:true,
	isBootMenu:true,
	isBootLogo:true,
	isLoading:false,
	isDev:true,
	isScreen:'overview',
	timer: '',
};

export const duoConfig = writable(true);
export const duoNode = writable(true);
export const dUoWallet = writable(true);
export const duoStatus = writable(true);
export const duoBalance = writable(true);
export const duoConnections = 0;
export const duoAddressBook = writable(true);
export const duoTransactions = writable(true);
export const duoPeerInfo = writable(true);
export const duoBlocks = [];