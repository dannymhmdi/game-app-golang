# game app

### migrations:
```bash
 go install github.com/rubenv/sql-migrate/...@latest
 
  sql-migrate status -env="production" -config=repository/mysql/dbconfig.yml
 
 sql-migrate up -env="production" -config=repository/mysql/dbconfig.yml -limit=1

 sql-migrate down -env="production" -config=repository/mysql/dbconfig.yml -limit=1
```

### Koanf :
With koanf we configure our app config :

1-Read default value and merge it to config struct:
```golang
	lErr := k.Load(confmap.Provider(map[string]interface{}{
		"auth.access_subject":  accessSubject,
		"auth.refresh_subject": refreshSubject,
	}, "."), nil)
```

2-Read config.yml file and merge it to config struct:
```golang
	err := k.Load(file.Provider("config.yml"), yaml.Parser())
	if err != nil {
		panic(err)
	}
```

3-Read envierment variables and merge it to config struct:

```golang
k.Load(env.Provider("GAMEAPP_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "GAMEAPP_")
		s = strings.ToLower(s)
		if index := strings.Index(s, "_"); index != -1 {
			s = s[:index] + s[index+1:]
		}
		return s
	}), nil)
```

