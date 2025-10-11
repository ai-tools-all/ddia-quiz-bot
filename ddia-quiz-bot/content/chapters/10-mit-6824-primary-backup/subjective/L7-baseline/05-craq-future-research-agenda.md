---
id: craq-subjective-L7-005
type: subjective
level: L7
category: baseline
topic: craq
subtopic: research-agenda
estimated_time: 12-15 minutes
---

# question_title - CRAQ Future Research Agenda

## main_question - Core Question
"Define a research agenda for the next decade that pushes CRAQ beyond current limits—consider hardware innovation, formal verification, and new consistency paradigms. Tie priorities to DDIA's themes of correctness, evolution, and real-world impact." 

## core_concepts - Must Mention (60%)
- **Hardware Exploration**: Persistent memory, NIC offload, programmable switches to accelerate clean propagation
- **Formal Methods**: Model checking/Coq proofs for CRAQ protocols, automated regression testing
- **Consistency Innovations**: Hybrid consistency models (bounded staleness, adaptive linearizability) for new workloads
- **Impact Evaluation**: Metrics and benchmarks demonstrating business value and reliability improvements

### expected_keywords
- Primary keywords: research agenda, hardware acceleration, formal verification, consistency paradigm
- Technical terms: RDMA, P4 switches, TLA+, model checking, bounded staleness

## peripheral_concepts - Nice to Have (40%)
- **Academic Partnerships**: Collaborate with universities, publish papers
- **Open Data Sets**: Benchmark contributions, reproducible research
- **Ethical Considerations**: Responsible innovation, socio-technical impact
- **Commercialization Path**: Transfer research to product roadmap

### bonus_keywords
- Implementation: SmartNIC, NVMe-oF, SMT solver, adaptive quorum, research KPIs
- Scenarios: financial trading, IoT edge, AI inference pipelines
- Trade-offs: research investment vs near-term ROI

## sample_excellent - Example Excellence
"Invest in hardware acceleration—prototype CRAQ on SmartNICs and programmable switches to shave propagation latency. Launch a formal verification track using TLA+ to prove clean/dirty protocol correctness, integrating tests into CI. Explore adaptive consistency modes that switch between linearizable and bounded staleness when workloads permit, aligning with DDIA's call for practical trade-offs. Publish benchmarks and partner with academia, ensuring research results feed into product roadmaps with clear success metrics (latency reduction, reliability improvements)." 

## sample_acceptable - Minimum Acceptable
"Focus on hardware acceleration, formal verification, and adaptive consistency research, measuring improvements and partnering with academia to turn findings into product features—matching DDIA's themes of correctness and evolution." 

## common_mistakes - Watch Out For
- Laundry-list ideas without strategic prioritization
- No link to DDIA or product impact
- Ignoring path from research to production
- Lack of metrics or proof goals

## follow_up_excellent - Depth Probe
**Question**: "How would you stage funding and milestones to balance exploratory work with accountable deliverables?"
- **Looking for**: Stage-gate process, portfolio approach, KPIs
- **Red flags**: Open-ended research without accountability

## follow_up_partial - Guided Probe  
**Question**: "Which workloads benefit most from adaptive consistency, and how will you test them?"
- **Hint embedded**: Trading, gaming, analytics—A/B testing, workload replay
- **Concept testing**: Application of research

## follow_up_weak - Foundation Check
**Question**: "Why plan both near-term experiments and long-term moonshots in R&D?"
- **Simplification**: Portfolio of bets analogy
- **Building block**: Balance risk and reward

## bar_raiser_question - L7→L8 Challenge
"Lead a cross-industry research consortium to standardize proofs and benchmarks for strongly consistent read-scaled systems." 

### bar_raiser_concepts
- Research coordination
- Industry benchmarks, open science
- Influence and standardization
- Long-term vision execution

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 6-8 min answer + 6-7 min discussion
- **Common next topics**: Research leadership, hardware acceleration, formal verification
