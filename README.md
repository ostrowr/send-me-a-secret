# Send Me a Secret

Messages are capped at about 500 bytes, since they're encrypted using RSA with a keysize of 4096 bits.

If longer messages are necessary, this could be modified to just use RSA for key exchange and then continue with normal AES encryption, but at that point you should probably use an encrypted channel.

- No authentication
- For when you already have a trusted channel, but maybe don't trust that channel not to be reading your messages
