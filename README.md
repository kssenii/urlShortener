**RUN**

```bash
$ docker-compose up --build
```


**USE**

-    Genearte url shortener
```bash
$ curl localhost:9090/encode -X POST -d '{"url":"http://www.vk.com"}'
$ {"url":"http://localhost:9090/aQ"}

```
-    Add custom url shortener
```bash
$ curl localhost:9090/encode -X POST -d '{"url":"http://www.google.com", "short":"googi"}'
$ {"url":"http://localhost:9090/googi"}

```
-    Decode url shortener
```bash
$ curl localhost:9090/decode -X POST -d '{"url":"http://localhost:9090/googi"}'
$ {"url":"http://www.google.com"}
```
-    Redirect
Open url in your browser http://localhost:9090/googi


**TESTS**

See *request_test.go*
