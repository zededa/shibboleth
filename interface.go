package shibboleth

type Shibboleth interface{
	// Wrap take a message, wrap it with encryption and signing, returning the json-encoded bytes to send.
	// Parameters include:
	//  - message: the message to prepare as payload
	//  - senderPubKey: the sender's ECDSA public key
	//  - senderPrivKey: the sender's ECDSA private key
	//  - signer: the certificate to use to sign the public key, can be blank, in which case the public key is sent as-is
	//  - includeCertificate: whether to include the certificate in the message, or just the hash of the certificate
  Wrap(message io.Reader, senderPubKey crypto.PublicKey, senderPrivKey crypto.PrivateKey, signer *x509.Certificate, includeCertificate bool) ([]byte, error)
	// Unwrap take a wrapped message, unwrap and verify it, returning the original message.
	// Parameters include:
	//  - message: the wrapped message
	//  - receiverPubKey: the receiver's ECDS public key
	//  - receiverPrivKey: the receiver's ECDSA private key
	//  - certificates: a slice of valid certificates to use to validate the sent hash
	//  - authorities: a slice of certificate authorities to use to validate the sent certificate, if included in a message
	//  - systemAuthorities: whether or not to use the system CAs, in addition to any passed to `authorities`, to validate the sent certificate
  Unwrap(message io.Reader, receiverPubKey crypto.PublicKey, receiverPrivKey crypto.PrivateKey, certificates, authorities []x509.Certificate, systemAuthorities bool) ([]byte, error)
}
