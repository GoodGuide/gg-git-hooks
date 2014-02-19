# The pattern used to check a commit message contains a Pivotal Tracker story ID
TAG_PATTERN = /
  \[
    (?:
      (?:
       (?:complete[sd]?|(?:finish|fix)(?:e[sd])?)\s+
      )?\#\d{4,}
    |
      no[ ]story
    )
  \]
/ix
