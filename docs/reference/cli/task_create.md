<!-- DO NOT EDIT | GENERATED CONTENT -->
# task create

Create a task

## Usage

```console
optimus-ide-collab task create [flags] [input]
```

## Description

```console
  - Create a task with direct input:

     $ optimus-ide-collab task create "Add authentication to the user service"

  - Create a task with stdin input:

     $ echo "Add authentication to the user service" | optimus-ide-collab task create

  - Create a task with a specific name:

     $ optimus-ide-collab task create --name task1 "Add authentication to the user service"

  - Create a task from a specific template / preset:

     $ optimus-ide-collab task create --template backend-dev --preset "My Preset" "Add authentication to the user service"

  - Create a task for another user (requires appropriate permissions):

     $ optimus-ide-collab task create --owner user@example.com "Add authentication to the user service"
```

## Options

### --name

|      |                     |
|------|---------------------|
| Type | <code>string</code> |

Specify the name of the task. If you do not specify one, a name will be generated for you.

### --owner

|         |                     |
|---------|---------------------|
| Type    | <code>string</code> |
| Default | <code>me</code>     |

Specify the owner of the task. Defaults to the current user.

### --template

|             |                                        |
|-------------|----------------------------------------|
| Type        | <code>string</code>                    |
| Environment | <code>$OPTIMUS-IDE-COLLAB_TASK_TEMPLATE_NAME</code> |

### --template-version

|             |                                           |
|-------------|-------------------------------------------|
| Type        | <code>string</code>                       |
| Environment | <code>$OPTIMUS-IDE-COLLAB_TASK_TEMPLATE_VERSION</code> |

### --preset

|             |                                      |
|-------------|--------------------------------------|
| Type        | <code>string</code>                  |
| Environment | <code>$OPTIMUS-IDE-COLLAB_TASK_PRESET_NAME</code> |
| Default     | <code>none</code>                    |

### --stdin

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Reads from stdin for the task input.

### -q, --quiet

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Only display the created task's ID.

### -O, --org

|             |                                  |
|-------------|----------------------------------|
| Type        | <code>string</code>              |
| Environment | <code>$OPTIMUS-IDE-COLLAB_ORGANIZATION</code> |

Select which organization (uuid or name) to use.
