
## working method

- create a markdown file in the issues dir and update with notes and tasks as work progresses in `docs/work` folder

- filename format: YYYY-MM-DD-[timestamp]-[category]-[four-word-summary].md
- to know today's date use  `date date +"%Y-%m-%d %H%M%S"` linux comand
- category can be `bug`, `feature`, `task`, `research` or anything suitable.
- numbered prefix is a running number starting from `001` (3 padded integer) for each day.


## building the project

`scripts/build.sh` to build the project 

## overall view of the project 

- this is a golang project. components of this projects are described below.

- the content is in folder: `ddia-quiz-bot/content/` 
  - containst the schedule to post questions in
  - mcq style questions + subjective questions as per level. 

- the tui for the mcq / subjective questions is in quiz-evaluator folder 

- spaced repetition is implemented for the MCQ style questions with persistence. 
  - how to create new questions / quiz ? - see the prompts/ folder on the general method used. verify that once by reading the folder `ddia-quiz-bot/content/chapters/10-mit-6824-primary-backup/subjective` (files in this fodler are large so do not read the entire file.) -
  - use the `./build/md-toc` binary to read the table of content of the markdown files in a directory. 