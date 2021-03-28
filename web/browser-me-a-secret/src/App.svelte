<script lang="ts">
import Encryptor from "./components/Encryptor.svelte";

	let encryptFn: ((message: string, publicKey: string) => string) | undefined;
	let getPublicKeyFn: ((x: string) => string) | undefined;
	const go = new Go();
	WebAssembly.instantiateStreaming(
		fetch("./registerEncryptor.wasm"),
		go.importObject
	).then((result) => {
		go.run(result.instance);
		encryptFn = encrypt
		getPublicKeyFn = getPublicKey
	});
	const queryString = window.location.search;
	console.log(queryString)
	let publicKey: string | undefined;
	$: publicKey = getPublicKeyFn?.("TODO");
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
	{#if !encryptFn || !publicKey}
		Loading...
	{:else}
		<Encryptor publicKey={publicKey} encryptFn={encryptFn}/>
	{/if}

</main>


