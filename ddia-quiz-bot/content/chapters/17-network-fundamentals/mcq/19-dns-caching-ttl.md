---
id: network-dns-caching-ttl-issues
day: 19
tags: [networking, dns, caching, ttl, troubleshooting, practical]
related_stories:
  - network-basics
  - dns-fundamentals
---

# DNS Caching and TTL Issues

## question
You updated your application's DNS record from IP 1.2.3.4 to 5.6.7.8 at 10:00 AM. The TTL is 3600 seconds (1 hour). At 10:30 AM, some users report they can still reach the old IP. Why?

## options
- A) DNS propagation takes 24-48 hours globally
- B) Users or their resolvers cached the old DNS record before your update, and the TTL hasn't expired yet
- C) Your DNS provider hasn't updated all their servers yet
- D) Browsers cache DNS indefinitely regardless of TTL

## answer
B

## explanation
DNS caching is distributed across multiple layers: user's browser, OS resolver cache, ISP's DNS resolver, intermediate caches. If a user's resolver queried your DNS record at 9:30 AM (30 minutes before your change), it cached the old IP (1.2.3.4) with TTL 3600 seconds (valid until 10:30 AM). At 10:00 AM you update the record, but users with cached entries won't see it until their cache expires. At 10:30 AM, caches from 9:30 AM are still valid. The old record won't fully expire from all caches until 11:00 AM (1 hour after your change). This is why pre-planning matters: lower TTL to 60 seconds hours before your change, make the change, then raise TTL back up after confirming propagation. The "24-48 hour propagation" myth is outdated - modern DNS respects TTL properly.

## hook
What strategies can you use to minimize DNS-related downtime during infrastructure changes?
