# Icinga Check Elastic
This is a learning repo for simple icinga checks in golang.

---
## Commands

**Max Command**

Help Output *(./icinga_check_elastic help mindocs)*:
```
Fails with exit code 2 , if doc count for a period is lower then defined treshold

Usage:
  check_elastic mindocs [flags]

Flags:
  -a, --auth string    basic auth for header authentication. format=username:password
  -e, --exit int       exit code to be used for fail (default 2)
      --field string   string of filed name like internal.created
  -h, --help           help for mindocs
      --index string   string of idex pattern or full index name like my-awesome-index-*
  -m, --min int        defines minimum amount of docs that are required when command fails (default 100)
  -p, --period int     sets mintues for period now - x minutes (default 30)
      --url string     string of url like http://localhost:9200
```

Example Call:
```shell
./icinga_check_elastic mindocs --url http://loclhost:9200 --index "my_index-*" --field my.time.field -a guest:guest
```

