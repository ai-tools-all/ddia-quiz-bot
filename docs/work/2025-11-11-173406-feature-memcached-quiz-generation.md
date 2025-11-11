# Memcached Quiz Generation

**Date:** 2025-11-11 17:34:06
**Category:** feature
**Status:** in-progress

## Objective
Generate MCQ and subjective quiz questions for MIT 6.824 Lecture 16: Cache Consistency - Memcached at Facebook

## Source Material
- Transcript: `transcripts/mit-6824-subtitles/016-Lecture_16_-_Cache_Consistency_-_Memcached_at_Facebook-summary.md`

## Reference Materials
- Prompts: `prompts/001-quiz-generator.md`, `prompts/001-quiz-flow.md`
- Example: `ddia-quiz-bot/content/chapters/14-mit-6824-optimistic-cc/`

## Tasks
- [ ] Analyze transcript and identify key concepts
- [ ] Generate 6 MCQ questions
- [ ] Generate subjective questions for L3-L7 (2 questions per level)
- [ ] Create folder structure: `ddia-quiz-bot/content/chapters/16-mit-6824-memcached/`
- [ ] Save MCQ files
- [ ] Save subjective files
- [ ] Commit and push changes

## Key Concepts from Transcript
1. Look-aside caching architecture
2. Write invalidation protocol (delete vs set)
3. Regional replication (primary/secondary regions)
4. Intra-region clustering strategy
5. Partitioning vs replication trade-offs
6. Read-your-writes consistency
7. MySQL async replication with memcached
8. Lease mechanism and thundering herd
9. Incast congestion and connection management
10. Cold cluster warmup

## Notes
- MCQ format: YAML frontmatter with id, day, tags + question/options/answer/explanation/hook
- Subjective format: YAML frontmatter + structured sections (main_question, core_concepts, sample answers, follow-ups)
- Difficulty progression: L3-L7 with increasing complexity
