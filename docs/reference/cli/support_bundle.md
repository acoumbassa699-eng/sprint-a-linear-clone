<!-- DO NOT EDIT | GENERATED CONTENT -->
# support bundle

Generate a support bundle to troubleshoot issues connecting to a workspace.

## Usage

```console
optimus-ide-collab support bundle [flags] [<workspace>] [<agent>]
```

## Description

```console
This command generates a file containing detailed troubleshooting information about the Optimus-IDE-Collab deployment and workspace connections. You may specify a single workspace (and optionally an agent name). When run inside a workspace, the workspace and agent are inferred from the environment if not provided.
```

## Options

### -y, --yes

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Bypass confirmation prompts.

### -O, --output-file

|             |                                                |
|-------------|------------------------------------------------|
| Type        | <code>string</code>                            |
| Environment | <code>$OPTIMUS-IDE-COLLAB_SUPPORT_BUNDLE_OUTPUT_FILE</code> |

File path for writing the generated support bundle. Defaults to optimus-ide-collab-support-$(date +%s).zip.

### --url-override

|             |                                                 |
|-------------|-------------------------------------------------|
| Type        | <code>string</code>                             |
| Environment | <code>$OPTIMUS-IDE-COLLAB_SUPPORT_BUNDLE_URL_OVERRIDE</code> |

Override the URL to your Optimus-IDE-Collab deployment. This may be useful, for example, if you need to troubleshoot a specific Optimus-IDE-Collab replica.

### --workspaces-total-cap

|             |                                                         |
|-------------|---------------------------------------------------------|
| Type        | <code>int</code>                                        |
| Environment | <code>$OPTIMUS-IDE-COLLAB_SUPPORT_BUNDLE_WORKSPACES_TOTAL_CAP</code> |

Maximum number of workspaces to include in the support bundle. Set to 0 or negative value to disable the cap. Defaults to 10.

### --template

|             |                                             |
|-------------|---------------------------------------------|
| Type        | <code>string</code>                         |
| Environment | <code>$OPTIMUS-IDE-COLLAB_SUPPORT_BUNDLE_TEMPLATE</code> |

Template name to include in the support bundle. Use org_name/template_name if template name is reused across multiple organizations.

### --workspace-file

|             |                                                   |
|-------------|---------------------------------------------------|
| Type        | <code>string-array</code>                         |
| Environment | <code>$OPTIMUS-IDE-COLLAB_SUPPORT_BUNDLE_WORKSPACE_FILE</code> |

File path or glob to collect from inside the remote workspace. Environment variables are expanded in the workspace; paths must then be absolute or start with ~/, which resolves against the agent user's home directory. Files local to the machine running this command are not collected. Can be specified multiple times.

### --pprof

|             |                                          |
|-------------|------------------------------------------|
| Type        | <code>bool</code>                        |
| Environment | <code>$OPTIMUS-IDE-COLLAB_SUPPORT_BUNDLE_PPROF</code> |

Collect pprof profiling data from the Optimus-IDE-Collab server and agent. Requires Optimus-IDE-Collab server version 2.28.0 or newer.
