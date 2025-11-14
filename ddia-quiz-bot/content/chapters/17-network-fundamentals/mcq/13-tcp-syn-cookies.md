---
id: network-tcp-syn-cookies-security
day: 13
tags: [networking, tcp, security, syn-flood, syn-cookies]
related_stories:
  - network-basics
  - security
---

# TCP SYN Cookies and SYN Flood Defense

## question
A web server under SYN flood attack enables SYN cookies. How does this mechanism allow the server to accept legitimate connections without maintaining state for half-open connections?

## options
- A) The server allocates minimal memory for each SYN, reducing memory consumption
- B) The server encodes connection state in the initial sequence number and only allocates resources upon receiving the ACK
- C) The server blocks all incoming SYN packets until the attack subsides
- D) The server forwards SYN packets to a dedicated defense appliance that filters malicious traffic

## answer
B

## explanation
SYN cookies cleverly encode TCP connection parameters (MSS, timestamp, etc.) into the initial sequence number (ISN) sent in the SYN-ACK. The server doesn't allocate any connection state or memory for the SYN request. When the client sends the ACK with sequence number ISN+1, the server can decode the original parameters from ISN and verify the connection is legitimate. Only then does it allocate resources. This defeats SYN flood attacks because attackers send SYN packets with spoofed source IPs that never complete the handshake with an ACK. Without SYN cookies, the server would exhaust its connection queue with half-open connections. With SYN cookies, it maintains zero state until receiving the ACK, making it immune to SYN floods.

## hook
What are the trade-offs of using SYN cookies, and why aren't they enabled by default?
