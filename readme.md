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
  -f, --no-fetch        Do not fetch before verifying repository status
  -n, --no-verify       Do not check repository status before creating tags
  -p, --prefix string   Prefix for the version tag
  -q, --quiet           Quiet - only output errors

Use "bump [command] --help" for more information about a command.
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
