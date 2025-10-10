---
id: ch04-textual-encoding-problems
day: 2
tags: [encoding, json, xml, numbers, ambiguity]
related_stories: []
---

# Problems with Textual Encoding

## question
What is a common problem when encoding numbers in JSON and XML?

## options
- A) They cannot encode floating-point numbers
- B) Ambiguity in distinguishing numbers from strings representing numbers
- C) They only support integers up to 32-bit
- D) Numbers cannot be negative

## answer
B

## explanation
JSON and XML do not distinguish between numbers and strings that happen to contain numeric digits. An application might encode a number as a string, and another application reading it might interpret it incorrectly. Additionally, JSON doesn't distinguish between integers and floating-point numbers, and doesn't specify precision.

## hook
How does your API handle large integers that exceed JavaScript's Number precision?
