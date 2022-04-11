# storage

Attempt at a kubernetes native object storage, with custom JSON API.

## postgres
1. the user and the databases need to be created in order for the service to work

## TODO
1. improve logging with default fiels - to differentiate between layers
2. UNIT TESTS
3. improve response writing - headers are fucked up
4. use buffered read/writes, to cap memory usage
5. improve container compatibility
6. object lifecycle policies
7. object versioning
8. implement buckets
9. check kubernetes compatibility
10. authentication and permissions
11. implement DB interface for other kinds of DBs (ideally: MySQL, CockroachDB, Couchbase, Cassandra)