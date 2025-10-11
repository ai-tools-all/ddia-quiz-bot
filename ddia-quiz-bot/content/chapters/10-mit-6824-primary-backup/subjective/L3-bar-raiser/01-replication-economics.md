---
id: primary-backup-subjective-L3-BR-001
type: subjective
level: L3
category: bar-raiser
topic: primary-backup
subtopic: economics-tradeoffs
estimated_time: 7-10 minutes
---

# question_title - Economics of Replication Decisions

## main_question - Core Question
"A startup is deciding whether to implement primary-backup replication for their service. The backup server costs $20,000/year. Their service generates $50,000/day in revenue. Industry data shows servers fail once every 2 years on average, taking 4 hours to repair. Should they implement replication? What factors beyond simple math should influence this decision?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **Downtime Cost Calculation**: Lost revenue during failures
- **Replication Cost**: Hardware, network, operational overhead
- **Failure Frequency**: Mean time between failures (MTBF)
- **Business Impact**: Beyond direct revenue loss

### expected_keywords
- Primary keywords: cost-benefit, downtime, MTBF, availability
- Business terms: revenue loss, customer trust, SLA, reputation

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Customer Lifetime Value**: Lost customers from outages
- **Reputation Damage**: Hard to quantify but significant
- **SLA Penalties**: Contractual obligations
- **Operational Complexity**: Maintenance overhead of replication
- **Peak Time Failures**: Not all downtime hours equal

### bonus_keywords
- Implementation: gradual degradation, partial failures, cascade effects
- Business: competitive advantage, market position, customer acquisition cost
- Technical debt: complexity cost, testing overhead, upgrade challenges

## sample_excellent - Example Excellence
"Based on pure math: Expected annual downtime = 0.5 failures/year × 4 hours = 2 hours. Lost revenue = 2 hours × $2,083/hour = $4,166/year. This is far less than $20,000 replication cost, suggesting no replication. However, this oversimplifies reality. Critical factors include: (1) Customer trust - one bad outage might lose customers worth far more than $4,166, (2) Revenue timing - failure during Black Friday could cost millions, (3) Growth trajectory - early-stage reputation damage could limit future growth, (4) Competitive landscape - if competitors have better availability, customers will switch, (5) SLA commitments - might face penalties exceeding direct revenue loss, (6) Recovery time variance - '4 hours average' might mean some failures take days. For a startup generating $18M annually, $20,000 for replication is just 0.1% of revenue - likely worthwhile for the peace of mind and customer trust alone. The decision depends more on business strategy than simple ROI calculation."

## sample_acceptable - Minimum Acceptable
"The math shows: 0.5 failures per year × 4 hours × $2,083/hour = $4,166 annual loss, which is less than $20,000 replication cost. But other factors matter: customer trust lost during outages, reputation damage, possible SLA penalties, and that some failures might take much longer than 4 hours. For a business making $18M/year, spending $20,000 on reliability seems reasonable to protect customer relationships."

## common_mistakes - Watch Out For
- Only considering direct revenue loss
- Ignoring customer acquisition costs
- Assuming all failures are equal
- Not considering growth impact
- Missing operational overhead of replication

## follow_up_excellent - Depth Probe
**Question**: "This startup is B2B with 10 enterprise customers, each paying $5,000/day. How does this change your analysis compared to B2C with 10,000 customers paying $5 each?"
- **Looking for**: Concentration risk, enterprise SLAs, switching costs
- **Red flags**: Not recognizing different customer dynamics

## follow_up_partial - Guided Probe  
**Question**: "What if the 4-hour average includes one 24-hour outage every 10 years? How does this change the calculation?"
- **Hint embedded**: Tail risks matter more than averages
- **Concept testing**: Understanding risk distribution

## follow_up_weak - Foundation Check
**Question**: "Your personal laptop crashes once every 2 years. Would you pay $100/year for automatic backup that restores everything in 1 hour instead of 1 day?"
- **Simplification**: Personal analogy to business decision
- **Building block**: Value of availability varies by context

## bar_raiser_question - L3→L4 Challenge
"Two scenarios: (A) E-commerce site with consistent $50K daily revenue, (B) Tax software making $30K/day normally but $500K/day in the week before tax deadline. Both have same failure statistics. How should their replication strategies differ? Consider cost structures and customer behavior."

### bar_raiser_concepts
- Temporal revenue concentration
- Customer behavior patterns
- Seasonal availability requirements
- Dynamic scaling possibilities
- Risk tolerance by time period
- Hybrid strategies

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 3-4 min answer + 4-6 min discussion
- **Common next topics**: SLA design, availability targets, cost optimization
