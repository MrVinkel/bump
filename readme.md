# BUMP

Bump those versions!

Utility for bumping and pushing git tags.

## Usage

```txt
Bump those versions! Utility for bumping and pushing git tags

Usage:
  bump [flags]
  bump [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  major       Bump the major version
  minor       Bump the minor version
  patch       Bump the patch version
  version     Print the version of bump

Flags:
  -d, --debug           Debug mode
  -x, --dry-run         Do not create tags, only print what would be done
  -h, --help            help for bump
  -c, --no-commit       Do not commit changes to the repository
  -f, --no-fetch        Do not fetch before verifying repository status
  -n, --no-verify       Do not check repository status before creating tags
  -p, --prefix string   Prefix for the version tag
  -q, --quiet           Quiet - only output errors
  -s, --skip-pre-hook   Skip any configured pre-hook

Use "bump [command] --help" for more information about a command.
```

## Config

Create a `.bump.json` in the root of the repository will enforce `bump` settings and gives the ability to configure a pre-hook which should run before the tagging. The pre-hook can create changes in files which will then be committed and pushed, before creating the tag.

Example config:

```json
{
  "$schema": "https://raw.githubusercontent.com/mrvinkel/bump/1.0.0/bump.schema.json",
  // Default commit message
  "message": "release ${VERSION}",
  // Enforce prefix
  "prefix": "v",
  // Enforce commit
  "commit": true,
  // Enforce fetch
  "fetch": true,
  // Default shell command
  "shell": "/bin/bash -c",
  // Pre-hooks runs in the shell and have access to the new and previous version env vars
  "preHook": [
    "echo $VERSION",
    "echo $PREVIOUS_VERSION"
  ]
}
```

## SSH agent

Bump requires a SSH agent to be running when using SSH for auth.

Add the following to .bashrc or similar.

```bash
# start and export ssh-agent env vars
eval $(ssh-agent)
# add private key
ssh-add ${HOME}/.ssh/id_ed25519
```
