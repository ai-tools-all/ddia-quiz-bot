### Security Context & Certificate Authority Problem
* Man-in-the-middle attacks pre-1995 allowed attackers to intercept DNS responses and impersonate legitimate servers like Gmail, capturing user credentials without detection.
* SSL/TLS certificates with public-private key pairs were introduced, requiring servers to prove identity by signing challenges with their private key, verified against public keys certified by certificate authorities.
* Certificate authorities (CAs) can be compromised, sloppy, or malicious—issuing bogus certificates for domains they don't own—enabling sophisticated attacks that bypass traditional certificate checks.

### Certificate Transparency Solution
* Certificate Transparency forces all certificates to be published in public append-only logs that anyone can audit, making rogue certificates visible to domain owners and monitors.
* Browsers require certificates to include a Signed Certificate Timestamp (SCT) proving the certificate was added to a public log before accepting it.
* Monitors continuously scan logs for unexpected certificates for domains they care about, alerting owners to potential attacks or misissuance.

### Merkle Tree Structure & Signed Tree Heads
* Logs use Merkle trees where leaf nodes contain certificate hashes, and internal nodes hash pairs of child nodes up to a single root called the Signed Tree Head (STH).
* The log server signs each STH, creating a compact cryptographic commitment to the entire log contents at that point in time.
* Merkle trees enable logarithmic-size inclusion proofs: to prove certificate X is at position i, the log provides only O(log n) sibling hashes along the path to the root.

### Fork Consistency Guarantees
* The critical threat is a malicious log server showing different log contents to browsers versus monitors—browsers see bogus certificates while monitors see a clean log.
* Cryptographic hash collision resistance prevents the log from producing valid inclusion proofs for different certificates at the same position under the same STH.
* Fork consistency ensures that if a log shows inconsistent views to different parties, gossiping STHs will reveal the discrepancy through failed consistency proofs.

### Consistency Proofs & Detection
* Given two STHs, a consistency proof demonstrates that the log corresponding to the earlier STH is a prefix of the log corresponding to the later STH.
* If STHs cannot be connected by a valid consistency proof despite being from the same log server, this proves the log has forked and is showing different views.
* The system remains secure even with untrusted log servers because cryptographic proofs make misbehavior detectable through gossip among participants.

### Gossip Protocol for Fork Detection
* Browsers and monitors exchange STHs they've received from log servers, comparing them to detect inconsistencies.
* When comparing two STHs, participants request consistency proofs from the log server to verify one is a legitimate extension of the other.
* Failed consistency proofs or incompatible STHs expose log server misbehavior, triggering alerts and potential revocation of the log's trusted status.

### Key Takeaways
* Certificate Transparency addresses CA compromise through transparency rather than trust, making all certificate issuance publicly auditable without requiring universal trust.
* Merkle trees provide efficient cryptographic proofs of log contents and consistency, enabling browsers and monitors to verify logs contain expected certificates.
* Fork consistency guarantees that malicious log servers cannot show different views to browsers and monitors without detection, as gossiped STHs will reveal inconsistencies through failed consistency proofs that prove the log has diverged.
