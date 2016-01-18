<!--
  Please don't hard wrap at 80 for this file:
  Vim: set wrap linebreak formatoptions-=tc tabstop=2 softtabstop=2 shiftwidth=2:
-->

# GoodGuide Git hooks

This is a collection of a few useful hooks for Git, for use across all projects.

They're all intended to be project agnostic; configuration happens via Git's config.

## Hooks Provided

- *pre-commit* &ndash; Git runs this right before setting up for you to enter your message. This hook can exit non-zero, which aborts the commit. The implementation in this repository has a couple checks for white-space errors, as well as a check for accidental committing of a `binding.pry` or `debugger`.

- *prepare-commit-msg* &ndash; Git runs this after the pre-commit hook, and it is able to modify the template message that it gives to your `$GIT_EDITOR`. The implementation in this repo will augment the default message with a commented-out list of your active stories in Pivotal Tracker.

- *commit-msg* &ndash; Git runs this, passing it the commit message you gave it, and the hook can abort the commit if it doesn't meet certain criteria. The implementation in this repo verifies that you have tagged your commit with a Tracker story ID, and offers you a dialog to select one if not

  ![Screenshot of story selection UI](http://f.cl.ly/items/0o3E3K0T2K05261y2g2g/Screen%20Shot%202014-09-03%20at%2010.58.25%20.png)

## Requirements

- Git `>= 1.9.0`
    - If you use Homebrew, you can update to this version easily: `brew update; brew upgrade git`

- You'll need to set up your Pivotal API token in git config. [Get your API Token here][pivotal-account-settings], then:

    ```shell
    $ git config --global pivotal.api-token [YOUR_TOKEN]
    ```

    The `--global` option sets the value in your `~/.gitconfig` as opposed to the current repository's `.git/config`.

## Installing

On a Mac with Homebrew, you may install via homebrew:

```bash
brew tap goodguide/tap
brew install goodguide-git-hooks
```

If you prefer to build yourself and have Go installed, you can simply `go get` the project.

```shell
go get github.com/goodguide/goodguide-git-hooks
```

Either way, there are [tagged releases with attached binaries on GitHub][releases]. Grab the build for your system, and just install the binary into your `~/.local/bin` or somewhere on your `PATH`.

## How to use

Once installed to your system, you can install the hooks to a particular local repo using the following from within the local repo in question:

```shell
goodguide-git-hooks install
```

Then, just use git normally.

### Migrating from previous Ruby version of these hooks

If you used the previous version of these tools, the recommended instructions were to make your repos' `.git/hooks` directory a symbolic link to the local clone of this repo. That is no longer necessary and that symlink should be deleted before you run `goodguide-git-hooks install`

### Integrating with existing hook logic

The `install` subcommand simply installs small shims into the `.git/hooks`
directory. For example, the `prepare-commit-message` shim looks like this:

```bash
#!/bin/bash
set -e

exec goodguide-git-hooks prepare-commit-message $@
```

You could easily just add the `goodguide-git-hooks CMD $@` command to your existing git hooks. (If you're not using `exec`, make sure you have `set -e` or manually check the exit status of this command so they can fail the commit if necessary.)

Similarly, you can add any additional logic to the generated shims. Existing hooks won't be rewritten by `goodguide-git-hooks install` unless you tell it to.

## Updating the cache of Tracker stories

The tracker-story-fetching is slow, and doesn't need to happen with every commit, as it did in the previous version of this project. The goal is to make it lazy but still automatic, but as of now it's extremely lazy, and manual. To fetch your stories manually:

```shell
goodguide-git-hooks update-pivotal-stories
```

## Development

To work on this project, you need Go installed and set up properly, then you should just be able to `go build` as usual. There are no external dependencies.

## Release process

1. To build the project and create a release, you'll need `goxc` installed:
    ```shell
    go get github.com/laher/goxc
    ```

2. Bump the version, commit the new version, and push that to github. Then create a tag based on the version and push that:
    ```shell
    make bump
    ```

3. Then, just run `goxc` to cross-compile for Linux/OSX and create tarballs in the `dist/` directory:
    ```shell
    make build
    ```

4. Then, go to the releases page on github, and edit the release you just made by pushing a tag. Add the contents of the `dist/` directory as individual binary attachments to the release.

5. Edit the `goodguide-git-hooks` formula in [goodguide/homebrew-tap](//github.com/goodguide/homebrew-tap) with the new version and SHA1 of the `darwin_amd64` archive.

[pivotal-account-settings]: https://www.pivotaltracker.com/profile#api
[releases]: //github.com/goodguide/goodguide-git-hooks/releases
