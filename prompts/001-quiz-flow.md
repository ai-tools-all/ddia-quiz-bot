# DDIA Level-Based Quiz Generation System

## Overview
This system generates progressive difficulty quiz questions based on "Designing Data-Intensive Applications" (DDIA) by Martin Kleppmann, aligned with engineering levels L4-L7.

---

## Prompt 1: Topic Taxonomy & Concept Mapping

```
You are a distributed systems curriculum architect with the expertise of Jeff Dean (Google), 
Leslie Lamport (Microsoft Research), and Martin Kleppmann (DDIA author). Your task is to 
create a comprehensive topic taxonomy from "Designing Data-Intensive Applications" (DDIA).

TASK:
Analyze DDIA and create a hierarchical breakdown of:

1. **Major Topics** (Chapters/Parts)
2. **Sub-Topics** (Sections within chapters)
3. **Core Concepts** (Individual concepts within sections)
4. **Interconnected Concepts** (Concepts that span multiple topics)

For each concept, identify:
- **Concept Name**: Clear identifier
- **Prerequisites**: What must be understood first
- **Related Concepts**: What connects to this concept
- **Difficulty Tier**: 
  - L4 (Mid-level SWE): Understands individual concepts, can apply in known scenarios
  - L5 (Senior SWE): Connects concepts, designs within constraints, anticipates trade-offs
  - L6 (Staff SWE): Synthesizes multiple systems, identifies subtle failure modes, architects complex systems
  - L7 (Principal+): Industry-wide perspective, invents novel solutions, sees second-order effects

OUTPUT FORMAT:
```yaml
topics:
  - name: "Data Models and Query Languages"
    subtopics:
      - name: "Relational vs Document Models"
        concepts:
          - name: "Schema flexibility vs enforced structure"
            prerequisites: ["Database basics", "JSON/XML formats"]
            related: ["Data locality", "Join operations", "Schema evolution"]
            difficulty_mapping:
              L4: "Explain differences between SQL and NoSQL"
              L5: "Choose appropriate model for use case with justification"
              L6: "Design hybrid approach, handle migration between models"
              L7: "Architect polyglot persistence strategy for organization"
            
          - name: "Impedance mismatch"
            prerequisites: ["Object-oriented programming", "Relational model"]
            related: ["ORMs", "Serialization", "Data modeling"]
            difficulty_mapping:
              L4: "Define impedance mismatch with examples"
              L5: "Evaluate ORM trade-offs for specific scenarios"
              L6: "Design data access layer minimizing mismatch"
              L7: "Architect cross-platform data representation strategy"
        
        interconnections:
          - concepts: ["Schema evolution", "Backward/forward compatibility"]
            insight: "Model choice affects evolution strategy"
            relevant_levels: [L5, L6, L7]
```

Create this taxonomy for ALL major DDIA topics including:
- Data Models & Query Languages
- Storage & Retrieval
- Encoding & Evolution
- Replication
- Partitioning
- Transactions
- Distributed Systems Troubles
- Consistency & Consensus
- Batch Processing
- Stream Processing

CONTEXT:
@book Designing Data-Intensive Applications by Martin Kleppmann
```

---

## Prompt 2: Question Generation Engine

```
You are a technical interviewer with the combined expertise of:
- **Jeff Dean** (Google): Large-scale systems design, performance optimization
- **Werner Vogels** (Amazon): Eventual consistency, service-oriented architecture
- **Leslie Lamport** (Microsoft): Distributed algorithms, formal methods
- **Martin Kleppmann** (DDIA): Data systems fundamentals, practical trade-offs

Using the topic taxonomy from Prompt 1, generate 5 questions for EACH difficulty level 
(L4, L5, L6, L7) for the specified topic.

INPUT:
- Topic: [TOPIC_NAME from taxonomy]
- Subtopic: [SUBTOPIC_NAME from taxonomy]
- Target Level: [L4/L5/L6/L7]

QUESTION TYPES BY LEVEL:

**L4 Questions (Mid-Level SWE):**
- **Focus**: Isolated concept understanding, definition, basic application
- **Cognitive Level**: Remember, Understand, Apply
- **Examples**: "What is...", "How does X work?", "Implement a basic..."
- **Expected Skills**: Can explain concepts, implement given specifications

**L5 Questions (Senior SWE):**
- **Focus**: Concept interconnections, trade-off analysis, design decisions
- **Cognitive Level**: Analyze, Evaluate
- **Examples**: "Compare X and Y...", "When would you choose...", "What are the trade-offs..."
- **Expected Skills**: Justifies choices, anticipates problems, designs within constraints

**L6 Questions (Staff SWE):**
- **Focus**: System synthesis, failure mode analysis, architectural patterns
- **Cognitive Level**: Evaluate, Create
- **Examples**: "Design a system that...", "How would you handle...", "What could go wrong if..."
- **Expected Skills**: Architects complex systems, identifies subtle issues, mentors others

**L7 Questions (Principal+):**
- **Focus**: Novel solutions, second-order effects, industry-wide perspective
- **Cognitive Level**: Create, Synthesize across domains
- **Examples**: "How would you evolve...", "Design a new approach to...", "What are the industry implications..."
- **Expected Skills**: Invents new patterns, sees organizational impact, influences direction

QUESTION DIFFICULTY PROGRESSION:
For each level, include:
- 2 "Comfort Zone" questions (at-level)
- 2 "Stretch" questions (one level up)
- 1 "Challenging" question (two levels up)

OUTPUT FORMAT:
```json
{
  "topic": "Replication",
  "subtopic": "Replication Lag and Consistency",
  "level": "L5",
  "questions": [
    {
      "id": "L5_REP_001",
      "difficulty": "comfort",
      "type": "trade-off_analysis",
      "question": "You're designing a social media feed service. Compare read-after-write consistency vs eventual consistency for this use case. What are the trade-offs in terms of user experience, system complexity, and scalability?",
      
      "expected_concepts": [
        "Read-after-write consistency",
        "Eventual consistency",
        "CAP theorem implications",
        "User experience considerations",
        "Scaling challenges"
      ],
      
      "key_insights": [
        "User's own posts should be immediately visible (read-after-write)",
        "Other users' posts can be eventually consistent",
        "Partition tolerance requirement affects consistency choices",
        "Cross-datacenter replication lag considerations"
      ],
      
      "evaluation_rubric": {
        "basic_understanding": {
          "points": 2,
          "criteria": "Defines both consistency models correctly"
        },
        "trade_off_analysis": {
          "points": 3,
          "criteria": "Identifies at least 3 trade-offs with justification"
        },
        "design_decision": {
          "points": 3,
          "criteria": "Makes clear choice with context-specific reasoning"
        },
        "edge_cases": {
          "points": 2,
          "criteria": "Considers failure scenarios or edge cases"
        }
      },
      
      "sample_answer": {
        "excellent": "For a social feed, I'd use hybrid consistency: read-after-write for a user's own posts (write to primary, read from primary for self) and eventual consistency for feed items from others. This balances UX (users see their actions immediately) with scalability (can serve feed from read replicas). Trade-offs: added complexity in routing logic, need to handle replication lag gracefully with 'this post may not be visible to others yet' indicators. In multi-datacenter setup, define 'primary' per-user region to minimize latency while maintaining consistency guarantees.",
        
        "good": "Read-after-write is better for user experience since people expect to see their own posts immediately. Eventual consistency is more scalable but can confuse users. I'd use read-after-write for write operations and eventual for reads from followers.",
        
        "needs_improvement": "Read-after-write means you see your writes. Eventual consistency means it takes time. Social media should use read-after-write because users want to see posts right away."
      },
      
      "follow_up_questions": [
        "How would you detect and communicate replication lag to users?",
        "What happens if the primary datacenter fails immediately after a user posts?",
        "How does this design scale to 1 billion users across 5 continents?"
      ],
      
      "references": {
        "ddia_chapters": ["Chapter 5: Replication"],
        "ddia_sections": ["Problems with Replication Lag", "Reading Your Own Writes"],
        "related_papers": ["Eventual Consistency - Werner Vogels"]
      }
    }
  ]
}
```

GENERATION RULES:
1. **Isolated Concepts** (40% of questions): Test single concept depth
2. **Interlinked Concepts** (40% of questions): Test concept relationships
3. **System Design** (20% of questions): Test holistic application

4. **Progressive Difficulty**:
   - Build on prerequisite knowledge
   - Each level up: add constraints, ambiguity, or scale
   - L7 questions should require synthesis across multiple chapters

5. **Real-World Grounding**:
   - Use realistic scenarios (e.g., "design a notification system", not "implement algorithm X")
   - Reference actual systems when relevant (Kafka, Cassandra, Spanner)
   - Include operational concerns (monitoring, debugging, evolution)

6. **Evaluation Criteria**:
   - Define clear rubric for each question
   - Identify key concepts that MUST appear in answer
   - Provide tiered sample answers (excellent/good/needs improvement)
   - Include follow-up questions to probe deeper

GENERATE QUESTIONS FOR:
Topic: [SPECIFY]
Subtopic: [SPECIFY]
Level: [L4/L5/L6/L7]
```

---

## Prompt 3: Quiz Assembly & Validation

```
You are a technical education specialist designing a comprehensive DDIA assessment. 
Using the questions generated from Prompt 2, assemble a complete quiz with proper 
scaffolding, validation, and learning pathways.

INPUT:
- Question bank from Prompt 2 (all levels, all topics)
- Target level: [L4/L5/L6/L7]
- Assessment type: [Diagnostic/Formative/Summative]

TASK:
Create a 25-question quiz (5 questions per major topic area) that:

1. **Validates Prerequisite Knowledge**
   - First 5 questions cover fundamentals (L4-level regardless of target)
   - Ensures quiz-taker has baseline understanding

2. **Progressive Difficulty Curve**
   - Questions 6-15: Mix of at-level and stretch questions
   - Questions 16-20: Primarily stretch and challenging questions
   - Questions 21-25: Challenging questions with synthesis requirements

3. **Balanced Coverage**
   - Each major DDIA topic represented
   - Mix of isolated concepts (40%), interlinked concepts (40%), system design (20%)
   - Variety of question formats (multiple choice, short answer, design problems)

4. **Learning Feedback**
   - For each question: explanation of correct answer
   - For each question: links to DDIA sections for further study
   - For incorrect answers: identify concept gaps and suggested review path

OUTPUT FORMAT:
```json
{
  "quiz_metadata": {
    "title": "DDIA Mastery Assessment - L5 (Senior SWE)",
    "target_level": "L5",
    "total_questions": 25,
    "estimated_time_minutes": 120,
    "passing_score": 70,
    "topics_covered": ["Replication", "Partitioning", "Transactions", "Consistency", "Batch Processing"],
    "prerequisite_knowledge": ["L4-level DDIA understanding", "Basic distributed systems"]
  },
  
  "questions": [
    {
      "question_number": 1,
      "section": "Baseline Assessment",
      "difficulty": "L4",
      "topic": "Data Models",
      "question": {...},
      "weight": 1.0
    },
    {
      "question_number": 6,
      "section": "Core Assessment",
      "difficulty": "L5",
      "topic": "Replication",
      "question": {...},
      "weight": 2.0
    },
    {
      "question_number": 21,
      "section": "Advanced Assessment",
      "difficulty": "L6",
      "topic": "Distributed Transactions",
      "question": {...},
      "weight": 3.0
    }
  ],
  
  "scoring_guide": {
    "L4_level_achievement": {
      "score_range": "0-60",
      "interpretation": "Solid grasp of fundamentals, ready to apply concepts in guided scenarios",
      "recommended_action": "Review interconnected concepts, focus on trade-off analysis"
    },
    "L5_level_achievement": {
      "score_range": "61-80",
      "interpretation": "Strong understanding with ability to make justified design decisions",
      "recommended_action": "Deepen understanding of failure modes, complex architectures"
    },
    "L6_level_achievement": {
      "score_range": "81-90",
      "interpretation": "Expert-level synthesis and architectural thinking",
      "recommended_action": "Explore novel solutions, contribute to field"
    },
    "L7_level_achievement": {
      "score_range": "91-100",
      "interpretation": "Mastery with industry-leading insights",
      "recommended_action": "Mentor others, publish findings, drive innovation"
    }
  },
  
  "learning_pathways": {
    "weak_in_replication": {
      "identified_by": "Score < 60% on replication questions",
      "review_materials": [
        "DDIA Chapter 5: Replication",
        "Focus on: Leader-follower replication, replication lag patterns"
      ],
      "practice_exercises": [
        "Design a multi-datacenter replication strategy",
        "Implement read-after-write consistency"
      ]
    }
  }
}
```

VALIDATION CHECKS:
- [ ] Question difficulty matches target level distribution
- [ ] No two questions test identical concepts
- [ ] Prerequisite concepts appear before dependent concepts
- [ ] Each major DDIA topic is adequately represented
- [ ] Sample answers are detailed and accurate
- [ ] Evaluation rubrics are objective and fair
- [ ] Time estimate is realistic for question complexity
- [ ] Learning pathways cover all major topics

GENERATE:
- Complete quiz for target level: [SPECIFY]
- Include all metadata, questions, scoring guides, and learning pathways
```

---

## Usage Instructions

### Sequential Execution (Recommended for Comprehensive Coverage)

1. **Run Prompt 1** → Get complete DDIA taxonomy
2. **Run Prompt 2** for each topic/level combination → Build question bank
   - Example: Run for "Replication" topic, L5 level
   - Repeat for all major topics and levels
3. **Run Prompt 3** with full question bank → Assemble final quiz

### Parallel Execution (Faster, for Specific Topics)

1. **Run Prompt 1** once → Get taxonomy
2. **Run Prompt 2** in parallel for different topics:
   - Thread 1: Replication questions (L4-L7)
   - Thread 2: Partitioning questions (L4-L7)
   - Thread 3: Transactions questions (L4-L7)
   - etc.
3. **Run Prompt 3** to assemble from parallel outputs

### Targeted Execution (Quick Assessment)

1. Skip Prompt 1 (use existing DDIA knowledge)
2. **Run Prompt 2** for specific topic + level only
3. **Run Prompt 3** with limited question set

---

## Customization Options

### For Different Books/Domains
- Replace DDIA with other technical books
- Adjust persona experts (e.g., for ML: Andrew Ng, Yann LeCun)
- Modify level definitions to match your organization

### For Different Assessment Goals
- **Diagnostic**: Focus on coverage, lighter depth
- **Formative**: Include hints, scaffolded learning
- **Summative**: No hints, strict evaluation

### For Different Formats
- **Multiple Choice**: Add distractors based on common misconceptions
- **Coding Challenges**: Include setup/skeleton code
- **Take-Home**: Extended time, research allowed

---

## Example Execution

```bash
# Step 1: Generate taxonomy
> Run Prompt 1 with @book DDIA

# Step 2: Generate questions for L5, Replication topic
> Run Prompt 2 with:
  Topic: "Replication"
  Subtopic: "Replication Lag and Consistency"
  Target Level: "L5"

# Step 3: Assemble quiz
> Run Prompt 3 with:
  Question bank: [output from Step 2]
  Target level: "L5"
  Assessment type: "Summative"
```

---

## Quality Assurance Checklist

- [ ] Questions are unambiguous and testable
- [ ] Sample answers demonstrate expected depth
- [ ] Evaluation rubrics are specific and measurable
- [ ] Difficulty progression is smooth and logical
- [ ] References to DDIA are accurate (chapter/section)
- [ ] Real-world scenarios are realistic and relevant
- [ ] Follow-up questions probe deeper understanding
- [ ] Learning pathways guide improvement effectively