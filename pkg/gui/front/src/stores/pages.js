import {writable} from 'svelte/store';


import PageOverview from '../components/pages/PageOverview.svelte'
import PageTransactions from '../components/pages/PageTransactions.svelte'
import PageAddressBook from '../components/pages/PageAddressBook.svelte'
import PageExplorer from '../components/pages/PageExplorer.svelte'
import PageSettings from '../components/pages/PageSettings.svelte'
import PageNotFound from '../components/pages/PageNotFound.svelte'


function createPages() {
    const {subscribe, set} = writable(PageOverview);

    return {
        subscribe,
        overview: () => set(PageOverview),
        transactions: () => set(PageTransactions),
        addressbook: () => set(PageAddressBook),
        explorer: () => set(PageExplorer),
        settings: () => set(PageSettings),
        notfound: () => set(PageNotFound),
    };
}

export const isPage = createPages();