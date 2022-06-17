# storage

Attempt at a kubernetes native object storage, with custom JSON API.

## postgres
1. the user and the databases need to be created in order for the service to work

## Unit tests
Unit tests require a db instance running on `$TEST_DB_ADDR:$TEST_DB_PORT`. The script `unit_tests.sh` spins up a db isntance automatically, and launches the unit tests. This is the preferred way to execute unit tests, and it should be used in a CI/CD environment. 

## TODO
1. UNIT TESTS 
2. <del>paged API to list objects,</del>
3. <del>API to set loglevel at runtime</del>
4. <del>improve response writing - headers are fucked up</del>
5. <del>remove default values and have them read from secrets as env variables</del>
6. improve read/write of large files - buffered IO operations to cap memory?
7. implement caching - check varnish compatibility
8. improve container compatibility
9. object lifecycle policies
10. object versioning
11. implement buckets
12. check kubernetes compatibility
13. helm chart
14. authentication and permissions
15. ad methods that prepare every possible query
16. implement DB interface for other kinds of DBs (ideally: MySQL, CockroachDB, Couchbase, Cassandra)