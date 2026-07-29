<!-- DO NOT EDIT | GENERATED CONTENT -->
# task resume

Resume a task

## Usage

```console
optimus-ide-collab task resume [flags] <task>
```

## Description

```console
  - Resume a task by name:

     $ optimus-ide-collab task resume my-task

  - Resume another user's task:

     $ optimus-ide-collab task resume alice/my-task

  - Resume a task without confirmation:

     $ optimus-ide-collab task resume my-task --yes
```

## Options

### --no-wait

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Return immediately after resuming the task.

### -y, --yes

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Bypass confirmation prompts.
