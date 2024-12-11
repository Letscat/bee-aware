<script lang="ts">
    let userName:string=$state("");
    let password:string=$state("");
    let error:string=$state("");
    async function handleSubmit():Promise<void>{
        if(!userName || !password) return;

       /*  setCookie("userSession", userName, 1); */
       const result=await fetch("/login", {
           method: "POST",
           headers: {
               "Content-Type": "application/json"
           },
           body: JSON.stringify({
               userName,
               password
           })
       })
       if(!result.ok){
           error= result.statusText;
       }
       window.location.href="#/";
    }

</script>

<main>
    <form onsubmit={handleSubmit}>
        <label for="username">Username:</label>
        <input type="text" id="username" bind:value={userName} />
        <br />
        <label for="password">Password:</label>
        <input type="password" id="password" bind:value={password} />
        <br />
        <button type="submit">Login</button>
        {#if error}
          <p style="color: red;">{error}</p>
        {/if}
      </form>
      <a href="#/register">Register</a>
</main>

<style>
</style>
