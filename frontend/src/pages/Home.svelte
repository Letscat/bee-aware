<script lang="ts">
 import { onMount } from "svelte";
 import {getCookie} from "../utils/cookie";
    import Camera from "../components/Camera.svelte";

 onMount(() => {
    const session=getCookie("session");
    console.log(session);
     if(!session){
        window.location.href="#/login";
     }
    })

</script>
<Camera />


