---
id: network-timeout-vs-refused-diagnosis
day: 16
tags: [networking, troubleshooting, tcp, debugging, practical]
related_stories:
  - network-basics
  - troubleshooting
---

# Connection Timeout vs Connection Refused

## question
Your application tries to connect to a database at 10.0.1.50:5432. You get "Connection refused" immediately. What does this tell you about the problem?

## options
- A) A firewall is blocking port 5432
- B) The host 10.0.1.50 is reachable, but nothing is listening on port 5432
- C) The network route to 10.0.1.50 is down
- D) The database is overloaded and rejecting new connections

## answer
B

## explanation
"Connection refused" means the TCP SYN packet reached the destination host successfully, but the host's OS responded with RST (reset) because no process is listening on that port. This is instant and definitive. In contrast: (A) a firewall would typically cause a timeout (SYN packets dropped silently) or send back ICMP "port unreachable"; (C) routing issues cause timeouts (SYN retransmits until timeout); (D) an overloaded database would either timeout (connection queue full) or accept-then-close, but not immediately refuse. "Connection refused" specifically means: host is up, network path works, but the port is closed. Check if the service is running, listening on the correct port/interface, or if you have the wrong port number.

## hook
What would you see if a firewall was blocking the connection instead?
