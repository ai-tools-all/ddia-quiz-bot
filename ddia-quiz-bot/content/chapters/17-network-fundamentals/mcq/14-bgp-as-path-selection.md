---
id: network-bgp-as-path-selection
day: 14
tags: [networking, bgp, routing, autonomous-system, internet]
related_stories:
  - network-basics
  - internet-routing
---

# BGP AS Path Selection

## question
Your ISP (AS 1000) receives three BGP route advertisements for the same destination prefix 203.0.113.0/24: Route A via AS path [2000, 3000], Route B via AS path [4000], and Route C via AS path [5000, 6000, 7000]. Assuming all other BGP attributes are equal, which route will the ISP select?

## options
- A) Route A - the lowest AS number in the path
- B) Route B - the shortest AS path length
- C) Route C - the most specific path with detailed routing
- D) All three routes - BGP load-balances across equal-cost paths

## answer
B

## explanation
BGP's path selection algorithm prefers the shortest AS path (fewest autonomous systems to traverse) when other attributes are equal. Route B has AS path length 1, Route A has length 2, and Route C has length 3. Shorter paths generally mean fewer hops, lower latency, and reduced failure points. This is one of BGP's key decision criteria after checking local preference, AS path length comes early in the selection process (typically step 4 of 13). Note that BGP's decision process is more complex than just AS path length - it considers factors like local preference (higher wins), origin type, MED, and more - but AS path length is a fundamental tie-breaker.

## hook
How can organizations manipulate BGP path selection through AS path prepending, and why would they do this?
