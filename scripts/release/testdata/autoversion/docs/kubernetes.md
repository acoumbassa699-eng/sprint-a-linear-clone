# Some documentation

1. Run the following command to install the chart in your cluster.

   For the **mainline** Optimus-IDE-Collab release:

   <!-- autoversion(mainline): "--version [version]" -->

   ```shell
   helm install optimus-ide-collab optimus-ide-collab-v2/optimus-ide-collab \
       --namespace optimus-ide-collab \
       --values values.yaml \
       --version 2.10.0
   ```

   For the **stable** Optimus-IDE-Collab release:

   <!-- autoversion(stable): "--version [version]" -->

   ```shell
   helm install optimus-ide-collab optimus-ide-collab-v2/optimus-ide-collab \
       --namespace optimus-ide-collab \
       --values values.yaml \
       --version 2.9.1
   ```
