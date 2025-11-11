---
id: memcached-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: memcached
subtopic: write-invalidation
estimated_time: 5-7 minutes
---

# question_title - Write Invalidation Protocol

## main_question - Core Question
"Describe Facebook's write invalidation protocol. Why does the system use deletes instead of updating the cached value directly?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- Front-end sends write to MySQL first
- Front-end immediately issues delete to memcached for affected keys
- MySQL servers also send deletes asynchronously by monitoring replication log
- Delete is preferred over set because concurrent writes can leave stale values

### expected_keywords
- invalidation, delete, MySQL, concurrent writes, replication log, stale data

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Front-end deletes ensure read-your-writes consistency
- Async MySQL deletes provide redundancy if front-end deletes fail
- Order of operations matters for consistency

### bonus_keywords
- asynchronous, redundancy, ordering, consistency guarantee

## sample_excellent - Example Excellence
"On writes, front-ends first update MySQL, then immediately delete affected keys from memcached. This invalidation forces subsequent reads to miss cache and fetch fresh data from the database. Deletes are preferred over cache updates (sets) because concurrent writes could cause consistency issues—if two clients' sets arrive out of order relative to their database commits, a stale value could remain cached. Additionally, MySQL servers asynchronously monitor the replication log and send deletes, ensuring invalidations propagate even if front-end deletes are lost due to network issues."

## sample_acceptable - Minimum Acceptable
"Write to MySQL then delete from memcached. Delete is used instead of update because concurrent writes can cause stale cached data if sets arrive in the wrong order."

## common_mistakes - Watch Out For
- Not explaining why delete is better than set
- Missing the dual invalidation (front-end + MySQL async)
- Not connecting invalidation to consistency

## follow_up_excellent - Depth Probe
**Question**: "What consistency guarantee does front-end invalidation specifically provide?"
- **Looking for**: Read-your-writes consistency—users see their own updates immediately

## follow_up_partial - Guided Probe
**Question**: "If front-end delete fails but MySQL async delete succeeds, what happens?"
- **Hint embedded**: User might briefly see stale data until async delete propagates

## follow_up_weak - Foundation Check
**Question**: "What are the two sources of delete operations in the system?"
