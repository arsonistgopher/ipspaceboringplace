## CreateDB

Code based generation for a SQLite DB with service based variables.

## ServeDB

REST webserver that serves variable creation, deletions and gets.

```bash`
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"cuid":"111", "sid":"111"}' \
  localhost:1323/vars

curl -X GET \
  -H 'Content-Type: application/json' \
  localhost:1323/vars/111n111

curl -X DELETE \
  -H 'Content-Type: application/json' \
  localhost:1323/vars/111n111
```
