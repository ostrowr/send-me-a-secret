<script lang="ts">
  export let username: string;
  export let getValidPublicKeyFn: (possibleKeys: string[]) => string;
  export let publicKeyPromise: Promise<string>;
  type Key = {id: number, key: string}

  async function getKeyForUser(username: string): Promise<string> {
    const response = await fetch(`https://api.github.com/users/${username}/keys`, {
      headers: {
        Accept: "application/vnd.github.v3+json"
      }
    })
    if (!response.ok) {
      throw new Error(response.statusText)
    }
    const keys: Key[] = await response.json()
    const validKey = getValidPublicKeyFn(keys.map(k => k.key))
    return validKey;
  }

  // todo: don't make a request every time this changes. Rate limited
  $: publicKeyPromise = getKeyForUser(username)

</script>

<style>

</style>

<input bind:value={username}>

{#await publicKeyPromise}
	<p>...waiting</p>
{:then key}
  <p>{username}'s key is</p>
  <pre>{key}</pre>
{:catch error}
	<p style="color: red">{error.message}</p>
{/await}
