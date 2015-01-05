This repository is my playground for experimenting with go.  It provides a simple data model based on banks and accounts.  In its current form this playground app compares the (micro) frameworks [goji](https://github.com/zenazn/goji) and [httprouter](https://github.com/julienschmidt/httprouter).   
 
## Usage

#### Start goji endpoint
```bash
▶ webplay bank serve goji
```

#### Start httprouter endpoint
```bash
▶ webplay bank serve httprouter
```

#### Tests and benchmarks
The tests ensure that both the goji and httprouter implementations are at a basic level functioning as expected.  The benchmark test while simple appears to validate the results found [here](https://github.com/julienschmidt/go-http-routing-benchmark]\).  
```bash
▶ go test -v
INFO[0000] Testing goji
2015/01/03 21:41:16.421811 profile: cpu profiling enabled, cpu.pprof
=== RUN TestBankList
2015/01/03 21:41:16.421990 Starting Goji on [::]:8000
--- PASS: TestBankList (0.00s)
=== RUN TestBankListAccount
--- PASS: TestBankListAccount (0.00s)
PASS
INFO[0000] Testing httprouter
=== RUN TestBankList
--- PASS: TestBankList (0.00s)
=== RUN TestBankListAccount
--- PASS: TestBankListAccount (0.00s)
PASS
ok      github.com/jrcookli/webplay 0.019s
```

```bash
▶ go test -v -bench .
INFO[0000] Testing goji
2015/01/03 21:43:11.941444 profile: cpu profiling enabled, cpu.pprof
=== RUN TestBankList
2015/01/03 21:43:11.941665 Starting Goji on [::]:8000
--- PASS: TestBankList (0.00s)
=== RUN TestBankListAccount
--- PASS: TestBankListAccount (0.00s)
PASS
BenchmarkBankList        200       8474195 ns/op
INFO[0002] Testing httprouter
=== RUN TestBankList
--- PASS: TestBankList (0.00s)
=== RUN TestBankListAccount
--- PASS: TestBankListAccount (0.00s)
PASS
BenchmarkBankList        200       8143276 ns/op
ok      github.com/jrcookli/webplay 5.015s
```

*Note: I added cpu profiling using [this](https://github.com/davecheney/profile) utility but the results haven't made sense to me yet.  My next step will be to try and profile on another platform (not OSX).  See [bug](http://golang.org/pkg/runtime/pprof/).*

#### list banks
Using another shell, in the first you are serving an endpoint, run the following.
```bash
▶ webplay bank list
INFO[0000] bank ACME has 5028840.53 total deposits in 100 accounts
+-----------------------------------+------+-----------------+----------------+
|               ID                  | NAME | NUM OF ACCOUNTS | TOTAL DEPOSITS |
+-----------------------------------+------+-----------------+----------------+
| E0955-27BA-E53B-822B-547F4D2AAB12 | ACME |             100 |     5028840.53 |
+-----------------------------------+------+-----------------+----------------+
```

#### list accounts
```bash
▶ webplay bank accounts E0000955-27BA-E53B-822B-547F4D2AAB12
+--------------------------------------+--------------------+----------+
|                  ID                  |        NAME        | BALANCE  |
+--------------------------------------+--------------------+----------+
| 123411E8-F250-8073-BEBF-A9C3884BDA6D | James Harris       | 84171.86 |
| 28B038A6-93C7-0671-7E60-5B2099EE57D4 | Jacob Jackson      | 95436.56 |
| 390404FD-B22B-A5B1-29B1-34BEBE018BD8 | Abigail Harris     | 80005.36 |
| 47AA5864-AE42-DFBA-2135-44699E28221B | Addison Taylor     | 14140.33 |
| 609E182B-907E-80B8-57A8-05610E428409 | Jayden Harris      | 24761.33 |
| 6B42DDB9-D55F-9236-D247-FF92EBBFEB64 | Olivia Harris      | 72752.19 |
... 
```
*In the init method of common/models.go we create 100 accounts with random data.*

## Swagger

[Here](webplay/blob/master/public/swagger.yaml) is the swagger 2.0 spec for this toy app.  If you serve from the httprouter ```webplay serve httprouter``` you can access the UI through localhost:8000/api/ however if you host through goji it can be access through localhost:8000/.  (The difference is due to the way that the two solutions manage routes)

## Notes, observations, next steps 
* live reload for go web servers can be achieved using [gin](https://github.com/codegangsta/gin)
* I like [github.com/codegansta/cli](https://github.com/codegangsta/cli) however it does not really support commands of the form 'bank <bankID> account <accountID>'...  will investigate further
* CORS middleware added to enable swagger editor
    *   For the httprouter impl middleware was integrated by using [github.com/codegangsta/negroni](https://github.com/codegangsta/negroni)
* If a fullstack web framework is needed [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) is worth investigating
* If a more idiomatic approach to the web in Go is needed (less is more) then [github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) + [github.com/codegangsta/negroni](https://github.com/codegangsta/negroni) should be investigated
  

