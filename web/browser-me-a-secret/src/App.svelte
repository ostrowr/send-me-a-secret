<script lang="ts">
	let encryptor: ((message: string) => string) | undefined;
	const go = new Go();
	WebAssembly.instantiateStreaming(
		fetch("/registerEncryptor.wasm"),
		go.importObject
	).then((result) => {
		go.run(result.instance);
		encryptor = encrypt
	});
	let message: string | undefined;
	let ciphertextElement: HTMLTextAreaElement;

	const copyCiphertext = () => {
			ciphertextElement.select();
			ciphertextElement.setSelectionRange(0, 99999);
			document.execCommand("copy");
	}
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

	p {
		overflow: auto;
		/* word-wrap: "break-word"; */
		max-width: 240px;
		/* white-space: nowrap; */
  	/* overflow: hidden; */
  	text-overflow: clip;
	}
</style>


<main>
	<h5>Send Me a Secret (<a href="https://github.com/ostrowr/send-me-a-secret">GitHub</a>)</h5>
	<h4>Encrypt a short message using my public key</h4>
	<textarea bind:value={message} placeholder="Type message here"></textarea>
	<h3>Encrypted</h3>
	{#if !encryptor}
		<p>Loading encryptor...</p>
	{:else}
		<textarea bind:this={ciphertextElement} readonly>{encryptor(message ?? "")}</textarea>
	{/if}
	<button on:click={copyCiphertext}>Copy</button>
</main>


