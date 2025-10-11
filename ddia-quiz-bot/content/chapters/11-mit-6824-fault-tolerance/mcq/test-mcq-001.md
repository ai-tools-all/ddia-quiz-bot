+++
id = "fault-tolerance-mcq-L3-001"
title = "Primary-Backup Failure Types"
level = "L3"
category = "baseline"
type = "mcq"
+++

## Question

What type of failures can primary-backup replication typically handle?

## Options

- A) Software bugs that cause incorrect calculations
- B) Hardware failures that cause the server to stop executing
- C) Network partitions that split the system
- D) All of the above

## Answer

B

## Explanation

Primary-backup replication handles fail-stop failures where servers stop cleanly. It does not handle software bugs (which would affect both primary and backup identically) or Byzantine failures. Network partitions require additional mechanisms like quorum systems.

## Hook

Understanding the limitations of replication is crucial for designing robust distributed systems.

## Core Concepts

- Fail-stop failures
- Replication limitations
- Fault models

## Peripheral Concepts

- Byzantine failures
- Network partitions
- Quorum systems
