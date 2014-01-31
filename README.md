<!--
  Please don't hard wrap at 80 for this file:
  Vim: set wrap linebreak formatoptions-=tc tabstop=2 softtabstop=2 shiftwidth=2:
-->

# GoodGuide Git hooks

This is a collection of a few useful hooks for Git, for use across all projects.

They're all intended to be project agnostic; configuration happens via Git's config.

Each is documented well, so please read the code to understand more.

## Hooks Provided

- *pre-commit* &ndash; Git runs this right before setting up for you to enter your message. This hook can exit non-zero, which aborts the commit. The implementation in this repository has a couple checks for white-space errors, as well as a check for accidental committing of a `binding.pry` or `debugger`.

- *prepare-commit-msg* &ndash; Git runs this after the pre-commit hook, and it is able to modify the template message that it gives to your `$GIT_EDITOR`. The implementation in this repo will augment the default message with a commented-out list of your active stories in Pivotal Tracker.

- *commit-msg* &ndash; Git runs this, passing it the commit message you gave it, and the hook can abort the commit if it doesn't meet certain criteria. The implementation in this repo verifies that you have tagged your commit with a Tracker story ID.

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

(This is Ruby 1.9+ compatible code)
