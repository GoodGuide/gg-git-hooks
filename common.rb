# This is common things used by more than one of the git hooks

# The pattern used to check a commit message contains a Pivotal Tracker story ID
TAG_PATTERN = /
  \[
    (?:
      (?:
       (?:complete[sd]?|(?:finish|fix)(?:e[sd])?)\s+
      )?\#\d{4,}
    |
      \#\s*no[ _-]?(?:tracker|story)
    )
  \]
/ix


