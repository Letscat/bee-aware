<script lang="ts">
	import Router, {
		location
	} from "svelte-spa-router";
	import Home from "./pages/Home.svelte";
	import Login from "./pages/Login.svelte";
	import Register from "./pages/Register.svelte";

</script>

<nav>
	<a href="#/">Home</a>
	<a href="#/login">Login</a>
	<a href="#/register">Register</a>

</nav>

<Router
	routes={{
		"/": Home,
		"/login": Login,
		"/register": Register
	}}
/>

<hr />
/#{$location}
