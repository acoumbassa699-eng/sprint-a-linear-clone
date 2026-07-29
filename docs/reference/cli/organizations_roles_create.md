<!-- DO NOT EDIT | GENERATED CONTENT -->
# organizations roles create

Create a new organization custom role

## Usage

```console
optimus-ide-collab organizations roles create [flags] <role_name>
```

## Description

```console
  - Run with an input.json file:

     $ optimus-ide-collab organization -O <organization_name> roles create --stdin < role.json
```

## Options

### -y, --yes

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Bypass confirmation prompts.

### --dry-run

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Does all the work, but does not submit the final updated role.

### --stdin

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Reads stdin for the json role definition to upload.
