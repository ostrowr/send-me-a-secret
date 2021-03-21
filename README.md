# Send Me a Secret

Easily encrypt a short message such that only I can read it.

Just want to encrypt a message? Go [here](https://ostro.ws/send-me-a-secret)

## Why?

I often find myself needing a key sent to me, and not particularly trusting normal mechanisms of communication (Slack, email, etc) to prevent anyone else from reading it.

Existing encryption tools are generally super complicated and unapproachable, especially for a non-technical user.

Other alternatives like Signal are probably good – but they require you to make an account! Plus, I've never actually read the code.

This library allows you to easily do 2 things:

1. Generate matching `encrypt` and `decrypt` binaries. The `decrypt` binary contains your private key, so don't share it! It's much easier than remembering opaque gpg commands.
2. Run a tiny web app that uses `encrypt` compiled into WebAssembly. This way, you can just send somebody a link (e.g. [ostro.ws/send-me-a-secret](https://ostro.ws/send-me-a-secret)) to easily encrypt a message using my public key. If they don't trust their browser (or their browser extensions) they can always download and run one of the `encrypt` binaries directly.

## Limitations

- Messages are capped at about 500 bytes, since they're encrypted using RSA with a keysize of 4096 bits. If longer messages are necessary, this could be modified to just use RSA for key exchange and then continue with normal AES encryption, but at that point you should probably use an encrypted channel.
- There is no attempt at authenticated encryption. These secrets are for communicating things like keys once you already have a trusted channel like Slack.
- Definitely no forward secrecy; we're just encrypting using these RSA keys.
- Definitely not audited. Probably lots of security issues. Dangerously close to rolling my own encryption.

## Setting up

1. Run `go run generate.go`. This will generate a private/public key pair and save them next to their respective templates. Don't commit these files!
2. Run ./build.sh
3. Save build/private/decrypt somewhere secure. This is your private key, and allows you to decrypt these secrets.

## Running in the browser

See [web/browser-me-a-secret](./web/browser-me-a-secret). This is a tiny web-app built in Svelte.

Run `npm run dev` from within `browser-me-a-secret` to run a development server. You'll also have to run `./build.sh` to generate `registerEncryptor.wasm`, which the web app needs.

To deploy to github pages, run `npm run deploy`.
