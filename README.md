# storage

Attempt at a kubernetes native object storage, with custom JSON API.

## postgres
1. the user and the databases need to be created in order for the service to work

## TODO
1. improve logging with default fiels - to differentiate between layers
2. improve response writing - headers are fucked up
3. use buffered read/writes, to cap memory usage
4. improve container compatibility
5. object lifecycle policies
6. object versioning
7. implement buckets
8. check kubernetes compatibility
9. authentication and permissions
10. implement DB interface for other kinds of DBs (ideally: MySQL, CockroachDB, Couchbase, Cassandra)