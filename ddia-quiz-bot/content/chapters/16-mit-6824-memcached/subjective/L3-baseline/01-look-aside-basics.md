---
id: memcached-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: memcached
subtopic: look-aside-caching
estimated_time: 5-7 minutes
---

# question_title - Look-Aside Caching Basics

## main_question - Core Question
"Explain Facebook's look-aside caching pattern. What are the client's responsibilities on reads and writes, and how does this differ from look-through caching?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- On reads: check memcached first, fetch from MySQL on miss, populate cache
- On writes: send to MySQL, then delete affected keys from memcached
- Front-end (client) manages the cache-database relationship
- Differs from look-through where cache itself fetches from backing store

### expected_keywords
- look-aside, cache miss, invalidation, front-end, MySQL, memcached

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Allows caching arbitrary transformations, not just raw rows
- Decoupling enables flexible caching strategies
- 99%+ of reads served from cache

### bonus_keywords
- decoupling, transformations, read-heavy workload

## sample_excellent - Example Excellence
"In look-aside caching, the front-end (PHP code) checks memcached first on reads. On a cache miss, it fetches data from MySQL, applies any transformations, and populates memcached. On writes, the front-end updates MySQL then deletes affected keys from memcached, forcing subsequent reads to fetch fresh data. This differs from look-through caching where the cache layer itself would handle database fetches. Look-aside decouples the cache from the database, allowing front-ends to cache transformed data like aggregations or joins, not just raw database rows."

## sample_acceptable - Minimum Acceptable
"Front-ends check memcached on reads and fetch from MySQL on misses. On writes, they update MySQL and delete from memcached. The client manages the cache-database relationship, unlike look-through where the cache does this."

## common_mistakes - Watch Out For
- Confusing look-aside with look-through semantics
- Not mentioning invalidation on writes
- Missing that front-end controls caching logic

## follow_up_excellent - Depth Probe
**Question**: "Why is decoupling the cache from the database beneficial for Facebook's use case?"
- **Looking for**: Flexibility to cache computed results, joins, application-specific structures

## follow_up_partial - Guided Probe
**Question**: "What happens on a write if you forget to invalidate memcached?"
- **Hint embedded**: Subsequent reads would return stale cached data

## follow_up_weak - Foundation Check
**Question**: "What are the two main operations a client performs in look-aside caching?"
