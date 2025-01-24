# game app

### migrations:
```bash
 go install github.com/rubenv/sql-migrate/...@latest
 
  sql-migrate status -env="production" -config=repository/mysql/dbconfig.yml
 
 sql-migrate up -env="production" -config=drepository/mysql/dbconfig.yml -limit=1

 sql-migrate down -env="production" -config=repository/mysql/dbconfig.yml -limit=1
```