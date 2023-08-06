### TODO (ideas n' such):


    - [ ] Add context/request id to toerr message reporting
    - [ ] Graceful shutdown via ctx
    - [ ] Make session remarshalling happen at the end of the request, not the beginning more seemlessly not via custom ServeHTTP Function
    - [ ] THIS WORKS -> what if : return context.WithValue(ctx, &m{}, m) - Change Contextable to use a interface instead of a string?
    - [ ] Change 'Repo' to Remotes db+auth svc
    - [ ] Completely confused if id's should or should not be included in user error Msg
    - [ ] Add context, and/or RouteInfo BodyInfo SessionInfo to toerrs for better error debugging
    - [ ] Should SessionInfo contain the full User?
    - [ ] Add panics to Bi and Ri check?
    - [ ] Consider adding more calendar providers than just google
    - [ ] Convert time ranges to use range tree
    - [ ] Project structure - https://www.gobeyond.dev/standard-package-layout/
    - [ ] Models need proper error returns for Not found, invalid, and server error, or already exists etc
    - [ ] https://golangci-lint.run/usage/linters/
    - [ ] Rename packages to not use snake case
    - [ ] Make transactional db calls with cancellable context
    - [ ] Divide routes by content type
    - [ ] Rename models to services
    - [ ] More integrated db migrationator
    - [ ] Add pagination to resources
    - [ ] Set Cookie expiration?
    - [ ] sql.Null types?
    - [ ] fix frame and model package ugly interdependency
    - [ ] Use context to limit total request time
    - [ ] Version api using content type:+[v]
    - [ ] is Request info context maybe needed by now?
    - [ ] Use interfaces for testing
    - [ ] consistent parsable error messages
    - [ ] Make Errors use custom error codes, and also make a custom Server http Error handler
    - [ ] Define Route strings in Route Package
    - [ ] Use a nice logger
    - [ ] Give better http errors, not revealing internals, while logging debug information in the server logs - use package/helper func
    - [ ] Use types for int64 IDs, String IDs, and auth-esque tokens - especially those damn int64s everywhere
    - [ ] Pointer vs Value best practices, use automated scanner  
    
