<!-- DO NOT EDIT | GENERATED CONTENT -->
# dotfiles

Personalize your workspace by applying a canonical dotfiles repository

## Usage

```console
optimus-ide-collab dotfiles [flags] <git_repo_url>
```

## Description

```console
  - Check out and install a dotfiles repository without prompts:

     $ optimus-ide-collab dotfiles --yes git@github.com:example/dotfiles.git
```

## Options

### --symlink-dir

|             |                                 |
|-------------|---------------------------------|
| Type        | <code>string</code>             |
| Environment | <code>$OPTIMUS-IDE-COLLAB_SYMLINK_DIR</code> |

Specifies the directory for the dotfiles symlink destinations. If empty, will use $HOME.

### -b, --branch

|      |                     |
|------|---------------------|
| Type | <code>string</code> |

Specifies which branch to clone. If empty, will default to cloning the default branch or using the existing branch in the cloned repo on disk.

### --repo-dir

|             |                                       |
|-------------|---------------------------------------|
| Type        | <code>string</code>                   |
| Environment | <code>$OPTIMUS-IDE-COLLAB_DOTFILES_REPO_DIR</code> |
| Default     | <code>dotfiles</code>                 |

Specifies the directory for the dotfiles repository, relative to global config directory.

### -y, --yes

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Bypass confirmation prompts.
