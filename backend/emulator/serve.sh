#!/bin/bash

set -euo pipefail

export firebase_port_database=${FIREBASE_PORT_DATABASE:-9000}
export firebase_port_firestore=${FIREBASE_PORT_DATABASE:-8090}
export firebase_port_auth=${FIREBASE_PORT_AUTH:-9099}
export firebase_port_ui=${FIREBASE_PORT_UI:-4000}
export firebase_project_id=${FIREBASE_PROJECT_ID:-demo-project}

cd /app/


(
    firebase emulators:start 2>&1 \
        --project "$firebase_project_id" \
        --import="/app/data" \
        --only auth,firestore \
            | sed -ur 's/^/:: [firebase] /'
) &
firebase_pid=$!

echo ":: ready"

wait $firebase_pid
