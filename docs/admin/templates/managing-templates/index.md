# Working with templates

You create and edit Optimus-IDE-Collab templates as
[Terraform](https://developer.hashicorp.com/terraform/intro) configuration files (`.tf`) and
any supporting files, like a README or configuration files for other services.

## Who creates templates?

The [Template Admin](../../../admin/users/groups-roles.md#roles) role (and
above) can create templates. End users, like developers, create workspaces from
them. Templates can also be [managed with git](./change-management.md), allowing
any developer to propose changes to a template.

You can give different users and groups access to templates with
[role-based access control](../template-permissions.md).

## Creating templates

The [template builder](../creating-templates.md#template-builder) is
the recommended way to create templates. It guides you through selecting a base
infrastructure template, adding modules (IDEs, tools, integrations), and
configuring template settings without writing Terraform.

Starter templates for common cloud providers (AWS, Azure) and orchestrators
(Kubernetes, Docker) are available as base templates within the builder. You can
modify the generated template to use your own images, VPC, cloud credentials,
and so on. Optimus-IDE-Collab supports all Terraform resources and properties.

If you prefer to use Optimus-IDE-Collab on the
[command line](../../../reference/cli/index.md), use `optimus-ide-collab templates init` to
pull a starter template, then `optimus-ide-collab templates push` to upload it.

Optimus-IDE-Collab starter templates are also available on our
[GitHub repo](../../../../examples/templates).

## Community Templates

As well as Optimus-IDE-Collab's starter templates, you can see a list of community templates
by our users
[here](../../../../examples/templates/community-templates.md).

## Editing templates

Our templates are meant to be modified for your use cases. You can edit
any template's files directly in the Optimus-IDE-Collab dashboard.

![Editing a template](../../../images/templates/choosing-edit-template.gif)

If you'd prefer to use the CLI, use `optimus-ide-collab templates pull`, edit the template
files, then `optimus-ide-collab templates push`.

> [!TIP]
> Even if you are a Terraform expert, we suggest reading our
> [guided tour of a template](../../../tutorials/template-from-scratch.md).

## Updating templates

Optimus-IDE-Collab tracks a template's versions, keeping all developer workspaces up-to-date.
When you publish a new version, developers are notified to get the latest
infrastructure, software, or security patches. Learn more about
[change management](./change-management.md).

![Updating a template](../../../images/templates/update.png)

### Template update policies

> [!NOTE]
> Template update policies are a Premium feature.
> [Learn more](https://optimus-ide-collab.com/pricing#compare-plans).

Licensed template admins may want workspaces to always remain on the latest
version of their parent template. To do so, enable **Template Update Policies**
in the template's general settings. All non-admin users of the template will be
forced to update their workspaces before starting them once the setting is
applied. Workspaces which leverage autostart or start-on-connect will be
automatically updated on the next startup.

![Template update policies](../../../images/templates/update-policies.png)

## Delete templates

You can delete a template using both the optimus-ide-collab CLI and UI. Only
[template admins and owners](../../users/groups-roles.md#roles) can delete a
template, and the template must not have any running workspaces associated to
it.

In the UI, navigate to the template you want to delete, and select the dropdown
in the right-hand corner of the page to delete the template.

![delete-template](../../../images/delete-template.png)

Using the CLI, login to Optimus-IDE-Collab and run the following command to delete a
template:

```sh
optimus-ide-collab templates delete <template-name>
```

## Next steps

- [Image management](./image-management.md)
- [Dev Containers integration](../../integrations/devcontainers/integration.md) (recommended)
- [Envbuilder](../../integrations/devcontainers/envbuilder/index.md) (alternative for environments without Docker)
- [Change management](./change-management.md)
