# storage

Attempt at a kubernetes native object storage, with custom JSON API.

## postgres
1. the user and the databases need to be created in order for the service to work

## Unit tests
Unit tests require a db instance running on `$TEST_DB_ADDR:$TEST_DB_PORT`. The script `unit_tests.sh` spins up a db isntance automatically, and launches the unit tests. This is the preferred way to execute unit tests, and it should be used in a CI/CD environment. 

## TODO
1. UNIT TESTS
2. IMPROVE (write :p) DOCUMENTATION
3. <del>paged API to list objects,</del>
4. <del>API to set loglevel at runtime</del>
5. <del>improve response writing - headers are fucked up</del>
6. <del>remove default values and have them read from secrets as env variables</del>
7. error management - have service methods return custom error type (ResponseError), so to avoid type assertions
8. improve read/write of large files - buffered IO operations to cap memory? write to binary?
9. implement caching - check varnish compatibility
10. improve container compatibility
11. object lifecycle policies
12. object versioning
13. implement buckets
14. check kubernetes compatibility
15. helm chart
16. authentication and permissions
17. ad methods that prepare every possible query
18. implement DB interface for other kinds of DBs (ideally: MySQL, CockroachDB, Couchbase, Cassandra)