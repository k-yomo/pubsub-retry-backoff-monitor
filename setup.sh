#!/bin/bash

set -eu

PROJECT_ID=$1

function main() {
  brew list pubsub_cli &>/dev/null || (brew tap k-yomo/pubsub_cli && brew install pubsub_cli)
  gcloud auth application-default login
  bq mk "$PROJECT_ID":playground.pubsub_request pubsub_request.json
  gcloud functions deploy RecordPubsubHandler --runtime go113 --trigger-http --region=asia-northeast1 --ingress-settings=all --project "$PROJECT_ID" --set-env-vars PROJECT_ID="$PROJECT_ID"
  pubsub_cli register_push test-message https://asia-northeast1-"$PROJECT_ID".cloudfunctions.net/RecordPubsubHandler -p "$PROJECT_ID"
  echo "Setup finished! start monitoring retry backoff by 'make start'"
}

main