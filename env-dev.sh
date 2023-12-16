export AUTH_JWT_SECRET="some_silly_secret"
export AUTH_API_RESTRICTED="false"

# Litestream not on in dev

export DB_FILE="./data/database.db" # just used in scripts, not config
export DB_URI="file:${DB_FILE}?mode=rw&cache=shared&_journal_mode=WAL&_fk=1"
export DB_DRIVER_NAME="sqlite3"

export HTTP_STATIC_DIR="_ui/dist"
export HTTP_IDLE_SHUTDOWN_MS="0"

export LOG_LEVEL="debug"
export LOG_OUTPUT_FORMAT="text"

export RPC_HOST=""
export RPC_PORT="5001"
export RPC_PATH_PREFIX="/rpc"
