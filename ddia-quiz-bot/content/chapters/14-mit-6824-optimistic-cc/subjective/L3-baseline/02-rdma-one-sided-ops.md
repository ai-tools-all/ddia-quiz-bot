---
id: farm-subjective-L3-002
type: subjective
level: L3
category: baseline
topic: rdma
subtopic: one-sided-operations
estimated_time: 6-8 minutes
---

# question_title - RDMA One-Sided Operations

## main_question - Core Question
"Explain what one-sided RDMA operations are and why they are crucial for FaRM's performance. What constraint do they impose on the concurrency control design, and how does FaRM address this constraint?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- One-sided RDMA: client NIC directly reads/writes server memory without server CPU involvement
- Kernel bypass allows application-level network access
- Constraint: cannot check or acquire locks during reads since server CPU is not involved
- Solution: OCC defers all coordination to commit time when server CPUs can be involved

### expected_keywords
- RDMA, NIC, one-sided, server CPU, bypass, optimistic, validation, commit time

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- Sub-microsecond network latency with RDMA
- Compare to traditional RPC requiring server CPU
- Object IDs contain memory addresses for direct access
- Version numbers enable conflict detection without locks

### bonus_keywords
- latency, RPC, memory address, version, performance

## sample_excellent - Example Excellence
"One-sided RDMA operations let the client's NIC directly read/write server RAM without interrupting the server CPU, achieving sub-microsecond latencies via kernel bypass. This fundamentally prohibits traditional locking—clients can't acquire locks during reads because no server code executes. FaRM resolves this by using OCC: reads are optimistic (no locks), and at commit time, server CPUs validate via version checks and lock bits."

## sample_acceptable - Minimum Acceptable
"RDMA lets clients read server memory directly without involving the server CPU. This prevents using locks during reads, so FaRM uses optimistic concurrency control instead."

## common_mistakes - Watch Out For
- Confusing one-sided with two-sided RDMA operations
- Not explaining the server CPU bypass
- Missing the connection to OCC design choice
- Claiming RDMA eliminates all server CPU involvement (commit phase still uses it)

## follow_up_excellent - Depth Probe
**Question**: "What other design alternatives could work with one-sided RDMA besides OCC?"
- **Looking for**: timestamp-based concurrency control, version vectors, application-level conflict resolution

## follow_up_partial - Guided Probe
**Question**: "How does FaRM use server CPUs during the commit phase if RDMA bypasses them?"
- **Hint embedded**: LOCK and COMMIT messages are not one-sided operations

## follow_up_weak - Foundation Check
**Question**: "What is the main performance benefit of bypassing the server CPU during reads?"
