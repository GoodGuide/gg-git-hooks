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

In the simple case, just clone this repo to a shared location, then replace your local repo's `.git/hooks` with a symlink to this repo:

```shell
[ -d ~/.git-hooks ] || git clone git@github.com:GoodGuide/git-hooks.git ~/.git-hooks
rm -rf .git/hooks
ln -s ~/.git-hooks .git/hooks
```

## Requirements

- Some features lean on the newest version of Git: 1.9.0; you can update to this version with Homebrew on OSX)
- You'll need to set up your Pivotal API token in git config. [Get your API Token here][pivotal-account-settings], then:

    ```
    $ git config pivotal.api-token [YOUR_TOKEN]
    ```

    You probably want to use the `--global` option as well, which sets the value in your `~/.gitconfig` as opposed to the current repository's `.git/config`.

- Ruby 1.9+ should be available.

- If you have [Selecta][] available on your PATH, it will be used to offer incremental search of an available story by name or ID, and will automatically tag the commit with the story chosen.

[pivotal-account-settings]: https://www.pivotaltracker.com/profile#api
[Selecta]: https://github.com/garybernhardt/selecta
