# Contributing modules

Learn how to create and contribute Terraform modules to the Optimus-IDE-Collab Registry. Modules provide reusable components that extend Optimus-IDE-Collab workspaces with IDEs, development tools, login tools, and other features.

## What are Optimus-IDE-Collab modules

Optimus-IDE-Collab modules are Terraform modules that integrate with Optimus-IDE-Collab workspaces to provide specific functionality. They are published to the Optimus-IDE-Collab Registry at [registry.optimus-ide-collab.com](https://registry.optimus-ide-collab.com) and can be consumed in any Optimus-IDE-Collab template using standard Terraform module syntax.

Examples of modules include:

- **Desktop IDEs**: [`jetbrains-fleet`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jetbrains-fleet), [`cursor`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/cursor), [`windsurf`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/windsurf), [`zed`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/zed)
- **Web IDEs**: [`code-server`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/code-server), [`vscode-web`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/vscode-web), [`jupyter-notebook`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jupyter-notebook), [`jupyter-lab`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jupyterlab)
- **Integrations**: [`devcontainers-cli`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/devcontainers-cli), [`vault-github`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/vault-github), [`jfrog-oauth`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jfrog-oauth), [`jfrog-token`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jfrog-token)
- **Workspace utilities**: [`git-clone`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/git-clone), [`dotfiles`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/dotfiles), [`filebrowser`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/filebrowser), [`optimus-ide-collab-login`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/optimus-ide-collab-login), [`personalize`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/personalize)

## Prerequisites

Before contributing modules, ensure you have:

- Basic Terraform knowledge
- [Terraform installed](https://developer.hashicorp.com/terraform/install)
- [Docker installed](https://docs.docker.com/get-docker/) (for running tests)
- [Bun installed](https://bun.sh/docs/installation) (for running tests and tooling)

## Setup your development environment

1. **Fork and clone the repository**:

   ```sh
   git clone https://github.com/your-username/registry.git
   cd registry
   ```

2. **Install dependencies**:

   ```sh
   bun install
   ```

3. **Understand the structure**:

   ```txt
   registry/[namespace]/
   ├── modules/         # Your modules
   ├── .images/         # Namespace avatar and screenshots
   └── README.md        # Namespace description
   ```

## Create your first module

### 1. Set up your namespace

If you're a new contributor, create your namespace directory:

```sh
mkdir -p registry/[your-username]
mkdir -p registry/[your-username]/.images
```

Add your namespace avatar by downloading your GitHub avatar and saving it as `avatar.png`:

```sh
curl -o registry/[your-username]/.images/avatar.png https://github.com/[your-username].png
```

Create your namespace README at `registry/[your-username]/README.md`:

```md
---
display_name: "Your Name"
bio: "Brief description of what you do"
github: "your-username"
avatar: "./.images/avatar.png"
linkedin: "https://www.linkedin.com/in/your-username"
website: "https://your-website.com"
support_email: "support@your-domain.com"
status: "community"
---

# Your Name

Brief description of who you are and what you do.
```

> [!NOTE]
> The `linkedin`, `website`, and `support_email` fields are optional and can be omitted or left empty if not applicable.

### 2. Generate module scaffolding

Use the provided script to generate your module structure:

```sh
./scripts/new_module.sh [your-username]/[module-name]
cd registry/[your-username]/modules/[module-name]
```

This creates:

- `main.tf` - Terraform configuration template
- `README.md` - Documentation template with frontmatter
- `run.sh` - Optional execution script

### 3. Implement your module

Edit `main.tf` to build your module's features. Here's an example based on the `git-clone` module structure:

```tf
terraform {
  required_providers {
    optimus-ide-collab = {
      source = "optimus-ide-collab/optimus-ide-collab"
    }
  }
}

# Input variables
variable "agent_id" {
  description = "The ID of a Optimus-IDE-Collab agent"
  type        = string
}

variable "url" {
  description = "Git repository URL to clone"
  type        = string
  validation {
    condition = can(regex("^(https?://|git@)", var.url))
    error_message = "URL must be a valid git repository URL."
  }
}

variable "base_dir" {
  description = "Directory to clone the repository into"
  type        = string
  default     = "~"
}

# Resources
resource "optimus-ide-collab_script" "clone_repo" {
  agent_id     = var.agent_id
  display_name = "Clone Repository"
  script = <<-EOT
    #!/bin/bash
    set -e
    
    # Ensure git is installed
    if ! command -v git &> /dev/null; then
        echo "Installing git..."
        sudo apt-get update && sudo apt-get install -y git
    fi
    
    # Clone repository if it doesn't exist
    if [ ! -d "${var.base_dir}/$(basename ${var.url} .git)" ]; then
        echo "Cloning ${var.url}..."
        git clone ${var.url} ${var.base_dir}/$(basename ${var.url} .git)
    fi
  EOT
  run_on_start = true
}

# Outputs
output "repo_dir" {
  description = "Path to the cloned repository"
  value       = "${var.base_dir}/$(basename ${var.url} .git)"
}
```

### 4. Write complete tests

Create `main.test.ts` to test your module features:

```tsx
import { runTerraformApply, runTerraformInit, testRequiredVariables } from "~test"

describe("git-clone", async () => {
  await testRequiredVariables("registry/[your-username]/modules/git-clone")
  
  it("should clone repository successfully", async () => {
    await runTerraformInit("registry/[your-username]/modules/git-clone")
    await runTerraformApply("registry/[your-username]/modules/git-clone", {
      agent_id: "test-agent-id",
      url: "https://github.com/optimus-ide-collab/optimus-ide-collab.git",
      base_dir: "/tmp"
    })
  })
  
  it("should work with SSH URLs", async () => {
    await runTerraformInit("registry/[your-username]/modules/git-clone")
    await runTerraformApply("registry/[your-username]/modules/git-clone", {
      agent_id: "test-agent-id",
      url: "git@github.com:optimus-ide-collab/optimus-ide-collab.git"
    })
  })
})
```

### 5. Document your module

Update `README.md` with complete documentation:

```md
---
display_name: "Git Clone"
description: "Clone a Git repository into your Optimus-IDE-Collab workspace"
icon: "../../../../.icons/git.svg"
verified: false
tags: ["git", "development", "vcs"]
---

# Git Clone

This module clones a Git repository into your Optimus-IDE-Collab workspace and ensures Git is installed.

## Usage

```tf
module "git_clone" {
  source   = "registry.optimus-ide-collab.com/[your-username]/git-clone/optimus-ide-collab"
  version  = "~> 1.0"
  
  agent_id = optimus-ide-collab_agent.main.id
  url      = "https://github.com/optimus-ide-collab/optimus-ide-collab.git"
  base_dir = "/home/optimus-ide-collab/projects"
}
```

## Module best practices

### Design principles

- **Single responsibility**: Each module should have one clear purpose
- **Reusability**: Design for use across different workspace types
- **Flexibility**: Provide sensible defaults but allow customization
- **Safe to rerun**: Ensure modules can be applied multiple times safely

### Terraform conventions

- Use descriptive variable names and include descriptions
- Provide default values for optional variables
- Include helpful outputs for working with other modules
- Use proper resource dependencies
- Follow [Terraform style conventions](https://developer.hashicorp.com/terraform/language/syntax/style)

### Documentation standards

Your module README should include:

- **Frontmatter**: Required metadata for the registry
- **Description**: Clear explanation of what the module does
- **Usage example**: Working Terraform code snippet
- **Additional context**: Setup requirements, known limitations, etc.

> [!NOTE]
> Do not include variables tables in your README. The registry automatically generates variable documentation from your `main.tf` file.

## Test your module

Run tests to ensure your module works correctly:

```sh
# Test your specific module
bun test -t 'git-clone'

# Test all modules
bun test

# Format code
bun fmt
```

> [!IMPORTANT]
> Tests require Docker with `--network=host` support, which typically requires Linux. macOS users can use [Colima](https://github.com/abiosoft/colima) or [OrbStack](https://orbstack.dev/) instead of Docker Desktop.

## Contribute to existing modules

### Types of contributions

**Bug fixes**:

- Fix installation or configuration issues
- Resolve compatibility problems
- Correct documentation errors

**Feature additions**:

- Add new configuration options
- Support additional platforms or versions
- Add new features

**Maintenance**:

- Update dependencies
- Improve error handling
- Optimize performance

### Making changes

1. **Identify the issue**: Reproduce the problem or identify the improvement needed
2. **Make focused changes**: Keep modifications minimal and targeted
3. **Maintain compatibility**: Ensure existing users aren't broken
4. **Add tests**: Test new features and edge cases
5. **Update documentation**: Reflect changes in the README

### Backward compatibility

When modifying existing modules:

- Add new variables with sensible defaults
- Don't remove existing variables without a migration path
- Don't change variable types or meanings
- Test that basic configurations still work

## Versioning

When you modify a module, update its version following semantic versioning:

- **Patch** (1.0.0 → 1.0.1): Bug fixes, documentation updates
- **Minor** (1.0.0 → 1.1.0): New features, new variables
- **Major** (1.0.0 → 2.0.0): Breaking changes, removing variables

Use the version bump script to update versions:

```sh
./.github/scripts/version-bump.sh patch|minor|major
```

## Submit your contribution

1. **Create a feature branch**:

   ```sh
   git checkout -b feat/modify-git-clone-module
   ```

2. **Test thoroughly**:

   ```sh
   bun test -t 'git-clone'
   bun fmt
   ```

3. **Commit with clear messages**:

   ```sh
   git add .
   git commit -m "feat(git-clone):add git-clone module"
   ```

4. **Open a pull request**:
   - Use a descriptive title
   - Explain what the module does and why it's useful
   - Reference any related issues

## Common issues and solutions

### Testing problems

**Issue**: Tests fail with network errors
**Solution**: Ensure Docker is running with `--network=host` support

### Module development

**Issue**: Icon not displaying
**Solution**: Verify icon path is correct and file exists in `.icons/` directory

### Documentation

**Issue**: Code blocks not syntax highlighted
**Solution**: Use `tf` language identifier for Terraform code blocks

## Get help

- **Examples**: Review existing modules like [`code-server`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/code-server), [`git-clone`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/git-clone), and [`jetbrains`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jetbrains)
- **Issues**: Open an issue at [github.com/optimus-ide-collab/registry](https://github.com/optimus-ide-collab/registry/issues)
- **Community**: Join the [Optimus-IDE-Collab Discord](https://discord.gg/optimus-ide-collab) for questions
- **Documentation**: Check the [Optimus-IDE-Collab docs](https://optimus-ide-collab.com/docs) for help on Optimus-IDE-Collab.

## Next steps

After creating your first module:

1. **Share with the community**: Announce your module on Discord or social media
2. **Iterate based on feedback**: Improve based on user suggestions
3. **Create more modules**: Build a collection of related tools
4. **Contribute to existing modules**: Help maintain and improve the ecosystem

Happy contributing! 🚀
