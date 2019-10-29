<script>
	import IcoOverview from "./ico/IcoOverview.svelte";
	import IcoHistory from "./ico/IcoHistory.svelte";
	import IcoAddressBook from "./ico/IcoAddressBook.svelte";
	import IcoSettings from "./ico/IcoSettings.svelte";

</script>

<style>
	#menu{

    }


    #menu svg{
    	width:48px;
    	height:48px;
    }

</style>

<nav class="Nav">
<ul id="menu" class="lsn noPadding">
	<li id='menuoverview' class='sidebar-item current'>
		<button onclick="external.invoke('overview')" class="noMargin noPadding noBorder bgTrans sXs"><IcoOverview /></button>
	</li>
	<li id='menutransactions' class='sidebar-item'>
		<button onclick="external.invoke('hisotry')" class="noMargin noPadding noBorder bgTrans sXs"><IcoHistory /></button>
	</li>
	<li id='menuaddressbook' class='sidebar-item'>
		<button onclick="external.invoke('addressbook')" class="noMargin noPadding noBorder bgTrans sXs"><IcoAddressBook /></button>
	</li>
	<li id='menublockexplorer' class='sidebar-item'>
		<button onclick="external.invoke('overview')" class="noMargin noPadding noBorder bgTrans sXs"><IcoOverview /></button>
	</li>
	<li id='menusettings' class='sidebar-item'>
		<button onclick="external.invoke('settings')" class="noMargin noPadding noBorder bgTrans sXs"><IcoSettings /></button>
	</li>
</ul>
</nav>