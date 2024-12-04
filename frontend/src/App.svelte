<script lang="ts">
	import Router, {
		location
	} from "svelte-spa-router";
	import Home from "./pages/Home.svelte";

</script>

<nav>
	<a href="#/">Home</a>

</nav>

<Router
	routes={{
		"/": Home,
	}}
/>

<hr />
/#{$location}
