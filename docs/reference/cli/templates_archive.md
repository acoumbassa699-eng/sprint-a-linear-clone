<!-- DO NOT EDIT | GENERATED CONTENT -->
# templates archive

Archive unused or failed template versions from a given template(s)

## Usage

```console
optimus-ide-collab templates archive [flags] [template-name...] 
```

## Options

### -y, --yes

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Bypass confirmation prompts.

### --all

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Include all unused template versions. By default, only failed template versions are archived.

### -O, --org

|             |                                  |
|-------------|----------------------------------|
| Type        | <code>string</code>              |
| Environment | <code>$OPTIMUS-IDE-COLLAB_ORGANIZATION</code> |

Select which organization (uuid or name) to use.
