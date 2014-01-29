<!--
  Please don't hard wrap at 80 for this file:
  Vim: set wrap linebreak formatoptions-=tc tabstop=2 softtabstop=2 shiftwidth=2:
-->

# GoodGuide Git hooks

This is a collection of a few useful hooks for Git, for use across all projects.

They're all intended to be project agnostic; configuration happens via Git's config.

Each is documented well, so please read the code to understand more.

## How to install

In the simple case, just replace your local repo's `.git/hooks` with a clone of this repository:

```shell
rm -rf .git/hooks
git clone git@github.com:GoodGuide/git-hooks.git .git/hooks
```

## Requirements

At least one of the hooks requires the following #! to work:

```
#!/usr/bin/env ruby
```

(This is Ruby 1.8+ compatible code)
