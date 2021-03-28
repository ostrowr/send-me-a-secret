<script lang="ts">
import Encryptor from "./components/Encryptor.svelte";
import GithubSelector from "./components/GithubSelector.svelte";

	let encryptFn: ((message: string, publicKey: string) => string) | undefined;
	let getValidPublicKeyFn: ((possibleKeys: string[]) => string) | undefined;
	const go = new Go();
	WebAssembly.instantiateStreaming(
		fetch("./registerEncryptor.wasm"),
		go.importObject
	).then((result) => {
		go.run(result.instance);
		encryptFn = encrypt
		getValidPublicKeyFn = getValidPublicKey
	});
	// const queryString = window.location.search;
	let publicKey: Promise<string>;
	let githubUsername: string = "ostrowr";
</script>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}
</style>


<main>
	<h5>Send Me a Secret (<a href="https://github.com/ostrowr/send-me-a-secret">GitHub</a>)</h5>
	{#if !encryptFn || !getValidPublicKeyFn}
		Loading...
	{:else}
		<GithubSelector bind:username={githubUsername} getValidPublicKeyFn={getValidPublicKeyFn} bind:publicKeyPromise={publicKey}/>
		{#await publicKey}
			<p>...waiting</p>
		{:then key}
			<Encryptor publicKey={key} encryptFn={encryptFn} username={githubUsername}/>
		{/await}
	{/if}

</main>


