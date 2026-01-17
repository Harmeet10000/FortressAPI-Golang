```go
17. check and correct the use of ENV variables with koanf sanke_case CAPITAL         DONE
3. undertand how env are loaded in structs                                           DONE
20. check air implementation                                                         DONE   
14. figure out what comes in database and db folder                                  DONE
23. rewrite all middlewares for echo
5. understand how logging is done and how should it be used
12. understand pkg folder role and use and figure out lib/helper required or not
9. check how it put useful(corelationalID, more) in request context like logger, db, cache, asynq, helpers/utils
21. check all echo features provided grouping "api/v1" and then "auth" and rest
11. rewrite httpResponse and understand echo.NewHTTPError
1. understand sqlc queries, schema, out
16. learn how to serve over https/http2 over trafik
4. make seed data and schema for DB
6. understand how middlewares are made used 
7. understand how error handling is done
8. rewrite health folder if necessary
10. swagger docs
24. use redis for rate limit
2. understand golangci.yml
18. do a go mod tidy carefully
15. make docker-compose working again
19. learn how syntax highlighting works
22. check all the available option in redis newClient and if it is correctly implemented 
13. last understand FX usage and usfulness
25. explore swaggo/swag and echo-swagger for OpenAPI docs


// Suppress a specific linter for one line
var foo string //nolint:errcheck
result := someFunc() // nolint

// Suppress multiple linters
result := someFunc() // nolint:errcheck,gosec

// Suppress for entire function
//nolint:funlen
func VeryLongFunction() {
    // ...
}

// Suppress for entire file
//nolint
package mypackage
```

