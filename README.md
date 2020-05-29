# pubsub-retry-backoff-monitor
This repo is for monitoring Cloud Pub/Sub's retry backoff by recording the requests.

## Usage
```shell script
$ ./setup.sh
$ pubsub_cli publish test-message 'test' -p {YOUR_GCP_PROJECT_ID}
```