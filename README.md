# storage

Attempt at a kubernetes native object storage, with custom JSON API.

## postgres
1. the user and the databases need to be created in order for the service to work

## TODO
1. UNIT TESTS
2. paged API to list objects
3. <del>API to set loglevel at runtime</del>
4. improve response writing - headers are fucked up
5. use buffered read/writes, to cap memory usage
6. improve container compatibility
7. object lifecycle policies
8. object versioning
9. implement buckets
10. check kubernetes compatibility
11. authentication and permissions
12. implement DB interface for other kinds of DBs (ideally: MySQL, CockroachDB, Couchbase, Cassandra)