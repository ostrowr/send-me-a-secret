<script lang="ts">
  export let publicKey: string;
  export let encryptFn: (message: string, publicKey: string) => string;
  export let username: string;

  let message: string | undefined;

  let ciphertextElement: HTMLTextAreaElement;

  const copyCiphertext = () => {
			ciphertextElement.select();
			ciphertextElement.setSelectionRange(0, 99999);
			document.execCommand("copy");
	}
</script>

<style>
	textarea { width: 100%; height: 200px; }

	textarea:read-only {
  	background-color: #ccc;
	}
</style>

<h4>Encrypt a short message using {username}'s public key</h4>
<textarea bind:value={message} placeholder="Type message here"></textarea>
<h3>Encrypted</h3>
<textarea bind:this={ciphertextElement} readonly>{encryptFn(message ?? "", publicKey ?? "")}</textarea>
<button on:click={copyCiphertext}>Copy</button>
