---
id: gfs-subjective-L3-bar
type: subjective
level: L3
category: bar-raiser
topic: gfs
subtopic: scalability
estimated_time: 7-10 minutes
---

# question_title - GFS Under Extreme Load

## main_question - Core Question
"It's Black Friday and your GFS cluster storing product images suddenly gets 100x normal traffic. The master server CPU hits 95%. What's happening and what would you do?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Master Bottleneck**: Single master handling all metadata requests
- **Client Caching**: Leverage client-side metadata caching
- **Read Replication**: Increase replicas for hot data
- **Load Distribution**: Spread load across chunk servers

### expected_keywords
- Primary keywords: master overload, metadata, caching, replicas
- Technical terms: hot spots, load balancing, bottleneck

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Emergency Response**: Immediate vs long-term solutions
- **Shadow Masters**: Read-only master replicas
- **Chunking Strategy**: Identify and replicate hot chunks
- **Application Changes**: CDN, caching layer
- **Monitoring**: Identify specific bottlenecks

### bonus_keywords
- Implementation: lease management, heartbeats
- Alternatives: Colossus multi-master, sharding
- Operations: traffic shaping, rate limiting

## sample_excellent - Example Excellence
"The master is overwhelmed with metadata requests as clients repeatedly ask for chunk locations. Immediate actions: 1) Increase client metadata cache TTL to reduce repeat requests, 2) Identify hot files (popular products) and increase their replication factor to spread read load, 3) Deploy shadow masters for read-only metadata operations. For the surge, consider adding a caching layer or CDN for popular images. Long-term, this suggests need for either multi-master architecture (like Colossus) or sharding the namespace. Monitor to identify if it's specific operations (opens vs reads) causing load."

## sample_acceptable - Minimum Acceptable
"The single master is getting too many metadata requests. I would increase client caching time for metadata and add more replicas for frequently accessed files to distribute the read load across more chunk servers."

## common_mistakes - Watch Out For
- Suggesting to scale the single master (can't easily scale)
- Missing that it's metadata not data traffic
- Proposing complete architecture changes during incident
- Not distinguishing immediate vs long-term solutions

## follow_up_excellent - Depth Probe
**Question**: "Your added replicas helped with reads, but now chunk servers are running out of network bandwidth. How do you handle this cascading problem?"
- **Looking for**: Rack-aware placement, traffic shaping, network topology optimization
- **Red flags**: Just adding more servers without considering network

## follow_up_partial - Guided Probe
**Question**: "You mentioned caching. The master is handling millions of requests per second. What specific information would you cache and where?"
- **Hint embedded**: Client-side vs server-side
- **Concept testing**: Understanding metadata types

## follow_up_weak - Foundation Check
**Question**: "Let's think about a simpler case. If a popular website is slow, what are different places you could add caches?"
- **Simplification**: Web architecture analogy
- **Building block**: Caching layers concept

## bar_raiser_question - L4→L5 Challenge
"Now design a system that can predict and prevent this Black Friday scenario. What monitoring, automation, and architectural changes would you implement?"

### bar_raiser_concepts
- Predictive scaling based on patterns
- Automated hot spot detection and mitigation
- Multi-tier architecture with degradation modes
- Cost optimization vs over-provisioning
- SLA definitions and guarantees

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-5 min discussion
- **Common next topics**: Multi-master architectures, CDN integration, capacity planning
