#!/usr/bin/env ruby

# Run this _from_ the git working tree of a repo on which you'd like to install
# these hooks. It will reconcile paths and install symlinks to the git-hooks
# repo checkout

require 'pathname'
require 'fileutils'

GIT_HOOKS_DIR = File.expand_path(File.dirname(__FILE__))

system 'git rev-parse'
raise 'This must be run from a git repo' unless $?.success?

def find_git_dir
  if File.basename(Dir.pwd) == '.git'
    git_dir = Dir.pwd
  else
    git_root = File.expand_path(`git rev-parse --show-cdup`.chomp)
    git_dir = File.join(git_root, '.git')
  end
  Pathname.new(git_dir)
end

def any_not_sample?(pathnames)
  pathnames.any? do |p|
    p.basename.to_s != '.' and p.basename.to_s != '..' and p.extname != '.sample'
  end
end

def symlink_is_correct?(pathname)
  return unless pathname.symlink?
  File.expand_path(pathname.readlink) == GIT_HOOKS_DIR
end

symlink_path = find_git_dir.join('hooks')

if symlink_is_correct?(symlink_path)
  warn "Already installed!"
  exit 0

elsif symlink_path.directory? and any_not_sample?(symlink_path.entries)
  warn ".git/hooks is a directory and has non-sample hooks installed already. Aborting!"
  exit 1


elsif symlink_path.exist? and !symlink_path.symlink? and !symlink_path.directory?
  warn ".git/hooks exists and isn't a symlink or directory. Aborting!"
  exit 1
end

FileUtils.rm_rf(symlink_path.to_s)

File.symlink(GIT_HOOKS_DIR, symlink_path)
