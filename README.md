# Status Checker

Status Checker is a Go application that schedules periodic status checks. If a check fails, it executes a predefined recovery command. The application includes a built-in UI for monitoring, supports Prometheus metrics, and integrates with Slack for notifications.

## Features

- **Scheduled Status Checks**: Monitor the health of your services at regular intervals using cron specs.
- **Automatic Recovery**: Execute recovery commands automatically when a check fails.
- **Built-in UI**: Visualize the status of your services through a user-friendly interface.
- **Prometheus Metrics**: Integrate with Prometheus for monitoring and alerting.
- **Slack Notifications**: Receive alerts and notifications directly in your Slack channels.

## Installation

To install Status Checker, run the following command:

```bash
curl -s https://raw.githubusercontent.com/sbnarra/status-checker/refs/heads/main/install.sh | bash
```

## Configuration

### User Interface

After installation, access the UI at [http://localhost:8000](http://localhost:8000).

### Check Configuration

Define your checks in the `checks.yaml` file located at `/opt/status-checker/checks.yaml`.

Here's an example of a minutely check configuration:
```yaml
ssh:
  schedule: "* * * * *"
  command: sudo systemctl status ssh
  recover: sudo systemctl restart ssh
```

The `schedule` field also supports seconds, e.g. `* * * * * *`.
* Pattern: `OpSec Min Hour DayOfMonth Month DayOfWeek`

### Environment Configuration

Configure the application settings in the `config.env` file located at `/opt/status-checker/config.env`. Below are the available environment variables:

| Variable                   | Description                                                                  | Default Value                     |
|----------------------------|------------------------------------------------------------------------------|-----------------------------------|
| `CHECKS_PATH`              | Path(s) to the checks configuration file. Supports multiple files separated by commas. | `/opt/status-checker/checks.yaml` |
| `BIND_ADDR`                | Address to bind the UI server                                                 | `:8000`                           |
| `SERVER_ENABLED`           | Enable or disable the server                                                  | `true`                            |
| `DEBUG`                    | Enable or disable debug mode                                                  | `false`                           |
| `HISTORY_DIR`              | Directory to store history logs                                               | `/opt/status-checker/history`     |
| `MIN_HISTORY`              | Minimum number of history entries to keep                                     | `100`                             |
| `HISTORY_CHECK_SIZE_LIMIT` | Maximum size for history check files                                          | `10MB`                            |
| `PROMETHEUS_ENABLED`       | Enable or disable Prometheus metrics                                          | `true`                            |
| `SLACK_HOOK_URL`           | URL for Slack webhook to send notifications                                   | _n/a_                             |

## Usage

### Running Locally for Development

To run Status Checker locally for development purposes, use the following command:

```bash
go run main.go
```

You can optionally pass check paths as arguments:

```bash
go run main.go config/checks.yaml /path/to/additional-checks.yaml
```

## Monitoring

Status Checker exposes Prometheus metrics for monitoring the application's performance and health.

### Configure Prometheus Scraping

Add the following job to your Prometheus configuration (`prometheus.yml`):

```yaml
scrape_configs:
  - job_name: 'status_checker'
    static_configs:
      - targets: ['<EXTERNAL_IP>:8000']
```

Replace `<EXTERNAL_IP>` with the IP address where Status Checker is running. Ensure that Prometheus can access this IP and port.

#### Using Prometheus in Kubernetes

If Prometheus is running in Kubernetes, you can create a `ServiceMonitor` to scrape the external Status Checker:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: status-checker
  labels:
    release: prometheus
spec:
  endpoints:
    - port: metrics
      path: /metrics
      interval: 30s
      static_config:
        targets:
          - <EXTERNAL_IP>:8000
```

Apply the `ServiceMonitor`:

```bash
kubectl apply -f service-monitor.yaml
```

Replace `<EXTERNAL_IP>` with the actual IP address of your Status Checker. Ensure network accessibility between Kubernetes and the external Status Checker instance.

## Notifications

Configure the `SLACK_HOOK_URL` in the `config.env` file to enable Slack notifications. This allows Status Checker to send alerts directly to your Slack channels when a check fails or a recovery action is taken.

### Setting Up Slack Incoming Webhooks

1. **Create a Slack App**
    - Navigate to the [Slack API: Applications](https://api.slack.com/apps) page.
    - Click on **"Create New App"**.
    - Select **"From scratch"**, enter an app name (e.g., "Status Checker"), and choose your workspace.
    - Click **"Create App"**.

2. **Enable Incoming Webhooks**
    - In your app's settings, go to **"Incoming Webhooks"** in the sidebar.
    - Toggle **"Activate Incoming Webhooks"** to **"On"**.

3. **Create a Webhook URL**
    - Click on **"Add New Webhook to Workspace"**.
    - Select the channel where you want to receive notifications and click **"Allow"**.
    - Copy the generated **Webhook URL**.

4. **Configure `config.env`**
    - Open the `config.env` file located at `/opt/status-checker/config.env`.
    - Set the `SLACK_HOOK_URL` variable to the copied Webhook URL:

    ```env
    SLACK_HOOK_URL=https://hooks.slack.com/services/your/webhook/url
    ```

Now, Status Checker will send notifications to your specified Slack channel whenever a check fails or a recovery action is executed.