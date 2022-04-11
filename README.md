# storage

Attempt at a kubernetes native object storage, with custom JSON API.

## postgres
1. the user and the databases need to be created in order for the service to work

## TODO
1. improve logging with default fiels - to differentiate between layers
2. object lifecycle policies
3. object versioning
4. buckets
5. authentication and permissions
6. implement DB interface for other kinds of DBs (ideally: MySQL, CockroachDB, Couchbase, Cassandra)