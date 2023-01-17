# KUBE EVENT ALERT

Kube event alert is a k8s controller that watches for events in the cluster and alerts about any abnormal events by calling a webhook with the event details. The webhook payload support slack incoming webhook payload. The controller can run both in and out of the k8s cluster.

## Usage

```sh
go run main.go [OPTIONS]

OPTIONS:
-kubeconfig
    path to a kubeconfig if running out of cluster
-masterUrl
    url to kube cluster if running out of cluster
    can be specified with KUBE_MASTER_URL environment variable as well
-webhookURL
    notification will be sent to this incoming webhook url
    can be specified with WEBHOOK_URL environment variable as well
```

## Installation

### In Cluster

Follow the instructions in [k8s](k8s) to configure `webhook.url` and deploy the manifest:

```bash
kustomize -o .deploymnet.yml ./k8s
kubectl apply -f deployment.yml
```

### Out of Cluster

#### Executable

Compile an executable file

```sh
go build -o kube-event-alert .
```

Run the executable according to the [usage](#usage) instructions above and satisfy `webhookURL` argument and either the `kubeconfig` or `masterUrl` arguments

#### Docker

```sh
docker build -t <name>:<tag> .
docker run -d -e WEBHOOK_URL=<webhook-url> -e KUBE_MASTER_URL=<kube-master-url> <name>:<tag>
```

## Notification

The event notification is sent to the webhook url via HTTP `POST` method with the following payload:

```json
{
    "text": "event details"
}
```

The payload supports the slack incoming webhook integration.
