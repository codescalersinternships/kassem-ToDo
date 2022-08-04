<script>
import Header from "./components/header.svelte"
import List from "./components/list.svelte"
</script>
 
<style>
	.container {
		width: 500px;
		background: #1594a5 ;
		margin: auto
	}
	:global(.container *){
		font-family: sans-serif;
	}
</style>
<div class="container">
	<Header />
	<List />

</div>