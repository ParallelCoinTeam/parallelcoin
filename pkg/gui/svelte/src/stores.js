import { writable } from "svelte/store";

export const right = writable(false);
export const persistent = writable(true);
export const elevation = writable(false);
export const showNav = writable(true);
export const showNavMobile = writable(false);
export const breakpoint = writable("");

export const duoSystem = {
	config:writable(true),
	node:writable(true),
	wallet:writable(true),
	status:writable(true),
	balance:writable(true),
	connectionCount:0,
	addressBook:writable(true),
	transactions:writable(true),
	peerInfo:writable(true),
	blocks:[],
	theme:false,
	isBoot:false,
	isLoading:false,
	isDev:true,
	isScreen:'overview',
	timer: '',
};
