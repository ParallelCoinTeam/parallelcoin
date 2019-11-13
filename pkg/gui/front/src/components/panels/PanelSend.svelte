<script>
  import { Button } from "smelte";
  import { TextField } from "smelte";

</script>

<div class="rwrap flx">
	<div class="flx flc fii justifyBetween">
		<TextField id="sendAddress" label="Test label" placeholder="Enter DUO address" outlined hint="Test hint" class="e-outline flx noMargin fii bgfff" />
 		<div class="flx fii">
			<TextField label="Test label" class="e-outline noMargin fii bgfff" />  
			<Button id="dialogbtn" v-on:click.native="btnClick">Send</Button>
  		</div>
	</div>
</div>