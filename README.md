# Send Me a Secret

Easily encrypt a short message such that only I can read it.

Just want to encrypt a message? Go [here](https://ostro.ws/send-me-a-secret)

## Why?

I often find myself needing a key sent to me, and don't particularly trust normal mechanisms of communication (Slack, email, etc) to prevent anyone else from reading it.

Existing encryption tools are generally super complicated and unapproachable, especially for a non-technical user.

Other alternatives like Signal are probably good – but they require you to make an account! Plus, I've never actually read the code.

This library allows you to easily do 2 things:

1. Generate a private-public key pair and (optionally) upload the public key to GitHub for anyone to use! After running `initialize`, you shouldn't ever have to remember opaque gpg commands or the path to your private key.
2. Run a tiny web app that uses `encrypt` compiled into WebAssembly. This way, you can just send somebody a link (e.g. [ostro.ws/send-me-a-secret](https://ostro.ws/send-me-a-secret)) to easily encrypt a message using my public key. If they don't trust their browser (or their browser extensions) they can always download and run the `send-me-a-secret` binary directly.

## Usage (sending secrets – encryption)

Anybody can send a secret without installing any software.

Just go to [ostro.ws/send-me-a-secret](https://ostro.ws/send-me-a-secret) and type in your message. Copy the encrypted message and send that to the recipient!

There is no tracking code on that website, but if you don't trust your browser (or browser extensions), you can install the `send-me-a-secret` binary and run

`./send-me-a-secret encrypt --user ostrowr "this message is going to get so encrypted"`

## Usage (receiving secrets – decryption)

First, run `./send-me-a-secret initialize`. This generates new keys and uploads your public key to your GitHub account.

```
$ ./send-me-a-secret initialize
Initializing send-me-a-secret
Enter passphrase (empty for no passphrase):
Generating new private key (this may take a moment)...
Private key generated. It's saved in /Users/robbie/.send-me-a-secret

send-me-a-secret adds a new public key to your GitHub account. If you'd prefer to manually add the key to your account, re-run `initialize` with the --skip-github flag.
Enter a GitHub token with write:public_key access. (You can generate one at https://github.com/settings/tokens/new):
Public key uploaded

Getting authenticated GitHub username...
Got username: ostrowr

Validating send-me-a-secret
Fetching public key from GitHub...
Public key successfully fetched

Encrypting a test message using that public key
Message successfully encrypted

Decrypting the test message using your private key
Message successfully decrypted

Validation succeeded! You're good to go.
```

Once you've initialized, anyone (without any authentication) can encrypt messages for you!

```
$ ./send-me-a-secret encrypt --user ostrowr "this message is going to get so encrypted" # you can also pass a message in via stdin
Getting public key from GitHub...
Encryption succeeded! Ciphertext below:

iq13TafMuPDE+0SE+gdGT024Nq63twOoxXpJrh4mm+Q/RzCp3500JhqqfTkoxBmxuq0oLL7UJ/DDvprvk7dlyPjFwHfnrzEAI3AFOGRYXIGcbxOrwFNefOTF8OEDX5cN1q+PPUlsklVXle9uuvb+iw1LB4wZjJviz56kEX0qC1vTFnT/IfdD/nnXACkaMWvZKe8rC4tCRBaZISgkQeYzUJNY6qqwvKv5PB95nrSU0Lh55FdvWcKvpIdCRFbt7FrxJcal14u59w7qfBQHe4NE9rk9AssDSJWh4Qakl2B2hDHUNR/LWHTOtKui31lHV1RfCuGhEHym88d19dlEe3niSzGfcwSovmvCUvfjxhPMsWGK8qvYkV8ALV9wrsLpUwP7kel3GSAVUtGgmJX7FoUOqgD+s41oFU12ul75T1gNI+z38sJCpyjqYiW6RjC/LRDDaVSAEJBrK652JxuLkbpu97xLuUYufPLUojWVVwHnBakoMXiZm14c2yqzjhHy4CQT8Rl++VZkL64Of2f1+MDbiagnynmysIO5qeSVVKD4l8gtMNz9uD2Kk0fgn7itt6NWdDtvOqXzc20ysYetlJoBj+fENQjhvgsi8aYOTL1vywaZOQRLSXZgwXuGv1JGHLSNkZdJuhvWRBYS+jDVU9D++Q9xAXvHrFr/xOR8j2ZL1DHl9DZEbWvrwkV97W3LNZJLa15jWWpoZ+5uqo7Km2F4WEo60YN9kfmrj9u62FkHk9l43z+dfuOBr7Jfpg==
```

Only you can decrypt this message, since your private key is saved locally:

```
$ pbpaste | ./send-me-a-secret decrypt # you can also pass in ciphertext/password/etc via an argument
this message is going to get so encrypted
```

## Running in the browser

Messages can be encrypted in the browser. Since public keys are stored in GitHub, you can encrypt a message for anyone who has initialized send-me-a-secret by knowing their GitHub username. If you don't have a GitHub account, you can (soon) also distribute a URL with your key baked in.

See [web/browser-me-a-secret](./web/browser-me-a-secret). This is a tiny web-app built in Svelte.

Run `npm run dev` from within `browser-me-a-secret` to run a development server. You'll also have to run `./build.sh` to generate `registerEncryptor.wasm`, which the web app needs.

To deploy to github pages, run `npm run deploy`.

## Limitations

- Messages are capped at about 500 bytes, since they're encrypted using RSA with a keysize of <5000 bits. If longer messages are necessary, this could be modified to just use RSA for key exchange and then continue with normal AES encryption, but at that point you should probably use an encrypted channel.
- There is no attempt at authenticated encryption. These secrets are for communicating things like keys once you already have a trusted channel like Slack.
- Definitely no forward secrecy; we're just encrypting using these RSA keys.
- Definitely not audited. Probably lots of security issues. Dangerously close to rolling my own encryption.

## Trivia

The unauthenticated API for GitHub public keys doesn't return any metadata about the keys. In order to ensure that the key is the right one, send-me-a-secret uses a nontraditional key length and just checks the GitHub account to see if there are any keys of that length. This is very brittle; I'm hoping for ideas here!
