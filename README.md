**RUN**

```bash
$ docker-compose up --build
```

**USE**

#### Genearte url shortener
```bash
$ curl localhost:9090/encode -X POST -d '{"url":"http://www.google.com"}'
$ {"url":"http://localhost:8080/aQ"}

```
#### Decode url shortener
```bash
$ curl localhost:9090/decode -X POST -d '{"url":"http://localhost:8080/aQ"}'
$ {"url":"http://www.google.com"}

```

#### Redirect
Open url in your browser http://localhost:8080/aQ
