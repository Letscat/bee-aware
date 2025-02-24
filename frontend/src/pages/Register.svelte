<script lang="ts">
    let email:string=$state("");
    let userName:string=$state("");
    let password:string=$state("");
    let error:string=$state("");

    async function handleSubmit():Promise<void>{
        if(!email || !userName || !password) return;
        let result=await fetch("/register", {
            method: "POST",
            body: JSON.stringify({
                email: email,
                userName: userName,
                password: password
            })
        })
        if(!result.ok){
            error=result.statusText || "Something went wrong";
        }
        window.location.href="#/login";
    }

</script>

<main>
    <h1>Register</h1>
    <form onsubmit={handleSubmit}>
        <label for="mail">Mail:</label>
        <input type="text" id="mail" bind:value={email} />
        <br/>
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

</main>

<style>
</style>
