export AUTH_JWT_SECRET="some_silly_secret"
export AUTH_API_RESTRICTED="false"

# Litestream not on in dev

export DB_FILE="./data/database.db" # NOTE: referenced from scripts, not app
export DB_URI="file:${DB_FILE}?_journal=WAL&_timeout=5000&_fk=true&_sync=NORMAL&_txlock=immediate"
export DB_DRIVER_NAME="sqlite3"

export HTTP_STATIC_DIR="_ui/dist"
export HTTP_IDLE_SHUTDOWN_MS="0"

export LOG_LEVEL="debug"
export LOG_OUTPUT_FORMAT="text"

export RPC_HOST=""
export RPC_PORT="5001"
export RPC_PATH_PREFIX="/rpc"
