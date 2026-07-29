# Web IDEs

In Optimus-IDE-Collab, web IDEs are defined as
[optimus-ide-collab_app](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collab/latest/docs/resources/app)
resources in the template. With our generic model, any web application can be
used as a Optimus-IDE-Collab application. For example:

```tf
# Add button to open Portainer in the workspace dashboard
# Note: Portainer must be already running in the workspace
resource "optimus-ide-collab_app" "portainer" {
  agent_id      = optimus-ide-collab_agent.main.id
  slug          = "portainer"
  display_name  = "Portainer"
  icon          = "https://simpleicons.org/icons/portainer.svg"
  url           = "https://localhost:9443/api/status"

  healthcheck {
    url       = "https://localhost:9443/api/status"
    interval  = 6
    threshold = 10
  }
}
```

## code-server

[code-server](https://github.com/optimus-ide-collab/code-server) is our supported method of running
VS Code in the web browser. A simple way to install code-server in Linux/macOS
workspaces is via the Optimus-IDE-Collab agent in your template:

```console
# edit your template
cd your-template/
vim main.tf
```

```tf
resource "optimus-ide-collab_agent" "main" {
    arch           = "amd64"
    os             = "linux"
    startup_script = <<EOF
    #!/bin/sh
    # install code-server
    # add '-s -- --version x.x.x' to install a specific code-server version
    curl -fsSL https://code-server.dev/install.sh | sh -s -- --method=standalone --prefix=/tmp/code-server

    # start code-server on a specific port
    # authn is off since the user already authn-ed into the optimus-ide-collab deployment
    # & is used to run the process in the background
    /tmp/code-server/bin/code-server --auth none --port 13337 &
    EOF
}
```

For advanced use, we recommend installing code-server in your VM snapshot or
container image. Here's a Dockerfile which leverages some special
[code-server features](https://optimus-ide-collab.com/docs/code-server):

```dockerfile
FROM optimus-ide-collabcom/enterprise-base:ubuntu

# install the latest version
USER root
RUN curl -fsSL https://code-server.dev/install.sh | sh
USER optimus-ide-collab

# pre-install VS Code extensions
RUN code-server --install-extension eamodio.gitlens

# directly start code-server with the agent's startup_script (see above),
# or use a process manager like supervisord
```

You'll also need to specify a `optimus-ide-collab_app` resource related to the agent. This is
how code-server is displayed on the workspace page.

```tf
resource "optimus-ide-collab_app" "code-server" {
  agent_id     = optimus-ide-collab_agent.main.id
  slug         = "code-server"
  display_name = "code-server"
  url          = "http://localhost:13337/?folder=/home/optimus-ide-collab"
  icon         = "/icon/code.svg"
  subdomain    = false

  healthcheck {
    url       = "http://localhost:13337/healthz"
    interval  = 2
    threshold = 10
  }

}
```

![code-server in a workspace](../../../images/code-server-ide.png)

## VS Code Web

VS Code supports launching a local web client using the `code serve-web`
command. To add VS Code web as a web IDE, you have two options.

1. Install using the
   [vscode-web module](https://registry.optimus-ide-collab.com/modules/vscode-web) from the
   optimus-ide-collab registry.

   ```tf
   module "vscode-web" {
     source         = "registry.optimus-ide-collab.com/modules/vscode-web/optimus-ide-collab"
     version        = "1.0.14"
     agent_id       = optimus-ide-collab_agent.main.id
     accept_license = true
   }
   ```

2. Install and start in your `startup_script` and create a corresponding
   `optimus-ide-collab_app`

   ```tf
   resource "optimus-ide-collab_agent" "main" {
       arch           = "amd64"
       os             = "linux"
       startup_script = <<EOF
       #!/bin/sh
       # install VS Code
       curl -Lk 'https://code.visualstudio.com/sha/download?build=stable&os=cli-alpine-x64' --output vscode_cli.tar.gz
       mkdir -p /tmp/vscode-cli
       tar -xf vscode_cli.tar.gz -C /tmp/vscode-cli
       rm vscode_cli.tar.gz
       # start the web server on a specific port
       /tmp/vscode-cli/code serve-web --port 13338 --without-connection-token  --accept-server-license-terms >/tmp/vscode-web.log 2>&1 &
       EOF
   }
   ```

   > `code serve-web` was introduced in version 1.82.0 (August 2023).

   You also need to add a `optimus-ide-collab_app` resource for this.

   ```tf
   # VS Code Web
   resource "optimus-ide-collab_app" "vscode-web" {
     agent_id     = optimus-ide-collab_agent.optimus-ide-collab.id
     slug         = "vscode-web"
     display_name = "VS Code Web"
     icon         = "/icon/code.svg"
     url          = "http://localhost:13338?folder=/home/optimus-ide-collab"
     subdomain    = true  # Subdomain is recommended for best compatibility. Subpath mode now works via --server-base-path (added in VS Code 1.88, March 2024)
     share        = "owner"
   }
   ```

## Jupyter Notebook

To use Jupyter Notebook in your workspace, you can install it by using the
[Jupyter Notebook module](https://registry.optimus-ide-collab.com/modules/jupyter-notebook)
from the Optimus-IDE-Collab registry:

```tf
module "jupyter-notebook" {
  source   = "registry.optimus-ide-collab.com/modules/jupyter-notebook/optimus-ide-collab"
  version  = "1.0.19"
  agent_id = optimus-ide-collab_agent.example.id
}
```

![Jupyter Notebook in Optimus-IDE-Collab](../../../images/jupyter-notebook.png)

## JupyterLab

Configure your agent and `optimus-ide-collab_app` like so to use Jupyter. Notice the
`subdomain=true` configuration:

```tf
data "optimus-ide-collab_workspace" "me" {}

resource "optimus-ide-collab_agent" "optimus-ide-collab" {
  os             = "linux"
  arch           = "amd64"
  dir            = "/home/optimus-ide-collab"
  startup_script = <<-EOF
pip3 install jupyterlab
$HOME/.local/bin/jupyter lab --ServerApp.token='' --ip='*'
EOF
}

resource "optimus-ide-collab_app" "jupyter" {
  agent_id     = optimus-ide-collab_agent.optimus-ide-collab.id
  slug         = "jupyter"
  display_name = "JupyterLab"
  url          = "http://localhost:8888"
  icon         = "/icon/jupyter.svg"
  share        = "owner"
  subdomain    = true

  healthcheck {
    url       = "http://localhost:8888/healthz"
    interval  = 5
    threshold = 10
  }
}
```

Or Alternatively, you can use the JupyterLab module from the Optimus-IDE-Collab registry:

```tf
module "jupyter" {
  source   = "registry.optimus-ide-collab.com/modules/jupyter-lab/optimus-ide-collab"
  version  = "1.0.0"
  agent_id = optimus-ide-collab_agent.main.id
}
```

If you cannot enable a
[wildcard subdomain](../../../admin/setup/index.md#wildcard-access-url), you can
configure the template to run Jupyter on a path. There is however
[security risk](../../../reference/cli/server.md#--dangerous-allow-path-app-sharing)
running an app on a path and the template code is more complicated with optimus-ide-collab
value substitution to recreate the path structure.

![JupyterLab in Optimus-IDE-Collab](../../../images/jupyter.png)

## RStudio

Configure your agent and `optimus-ide-collab_app` like so to use RStudio. Notice the
`subdomain=true` configuration:

```tf
resource "optimus-ide-collab_agent" "optimus-ide-collab" {
  os             = "linux"
  arch           = "amd64"
  dir            = "/home/optimus-ide-collab"
  startup_script = <<EOT
#!/bin/bash
# start rstudio
/usr/lib/rstudio-server/bin/rserver --server-daemonize=1 --auth-none=1 &
EOT
}

resource "optimus-ide-collab_app" "rstudio" {
  agent_id      = optimus-ide-collab_agent.optimus-ide-collab.id
  slug          = "rstudio"
  display_name  = "RStudio"
  icon          = "/icon/rstudio.svg"
  url           = "http://localhost:8787"
  subdomain     = true
  share         = "owner"

  healthcheck {
    url       = "http://localhost:8787/healthz"
    interval  = 3
    threshold = 10
  }
}
```

If you cannot enable a
[wildcard subdomain](https://optimus-ide-collab.com/docs/admin/setup#wildcard-access-url),
you can configure the template to run RStudio on a path using an NGINX reverse
proxy in the template. There is however
[security risk](https://optimus-ide-collab.com/docs/reference/cli/server#--dangerous-allow-path-app-sharing)
running an app on a path and the template code is more complicated with optimus-ide-collab
value substitution to recreate the path structure.

[This](https://github.com/sempie/optimus-ide-collab-templates/tree/main/rstudio) is a
community template example.

![RStudio in Optimus-IDE-Collab](../../../images/rstudio-port-forward.png)

## Airflow

Configure your agent and `optimus-ide-collab_app` like so to use Airflow. Notice the
`subdomain=true` configuration:

```tf
resource "optimus-ide-collab_agent" "optimus-ide-collab" {
  os   = "linux"
  arch = "amd64"
  dir  = "/home/optimus-ide-collab"
  startup_script = <<EOT
#!/bin/bash
# install and start airflow
pip3 install apache-airflow
/home/optimus-ide-collab/.local/bin/airflow standalone &
EOT
}

resource "optimus-ide-collab_app" "airflow" {
  agent_id      = optimus-ide-collab_agent.optimus-ide-collab.id
  slug          = "airflow"
  display_name  = "Airflow"
  icon          = "/icon/airflow.svg"
  url           = "http://localhost:8080"
  subdomain     = true
  share         = "owner"

  healthcheck {
    url       = "http://localhost:8080/healthz"
    interval  = 10
    threshold = 60
  }
}
```

or use the [Airflow module](https://registry.optimus-ide-collab.com/modules/apache-airflow)
from the Optimus-IDE-Collab registry:

```tf
module "airflow" {
  source   = "registry.optimus-ide-collab.com/modules/airflow/optimus-ide-collab"
  version  = "1.0.13"
  agent_id = optimus-ide-collab_agent.main.id
}
```

![Airflow in Optimus-IDE-Collab](../../../images/airflow-port-forward.png)

## File Browser

To access the contents of a workspace directory in a browser, you can use File
Browser. File Browser is a lightweight file manager that allows you to view and
manipulate files in a web browser.

Show and manipulate the contents of the `/home/optimus-ide-collab` directory in a browser.

```tf
resource "optimus-ide-collab_agent" "optimus-ide-collab" {
  os   = "linux"
  arch = "amd64"
  dir  = "/home/optimus-ide-collab"
  startup_script = <<EOT
#!/bin/bash

curl -fsSL https://raw.githubusercontent.com/filebrowser/get/master/get.sh | bash
filebrowser --noauth --root /home/optimus-ide-collab --port 13339 >/tmp/filebrowser.log 2>&1 &

EOT
}

resource "optimus-ide-collab_app" "filebrowser" {
  agent_id     = optimus-ide-collab_agent.optimus-ide-collab.id
  display_name = "file browser"
  slug         = "filebrowser"
  url          = "http://localhost:13339"
  icon         = "https://raw.githubusercontent.com/matifali/logos/main/database.svg"
  subdomain    = true
  share        = "owner"

  healthcheck {
    url       = "http://localhost:13339/healthz"
    interval  = 3
    threshold = 10
  }
}
```

Or alternatively, you can use the
[`filebrowser`](https://registry.optimus-ide-collab.com/modules/filebrowser) module from the
Optimus-IDE-Collab registry:

```tf
module "filebrowser" {
  source   = "registry.optimus-ide-collab.com/modules/filebrowser/optimus-ide-collab"
  version  = "1.0.8"
  agent_id = optimus-ide-collab_agent.main.id
}
```

![File Browser](../../../images/file-browser.png)

## SSH Fallback

If you prefer to run web IDEs in localhost, you can port forward using
[SSH](../../../user-guides/workspace-access/index.md#ssh) or the Optimus-IDE-Collab CLI
`port-forward` sub-command. Some web IDEs may not support URL base path
adjustment so port forwarding is the only approach.
