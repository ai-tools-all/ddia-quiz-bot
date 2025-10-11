---
id: primary-backup-subjective-L3-001
type: subjective
level: L3
category: baseline
topic: primary-backup
subtopic: fail-stop-failures
estimated_time: 5-7 minutes
---

# question_title - Understanding Fail-Stop Failures in Replication

## main_question - Core Question
"What are fail-stop failures, and why are they important in the context of primary-backup replication? Give examples of what types of failures can and cannot be handled by replication."

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Fail-Stop Definition**: Computer stops executing completely when failure occurs
- **No Incorrect Results**: System stops rather than computing wrong answers
- **Replication Scope**: Can only handle fail-stop failures effectively

### expected_keywords
- Primary keywords: fail-stop, crash, stop executing, availability
- Technical terms: power failure, network disconnection, hardware failure

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Cannot Handle Bugs**: Software bugs replicate across all copies
- **Hardware Design Defects**: Cannot handle systematic hardware issues
- **Independent Failures**: Assumes failures in replicas are uncorrelated
- **Checksums/Detection**: How some errors are converted to fail-stop

### bonus_keywords
- Implementation: checksum, error detection, error correction
- Examples: CPU overheating, kernel panic, network corruption
- Limitations: correlated failures, earthquakes, power grid failures

## sample_excellent - Example Excellence
"Fail-stop failures are failures where a computer simply stops executing when something goes wrong, rather than continuing to run and producing incorrect results. This is crucial for primary-backup replication because the system can detect when a replica has failed and switch to the backup. Examples that CAN be handled include power cable disconnection, network cable unplugging, or CPU overheating that causes clean shutdown. What CANNOT be handled are software bugs (which would affect both primary and backup identically), hardware design defects that cause incorrect computation, or correlated failures like earthquakes affecting all replicas in the same datacenter. The key assumption is that when something fails, it stops cleanly, allowing the system to detect the failure and failover to a working replica."

## sample_acceptable - Minimum Acceptable
"Fail-stop failures mean the computer completely stops working when there's a problem, rather than giving wrong answers. Replication can handle these - like when power is lost or network cable is unplugged. It cannot handle software bugs or hardware that computes incorrectly, since these would affect all replicas the same way."

## common_mistakes - Watch Out For
- Confusing fail-stop with Byzantine failures
- Thinking replication handles all types of failures
- Not understanding that bugs affect all replicas
- Missing the "clean stop" aspect

## follow_up_excellent - Depth Probe
**Question**: "How might a system convert potentially Byzantine failures (like bit flips in network packets) into fail-stop failures?"
- **Looking for**: Checksums, error detection codes, validation layers
- **Red flags**: Not understanding detection vs correction

## follow_up_partial - Guided Probe  
**Question**: "If we replicate a MapReduce master on two machines and there's a bug in the code, what happens?"
- **Hint embedded**: Both compute same wrong answer
- **Concept testing**: Understanding bug propagation

## follow_up_weak - Foundation Check
**Question**: "Imagine a server that sometimes gives wrong answers vs one that just stops working. Which is easier to handle with backup servers?"
- **Simplification**: Detectability of failures
- **Building block**: Why fail-stop is manageable

## bar_raiser_question - L3→L4 Challenge
"A company runs primary-backup replication with both servers in the same rack. They experience: (1) a software bug causing incorrect calculations, (2) rack power supply failure, (3) one server's fan failure. Which failures would the replication handle, and why?"

### bar_raiser_concepts
- Correlated vs independent failures
- Physical proximity risks
- Bug propagation across replicas
- Fail-stop conversion mechanisms

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 2-3 min answer + 3-4 min discussion
- **Common next topics**: State transfer, replicated state machine, VMware FT
