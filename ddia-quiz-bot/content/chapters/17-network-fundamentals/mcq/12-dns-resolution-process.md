---
id: network-dns-recursive-iterative
day: 12
tags: [networking, dns, resolution, recursive, iterative]
related_stories:
  - network-basics
  - dns-fundamentals
---

# DNS Resolution: Recursive vs Iterative

## question
Your browser queries your ISP's DNS resolver for "api.example.com". The resolver doesn't have it cached. What type of queries does the resolver typically make to find the answer?

## options
- A) Recursive queries to root servers, TLD servers, and authoritative servers
- B) Iterative queries to root servers, TLD servers, and authoritative servers
- C) Only a single query to the authoritative nameserver for example.com
- D) Broadcast queries to all DNS servers simultaneously

## answer
B

## explanation
DNS resolvers typically perform iterative resolution when querying the DNS hierarchy. The resolver asks the root server "who handles .com?" and gets a referral to .com TLD servers. It then asks .com TLD "who handles example.com?" and gets a referral to example.com's authoritative servers. Finally, it asks the authoritative server for "api.example.com" and gets the actual IP. Each query returns either an answer or a referral (iterative). In contrast, recursive queries expect the queried server to do all the work and return a final answer. Client-to-resolver queries are typically recursive (client expects a complete answer), but resolver-to-hierarchy queries are iterative (each server provides next step).

## hook
Why do DNS resolvers cache responses, and how does TTL affect caching behavior?
