# JFrog kubexray scanner on Kubernetes Helm Chart

## Prerequisites Details

* Kubernetes 1.10+

## Chart Details

This chart will do the following:

* Deploy JFrog KubeXray

## Requirements

- A running Kubernetes cluster
- Running Artifactory and Xray
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed and setup to use the cluster
- [Helm](https://helm.sh/) installed and setup to use the cluster (helm init)
- Configuration file `xray_config.yaml` with Xray server connection settings:

```
url: https://xray.mydomain.com
user: admin
password: password
slackWebhookUrl: ""
xrayWebhookToken: ""
```

**Note:** Configuration file `xray_config.yaml` must be provided.

## Configuration

### Slack notifications

Notification by Slack can be enabled by providing `slackWebhookUrl`:

```
url: https://xray.mydomain.com
user: admin
password: password
slackWebhookUrl: https://hooks.slack.com/services/your_slack_webhook_url 
xrayWebhookToken: ""
```

### Enable kubexray WebHook

If you want kubexray to react on Xray policy changes generate `xrayWebhooToken` with `openssl rand -base64 64 | tr -dc A-Za-z0-9`:

```
url: https://xray.mydomain.com
user: admin
password: password
slackWebhookUrl: https://hooks.slack.com/services/your_slack_webhook_url 
xrayWebhookToken: replace_with_generated_token
```

**Note:** Also you need to add `kubexray_url` your generated `xrayWebhooToken` to your Xray server under `Admin/Webhooks`.

## Install JFrog KubeXray

### Add JFrog Helm repository

Before installing JFrog helm charts, you need to add the [JFrog helm repository](https://charts.jfrog.io/) to your helm client

```bash
helm repo add jfrog https://charts.jfrog.io
```

### Install Chart

#### Install JFrog KubeXray

```bash
helm install --name kubexray --namespace kubexray jfrog/kubexray \
    --set xrayConfig="$(cat path_to_your/xray_config.yaml | base64)"
```

#### Installing with existing secret

You can deploy the KubeXray configuration file `xray_config.yaml` as a [Kubernetes secret](https://kubernetes.io/docs/concepts/configuration/secret/).


Create the Kubernetes secret

```bash
kubectl create secret generic kubexray --from-file=path_to_your/xray_config.yaml
```

Pass the configuration file to helm

```bash
 helm install --name kubexray --namespace kubexray jfrog/kubexray \
    --set existingSecret="kubexray"
```

**NOTE:** You have to keep passing the configuration file secret parameter as `--set existingSecret="kubexray"` on all future calls to `helm install` and `helm upgrade` or set it in `values.yaml` file `existingSecret: kubexray`!

## Status

See the status of your deployed **helm** release

```bash
helm status kubexray
```

## Upgrade

E.g you have changed scan policy rules and to need upgrade an existing kubexray release

```bash
helm upgrade kubexray --namespace kubexray jfrog/kubexray \
    --set xrayConfig="$(cat path_to_your/xray_config.yaml | base64)"
```

Upgrading with existing secret

```bash
helm upgrade --install kubexray --namespace kubexray jfrog/kubexray \
    --set existingSecret="kubexray"
```

## Remove

Removing a **helm** release is done with

```bash
# Remove the Xray services and data tools
helm delete --purge kubexray
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the xray chart and their default values.

|         Parameter            |                    Description                   |           Default                  |
|------------------------------|--------------------------------------------------|------------------------------------|
| `image.PullPolicy`| Container pull policy | `IfNotPresent` |
| `xrayConfig` | base64 encoded `xray_config.yaml` file |  |
| `existingSecret` | Specifies an existing secret holding the Xray config |  |
| `securityContext.enabled` | Enables Security Context  | `false` |
| `securityContext.enabled` |  Security UserId | `1000` |
| `securityContext.kubeXrayUserId` |  Security GroupId | `1000` |
| `scanPolicy.unscanned.deployments` | Specifies unscanned Deployments policy | `ignore` |
| `scanPolicy.unscanned.statefulSets` | Specifies unscanned StatefulSets policy | `ignore` |
| `scanPolicy.unscanned.whiltelistNamespaces` | Specifies unscanned whiltelist Namespaces list | `kube-system` |
| `scanPolicy.security.deployments` | Specifies Deployments with security issues policy | `ignore` |
| `scanPolicy.security.statefulSets` | Specifies Deployments with security issues policy  | `ignore` |
| `scanPolicy.license.deployments` | Specifies Deployments with license issues policy | `ignore` |
| `scanPolicy.license.statefulSets` | Specifies StatefulSets with license issues policy | `ignore` |
| `rbac.enabled` | Specifies whether RBAC resources should be created | `true` |
| `resources.limits.cpu` | Specifies CPU limit | `256m` |
| `resources.limits.memory` | Specifies memory limit | `128Mi` |
| `resources.requests.cpu` | Specifies CPU request | `100m` |
| `resources.requests.memory` | Specifies memory request | `128Mi` |
| `nodeSelector` | kubexray micro-service node selector | `{}` |
| `tolerations` | kubexray micro-service node tolerations | `[]` |
| `affinity` | kubexray micro-service node affinity | `{}` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install/upgrade`.

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example

```bash
helm upgrade kubexray --namespace kubexray jfrog/kubexray \
    --set --set existingSecret="kubexray",existingSecretKey="xray_config.yaml" -f override-values.yaml 
```
