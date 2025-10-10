---
id: ch04-schema-evolution-vs-versioning
day: 19
tags: [schema-evolution, versioning, api-versioning]
related_stories: []
---

# Schema Evolution vs API Versioning

## question
What is the difference between schema evolution and API versioning?

## options
- A) They are the same thing with different names
- B) Schema evolution allows changes within a version while maintaining compatibility; API versioning creates new incompatible versions
- C) Schema evolution is only for databases, API versioning is only for web services
- D) Schema evolution is faster than API versioning

## answer
B

## explanation
Schema evolution allows you to make backward and forward compatible changes within a single version (e.g., adding optional fields), avoiding the complexity of maintaining multiple API versions. API versioning creates explicitly different versions (v1, v2) when you need to make breaking changes. Schema evolution is preferable when possible since supporting multiple versions increases maintenance burden.

## hook
How many API versions are you maintaining in production right now?
