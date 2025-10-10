---
id: ch06-partition-key-design
day: 14
tags: [partitioning, schema-design, best-practices]
related_stories: []
---

# Partition Key Design

## question
An e-commerce platform needs to partition order data. Which partition key would likely provide the best balance of load distribution and query efficiency?

## options
- A) order_timestamp
- B) customer_id
- C) product_category
- D) order_total_amount

## answer
B

## explanation
Customer_id provides good load distribution (assuming many customers) while keeping all of a customer's orders on the same partition, enabling efficient customer-specific queries which are common in e-commerce (order history, customer service, analytics). Timestamp would create hot partitions for recent data, category might be too coarse with skewed distribution, and order amount would scatter related orders making customer queries inefficient.

## hook
How does Amazon DynamoDB handle partition key selection for shopping cart data?
