---
id: zookeeper-subjective-L5-002
type: subjective
level: L5
category: baseline
topic: zookeeper
subtopic: production-deployment
estimated_time: 10-12 minutes
---

# question_title - Production Zookeeper Patterns and Best Practices

## main_question - Core Question
"Your company is deploying Zookeeper to production for critical services. Design a production-ready deployment covering high availability, disaster recovery, monitoring, and operational procedures. What are the key risks and how do you mitigate them?"

## core_concepts - Must Mention (60%)
### mandatory_concepts
- **HA Setup**: Odd number of nodes, multi-rack/AZ deployment
- **Backup Strategy**: Snapshots and transaction logs
- **Monitoring**: Key metrics and alerting thresholds
- **Operational Procedures**: Rolling upgrades, failure recovery

### expected_keywords
- Primary keywords: availability, disaster recovery, monitoring, operations
- Technical terms: quorum, split brain, backup, restore, metrics

## peripheral_concepts - Nice to Have (40%)
### bonus_concepts
- **Security Hardening**: ACLs, SSL, authentication, audit logs
- **Capacity Management**: Growth planning, limits, quotas
- **Automation**: Infrastructure as code, automated recovery
- **Testing**: Chaos engineering, load testing, failure injection
- **Documentation**: Runbooks, architecture diagrams, SLAs
- **Compliance**: Data residency, encryption, retention policies
- **Performance Tuning**: JVM settings, OS tuning, network optimization

### bonus_keywords
- Tools: Ansible, Terraform, Prometheus, Grafana, PagerDuty
- Practices: GitOps, blue-green deployment, canary releases
- Standards: SOC2, GDPR, HIPAA compliance requirements

## sample_excellent - Example Excellence
"Here's a comprehensive production deployment strategy:

Infrastructure Design:
```
Production: 5 nodes (2-2-1 across 3 AZs)
├── zk-prod-1a (AZ-1) - Leader eligible
├── zk-prod-1b (AZ-1) - Follower
├── zk-prod-2a (AZ-2) - Follower
├── zk-prod-2b (AZ-2) - Follower
└── zk-prod-3a (AZ-3) - Follower

DR Site: 5 nodes (separate region, async replication)
```

High Availability:
- **Quorum**: Can lose 2 nodes and maintain service
- **Anti-affinity**: No two nodes on same physical host
- **Network**: Redundant network paths between nodes
- **Storage**: RAID-10 for data directories, separate disk for logs

Disaster Recovery:
1. **Backups**: Hourly snapshots to S3 with 30-day retention
2. **Transaction Logs**: Ship to remote storage every 5 minutes
3. **Recovery Time Objective**: 1 hour for full restoration
4. **Recovery Point Objective**: Maximum 5 minutes data loss
5. **DR Testing**: Monthly failover drills

Monitoring Stack:
```yaml
Critical Alerts (Page immediately):
- Quorum loss
- Leader election failures > 2
- Disk usage > 90%
- Connection count > 80% limit

Warning Alerts (Ticket):
- Latency p99 > 100ms
- Memory usage > 75%
- Snapshot failures
- Certificate expiry < 30 days

Dashboards:
- Real-time: Operations/sec, latency, connections
- Historical: Growth trends, capacity planning
- Business: SLA adherence, incident metrics
```

Operational Procedures:

Rolling Upgrade Process:
1. Test in staging environment
2. Backup current state
3. Upgrade observers first
4. Upgrade followers one at a time
5. Upgrade leader last (triggers election)
6. Verify ensemble health between each step
7. Automated rollback on failure

Security Measures:
- **Network**: Private subnet, security groups, NACLs
- **Authentication**: Kerberos for service accounts
- **Encryption**: TLS 1.3 for client connections
- **ACLs**: Least privilege per service
- **Audit**: All admin operations logged to SIEM

Risk Mitigation:

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Hardware failure | High | Low | 5-node ensemble, hot spares |
| Network partition | Medium | High | Multi-AZ, redundant networks |
| Data corruption | Low | Critical | Checksums, backups, replicas |
| Operator error | Medium | High | Automation, peer review, rollback |
| Capacity exceeded | Medium | Medium | Monitoring, auto-scaling observers |

Automation:
```bash
# Infrastructure as Code (Terraform)
module "zookeeper" {
  source = "./modules/zookeeper"
  environment = "production"
  instance_count = 5
  enable_backups = true
  enable_monitoring = true
}

# Configuration Management (Ansible)
ansible-playbook -i production zookeeper-deploy.yml
```

Performance Tuning:
- JVM: G1GC, 4GB heap, optimize for low latency
- OS: Increase file descriptors, disable swap
- Network: Jumbo frames for cross-DC traffic
- Storage: Separate disks for snapshots and logs

Documentation:
- Architecture diagrams with failure scenarios
- Runbooks for common operations
- Incident response procedures
- Performance baselines and limits

This approach ensures 99.95% availability with comprehensive operational safety."

## sample_acceptable - Minimum Acceptable
"Deploy 5 Zookeeper nodes across multiple availability zones for high availability. Set up regular backups of snapshots and transaction logs. Monitor key metrics like quorum health, latency, and disk usage with appropriate alerts. Document procedures for upgrades and failure recovery. Use security best practices like encryption and authentication."

## common_mistakes - Watch Out For
- Even number of nodes (split brain risk)
- No backup strategy
- Insufficient monitoring coverage
- No documented procedures
- Ignoring security considerations

## follow_up_excellent - Depth Probe
**Question**: "How would you perform a zero-downtime migration from an old Zookeeper cluster to a new one with different hardware?"
- **Looking for**: Dual-write strategy, data verification, gradual cutover, rollback plan
- **Red flags**: Big bang cutover, no verification step

## follow_up_partial - Guided Probe  
**Question**: "What would happen if your backup site loses connectivity to the primary site? How do you prevent split brain?"
- **Hint embedded**: Quorum-based decisions prevent split brain
- **Concept testing**: Understanding distributed system failures

## follow_up_weak - Foundation Check
**Question**: "If you were running a critical website, what would you monitor to know it's healthy?"
- **Simplification**: Basic monitoring concepts
- **Building block**: Metrics and alerting fundamentals

## bar_raiser_question - L5→L6 Challenge
"Design an automated system that can detect and recover from Byzantine failures in your Zookeeper cluster, including corrupted nodes sending valid-looking but incorrect data. How do you detect, isolate, and recover?"

### bar_raiser_concepts
- Byzantine fault detection
- Checksum verification strategies
- Automated quarantine procedures
- Forensic analysis tools
- Self-healing mechanisms
- Trust verification protocols

## metadata
- **Evaluation rubric**: See GUIDELINES.md
- **Time expectation**: 4-5 min answer + 5-7 min discussion
- **Common next topics**: Site reliability engineering, chaos engineering, distributed system operations
