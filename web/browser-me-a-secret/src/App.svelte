<script lang="ts">
	let encryptor: ((message: string, publicKey: string) => string) | undefined;
	let publicKeyer: ((x: string) => string) | undefined;
	const go = new Go();
	WebAssembly.instantiateStreaming(
		fetch("./registerEncryptor.wasm"),
		go.importObject
	).then((result) => {
		go.run(result.instance);
		encryptor = encrypt
		publicKeyer = getPublicKey
	});
	let message: string | undefined;
	let ciphertextElement: HTMLTextAreaElement;

	const copyCiphertext = () => {
			ciphertextElement.select();
			ciphertextElement.setSelectionRange(0, 99999);
			document.execCommand("copy");
	}

	const queryString = window.location.search;
	console.log(queryString)
	console.log(publicKeyer?.("sdf"))
</script>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}
	textarea { width: 100%; height: 200px; }

	textarea:read-only {
  	background-color: #ccc;
	}
</style>


<main>
	<h5>Send Me a Secret (<a href="https://github.com/ostrowr/send-me-a-secret">GitHub</a>)</h5>
	<h4>Encrypt a short message using my public key</h4>
	<textarea bind:value={message} placeholder="Type message here"></textarea>
	<h3>Encrypted</h3>
	{#if !encryptor || !publicKeyer}
		<p>Loading encryptor...</p>
	{:else}
		<textarea bind:this={ciphertextElement} readonly>{encryptor(message ?? "", publicKeyer("d"))}</textarea>
	{/if}
	<button on:click={copyCiphertext}>Copy</button>
</main>


