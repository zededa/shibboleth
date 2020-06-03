# Shibboleth

Shibboleth is a golang library to package messages with body protection: encryption and digital signing.

Why? Don't we have TLS?

* Sometimes TLS just isn't available
* Sometimes TLS just isn't enough

TLS, as [Erik](https://twitter.com/nordmark_erik) described it, is a fighter jet. It is a modern platform that contains
everything and anything, can fly hundreds of miles at Mach 3 and carry enough weapons or observation platforms to win a
war (or at least a decent-sized battle).

Sometimes, you just need a good bicycle: solid, reliable, and has just what you need.

Shibboleth is here to provide two of the most basic TLS guarantees when you either don't have full TLS or cannot
rely on it.

* Protection: encrypting a message so prying eyes cannot read it
* Authentication: validating a message so you can be sure the sender actually sent it and it has not been tampered with

You can use one, the other or both.

When would you use it?

First, if you do not have TLS. For example, you need to distribute a blob of data via non-network channels such as a USB key or
hard drive. None of the network envelope methods provide by TLS is relevant, but you still want to protect and validate the body
of the data.

Second, if you cannot rely on TLS. For example, you have a man-in-the-middle (mitm) proxy whose certificate is on your endpoint,
removing all of your protection. Or perhaps your organization requires you to work through a proxy that prevents TLS entirely.

Shibboleth helps you build applications that provide the two most basic protections in the message body itself.

## What It Does

Shibboleth either wraps or unwraps raw data in a package. That package is encrypted, signed or both.

The package format is JSON for universal readability and simplicity. The contents of your original message,
of course, can be anything you want.

## How To Use It

Shibboleth comes in two forms:

* CLI: creates and reads Shibboleth messages
* Library: golang library to do the same

For the library, see the [godoc](https://godoc.org/github.com/zededa/shibboleth).

The CLI has two main commands:

```console
$ echo cleartest_message | shibboleth wrap
$ echo wrapped_message | shibboleth unwrap
```

You need to pass the CLI parameters for:

For both wrapping and unwrapping:

* where to find your private key to wrap or unwrap
* where to find your public key to wrap or unwrap

For wrapping:

* where to find the certificate to use to sign your public key; optional
* whether to include your certificate/public key in the actual message, or just the hash of it

For unwrapping:

* where to find a list of certificates that you already trust, if any; optional
* where to find a list of certificate authorities you trust, to validate a certificate passed; optional
* whether to trust the CAs installed in your system

Run `shibboleth --help` to get any of the options.

## How It Works

Shibboleth uses [Elliptic-Curve](https://en.wikipedia.org/wiki/Elliptic-curve_cryptography) keys. Both endpoints, Sender and Receiver,
must already have public keys that are trusted by the other endpoint. These can be via certificates that are signed by a trusted CA,
or simply having the keys.

Shibboleth then uses [Elliptic-Curve Diffie Hellman](https://en.wikipedia.org/wiki/Elliptic-curve_Diffie–Hellman) (ECDH) to exchange
encryption keys. Notably, the keys are ephemeral, and are valid only for each communication. The processes follow below for each party.

Sender

1. Creates the original unprotected data to send.
1. Generates an ephemeral ECC key-pair
1. OPTIONAL: Signs the ECC public key with the well-known certificate. As stated above, the signer certificate must either be shared in advance with the Receiver, or have been signed by a CA trusted by the Receiver.
1. Uses the Receiver's public key, the Sender's private key and the pre-agreed parameters to generate a session key that will be used only for this message.
1. Encrypts the payload with the session key
1. Hashes the payload
1. Signs the hash of the payload with the Sender's private key
1. Constructs the Shibboleth message

The Shibboleth message is as follows:

```json
{
  "payload": {
    "body": "", // body of the message, the encrypted payload, base64-encoded
    "signature": "", // signature of the hash of the payload
    "algo": "", // the algorithm used
   },
  "signer": {
    "body": "", // the Sender's signer, either the actual ECC public key, or a signed certificate, PEM-encoded; optional
    "hash": "", // the hash of the Sender's signer body, base64-encoded
    "type": "", // the type of the Sender's signer body, one of: "C" (certificate), "K" (key)
  }
}
```

To save on space, the message will be minimized to eliminate unnecessary newlines and whitespace.

The Sender has the option to include its signer - the ECC public key or a signed certificate - in the message, or not. It MUST include the hash of 
the signer body. This enables the Receiver to know whether or not it can validate the message. If the Receiver has a hash of a signer
that matches the hash, then it knows the signer; if it does not, it does not. If the signer is included and is a certificate, the Receiver can
validate the certificate in whatever way it normally would.

If the signer is not included, and the Receiver does not have a signer whose hash matches that of the signer
used to sign the payload, then the Receiver cannot use the payload, unless it has other channels to retrieve an appropriate signer.
Retrieving that signer is beyond the scope of Shibboleth.

Receiver

1. Receives the message
1. Checks if it has a signer it trusts, whose hash matches the field `signer.hash`. If it does, it can continue. If it does not:
   * Retrieve the signer in the `signer.body` field and try to validate it. If it cannot be validated, nothing further can be done.
1. Use the Sender's public key in the signer, the Receiver's private key and the pre-agreed parameters to regenerate the session key for this message.
1. Use the session key to decrypt the message.
1. Use the Sender's public key in the signer to validate the hash of the payload.

The algorithms supported are identical to those supported by [TLS 1.3](https://tools.ietf.org/html/rfc8446#appendix-B.4). The `TLS_` prefix
is kept, even though this is not TLS, for simplicity and consistency's sake.

These are:

* `TLS_AES_256_GCM_SHA384`
* `TLS_CHACHA20_POLY1305_SHA256`
* `TLS_AES_128_GCM_SHA256`
* `TLS_AES_128_CCM_8_SHA256`
* `TLS_AES_128_CCM_SHA256`


