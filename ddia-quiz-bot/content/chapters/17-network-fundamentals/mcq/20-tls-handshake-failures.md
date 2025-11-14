---
id: network-tls-handshake-failure-diagnosis
day: 20
tags: [networking, tls, ssl, security, troubleshooting, practical]
related_stories:
  - network-basics
  - security
---

# TLS Handshake Failures

## question
A client fails to establish a TLS connection with error "certificate verify failed: unable to get local issuer certificate". The server certificate is valid and properly configured. What is the most likely cause?

## options
- A) The server's private key doesn't match its certificate
- B) The client doesn't have the root CA certificate that signed the server's certificate in its trust store
- C) The server is using an outdated TLS protocol version
- D) A firewall is blocking TLS traffic

## answer
B

## explanation
During TLS handshake, the client validates the server's certificate chain: server cert → intermediate CA cert(s) → root CA cert. The error "unable to get local issuer certificate" means the client can't find the root CA (or intermediate CA) certificate in its local trust store to complete the chain. This commonly happens when: (1) using self-signed certificates without importing them, (2) the server doesn't send intermediate certificates (only sends leaf certificate), (3) client's CA bundle is outdated or incomplete, (4) using internal/corporate CAs without distributing their root cert. The fix: either add the missing CA certificate to the client's trust store, or ensure the server sends the complete certificate chain including intermediates. Option A would cause a different error ("certificate signature failure"). Option C triggers "protocol version" errors. Option D causes connection timeouts, not certificate errors.

## hook
Why do servers need to send intermediate CA certificates along with their leaf certificate?
