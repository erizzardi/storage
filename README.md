# storage

Attempt at a kubernetes native object storage, with custom JSON API.

## APIs
TODO

## postgres
1. the user and the databases need to be created in order for the service to work

## Unit tests
Unit tests require a db instance running on `$TEST_DB_ADDR:$TEST_DB_PORT`. The script `unit_tests.sh` spins up a db isntance automatically, and launches the unit tests. This is the preferred way to execute unit tests, and it should be used in a CI/CD environment. 

## Caveats on error handling
Errors should NEVER be returned with go-kit regular funcions. They should be embedded in the request payload and handled in the correct layer. The `service` layer should handle all the 500-errors, and the `endpoint` layer should handle all the 400-errors (probably).


## TODO
1. UNIT TESTS
2. IMPROVE (write :p) DOCUMENTATION
3. <del>paged API to list objects,</del>
4. <del>API to set loglevel at runtime</del>
5. <del>Endpoints layer logging</del>
6. <del>improve response writing - headers are fucked up</del>
7. <del>remove default values and have them read from secrets as env variables</del>
8. APIs to manipulate files by name
9. DB middleware to implement retry and timeouts
10. object versioning
11. <del>error management - have service methods return custom error type (ResponseError), so to avoid type assertions</del>
12. improve read/write of large files - buffered IO operations to cap memory? write to binary?
13. implement caching - check varnish compatibility
14. improve container compatibility
15. object lifecycle policies
16. implement buckets
17. check kubernetes compatibility
18. helm chart
19. authentication and permissions
20. add methods that prepare every possible query
21. implement DB interface for other kinds of relational DBs (ideally: MySQL, CockroachDB)