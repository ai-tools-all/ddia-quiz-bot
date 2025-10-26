---
id: fault-tolerance-subjective-L6-001
type: subjective
level: L6
category: baseline
topic: fault-tolerance
subtopic: byzantine-faults
estimated_time: 8-10 minutes
---

# question_title - Byzantine Fault Tolerance vs Crash Fault Tolerance

## main_question - Core Question
"Compare and contrast Byzantine Fault Tolerant (BFT) consensus with crash fault tolerant protocols like Raft. When would you choose BFT despite its significantly higher complexity and cost? Design a hybrid system that uses both approaches effectively."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Fault Model Differences**: Crash-stop vs arbitrary/malicious behavior
- **3f+1 vs 2f+1**: BFT needs more replicas for same fault tolerance
- **Performance Impact**: Higher latency, message complexity in BFT
- **Trust Boundaries**: When participants can't trust each other

### expected_keywords
- Primary keywords: Byzantine, crash fault, trust, malicious, integrity
- Technical terms: PBFT, Tendermint, proof of work, state machine replication

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Hybrid Approaches**: BFT across organizations, CFT within
- **Cryptographic Overhead**: Digital signatures, MACs, merkle trees
- **View Changes**: More complex in BFT due to Byzantine primary
- **Blockchain Consensus**: BFT in permissionless settings
- **Economic Incentives**: Aligning behavior without trust
- **Hardware Security**: SGX, HSMs as trust anchors

### bonus_keywords
- Protocols: PBFT, HotStuff, Tendermint, Algorand
- Attacks: equivocation, censorship, DOS
- Optimizations: threshold signatures, speculation

## sample_excellent - Example Excellence
"Byzantine Fault Tolerance assumes nodes can exhibit arbitrary behavior - lying, colluding, or sending conflicting messages - while crash fault tolerance (CFT) like Raft assumes nodes either work correctly or stop. This fundamental difference drives dramatic complexity increases in BFT:

**Resource Requirements**: BFT needs 3f+1 nodes to tolerate f failures (need 2f+1 correct nodes for majority among non-faulty), while CFT needs only 2f+1. BFT requires cryptographic signatures on all messages, adding computational overhead and message size.

**Protocol Complexity**: BFT protocols like PBFT require multiple phases (pre-prepare, prepare, commit) with different quorum intersection properties. View changes are complex because the old primary might be Byzantine and actively disrupting. CFT protocols like Raft are simpler because failed nodes can't interfere.

**When to Use BFT**: 
1. Cross-organization systems where participants have conflicting interests (supply chain, inter-bank settlement)
2. Public/permissionless networks (cryptocurrencies)
3. High-stakes systems where insider threats are concerns (military, critical infrastructure)

**Hybrid Design**: I'd architect a multi-tier system:
- **Inter-organization layer**: BFT consensus (e.g., Tendermint) across company boundaries for critical shared state
- **Intra-organization layer**: Raft clusters within each company for internal services
- **Bridge components**: Translate between layers, enforce policy
- **Optimization**: Use BFT only for conflict-prone operations; route trusted operations through faster CFT path

This provides Byzantine resistance where needed while maintaining performance for trusted operations. Example: A supply chain system uses BFT for ownership transfers between companies but CFT for inventory updates within each company."

## sample_acceptable - Minimum Acceptable
"Byzantine fault tolerance handles nodes that can lie or act maliciously, while crash fault tolerance like Raft only handles nodes that stop working. BFT requires 3f+1 nodes versus 2f+1 for CFT, needs cryptographic signatures, and has more complex protocols with multiple rounds of messages. Use BFT when participants don't trust each other, like in cryptocurrencies or cross-company systems. A hybrid approach could use BFT between organizations but Raft within each organization's systems for better performance."

## common_mistakes - Watch Out For
- Thinking BFT is always "better" (it has major costs)
- Not understanding the 3f+1 requirement derivation
- Ignoring cryptographic overhead
- Missing trust boundary analysis

## follow_up_excellent - Depth Probe
**Question**: "In your hybrid design, how do you handle consistency between the BFT and CFT layers? What happens if the CFT layer commits something the BFT layer later rejects?"
- **Looking for**: Two-phase commit, compensation, eventual consistency strategies
- **Red flags**: Not seeing the consistency challenge

## follow_up_partial - Guided Probe  
**Question**: "You mentioned BFT needs 3f+1 nodes. Walk through why 2f+1 isn't sufficient when nodes can lie."
- **Hint embedded**: Byzantine nodes can send different messages to different nodes
- **Concept testing**: Understanding quorum intersection with Byzantine nodes

## follow_up_weak - Foundation Check
**Question**: "Imagine you're organizing a vote where some people might lie about the results. How many honest vote counters would you need?"
- **Simplification**: Majority of honest nodes
- **Building block**: Trust and verification

## bar_raiser_question - L6→L7 Challenge
"Design a consensus protocol for a federated system where different organizations run validator nodes with varying trust levels and computational resources. How do you weight votes while maintaining Byzantine fault tolerance?"

### bar_raiser_concepts
- Weighted voting with trust scores
- Adaptive protocols based on threat detection  
- Economic incentives and slashing conditions
- Heterogeneous fault models per participant
- Performance vs security trade-offs

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 5-7 min answer + 3-3 min discussion
- **Common next topics**: Blockchain consensus, threshold cryptography, secure multiparty computation

## assistant_answer
BFT tolerates arbitrary/malicious behavior (needs 3f+1 replicas, cryptographic auth, multi-phase commits), while Raft tolerates crash faults with 2f+1 and simpler protocols. Choose BFT across untrusted parties (cross-org, public chains); a hybrid uses BFT between organizations and Raft within each organization, with bridges and selective anchoring of critical state.

## improvement_suggestions
- Ask for a specific BFT protocol choice (PBFT/HotStuff/Tendermint) and view-change handling details.
- Require analysis of cryptographic overhead and the use of signature aggregation/threshold signatures.

## improvement_exercises
### exercise_1 - Protocol Choice and View Change
**Question**: "Pick PBFT, HotStuff, or Tendermint for an inter-org consortium. Justify the choice and outline its view/leader change process."

**Sample answer**: "HotStuff: linear message complexity with chained pipelining; simpler view change by carrying QC (quorum certificate) to the next leader. PBFT: 3-phase with complex view change; Tendermint: lock/unlock on rounds with proposer rotation. HotStuff balances simplicity and throughput for WAN."

### exercise_2 - Crypto Overhead and Aggregation
**Question**: "Estimate the cost of per-message signatures at 10k TPS and explain how threshold signatures or BLS aggregation help."

**Sample answer**: "At 10k TPS with 3 phases and 7 validators, naive signing/verifying saturates CPU. BLS aggregation combines many signatures into one, reducing verification cost and bandwidth; threshold signatures allow a single combined signature once ≥t shares collected, lowering message size and verification work."
