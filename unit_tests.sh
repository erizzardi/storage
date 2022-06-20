#! /bin/bash

# N.B. at the moment (Jun 2022) this works EXCLUSIVELY with Postgres.

cleanup() {
    CONTAINER_NAME="test-$TEST_DB"
    # Checks if a container named $CONTAINER_NAME is running
    RUNNIG_CONTAINER_ID=$(docker ps -q -f name="$CONTAINER_NAME")
    if [[ -n "$RUNNIG_CONTAINER_ID" ]]; then
        echo "Stopping running $TEST_DB container: $RUNNIG_CONTAINER_ID"
        docker stop $RUNNIG_CONTAINER_ID > /dev/null
    fi

    # Checks if a container named $CONTAINER_NAME is stored in the local registry
    STORED_CONTAINER_ID=$(docker container ls -a -q -f name="$CONTAINER_NAME")
    if [[ -n "$STORED_CONTAINER_ID" ]]; then
        echo "Deleting stored $TEST_DB container: $STORED_CONTAINER_ID"
        docker rm "$STORED_CONTAINER_ID" > /dev/null
    fi
}
trap cleanup EXIT


DB_SUPPORTED=("postgres")

# Default values for the database name and version to be used for the tests.
# The service supports (maybe in the future!) different databases, so override this value
# to use another kind of db (cockroach, mysql, cassandra etc)
if [[ -z $TEST_DB ]]; then
    TEST_DB="postgres"
fi
if [[ -z $TEST_DB_VERSION ]]; then
    TEST_DB_VERSION="alpine3.16"
fi
# Defaults the test db to localhost:5432.
# This values can be overridden if using, for example, 
# a dummy db in a gitlab runner instance. 
if [[ -z $TEST_DB_HOST ]]; then
    TEST_DB_HOST="127.0.0.1"
fi
if [[ -z $TEST_DB_PORT ]]; then
    TEST_DB_PORT="5432"
fi
# Default values for user authentication.
if [[ -z $TEST_DB_USER ]]; then
    TEST_DB_USER="postgres"
fi
if [[ -z $TEST_DB_PASSWORD ]]; then
    TEST_DB_PASSWORD="password"
fi
if [[ -z $TEST_DB_DATABASE ]]; then
    TEST_DB_DATABASE="test-db"
fi

CONTAINER_NAME="test-$TEST_DB"


cleanup

    
# Pulls the image for the test database, if allowed.
# The docker client can already tell whether
# the image is already present in the local registry (exit code still 0)
if [[ " ${DB_SUPPORTED[*]} " =~ " ${TEST_DB} " ]]; then
    docker pull "$TEST_DB:$TEST_DB_VERSION"
fi

if [[ ! " ${DB_SUPPORTED[*]} " =~ " ${TEST_DB} " ]]; then
    echo "Database not supported! Allowed values: ${DB_SUPPORTED[*]}"
    exit 1
fi




case "$TEST_DB" in

  "postgres")

    # Connection string for the test database is created here, then injected in the testing execution.
    TEST_DB_CONN_STR="postgresql://${TEST_DB_USER}:${TEST_DB_PASSWORD}@${TEST_DB_HOST}:${TEST_DB_PORT}/${TEST_DB_DATABASE}?sslmode=disable"

    docker run -d \
	    --name "$CONTAINER_NAME" \
        -e POSTGRES_USER="$TEST_DB_USER" \
	    -e POSTGRES_PASSWORD="$TEST_DB_PASSWORD" \
        -p "$TEST_DB_PORT:$TEST_DB_PORT" \
	    "$TEST_DB:$TEST_DB_VERSION"
    ;;
esac

# TODO Add unit tests execution here
# TODO Add unit tests
# TEST_DB_CONN_STR="$TEST_DB_CONN_STR" go test -v
echo "Running tests..."

# The end! The trap should now activate