# booQ

[![GitHub release](https://img.shields.io/github/release/traPtitech/booQ-v3.svg)](https://GitHub.com/traPtitech/booQ-v3/releases/)
![CI](https://github.com/traPtitech/booQ-v3/workflows/CI/badge.svg)
![master](https://github.com/traPtitech/booQ-v3/workflows/master/badge.svg)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=traPtitech/booQ-v3)](https://dependabot.com)

management tool for equipment and book rental

## Development environment

### Setup with docker (compose)

#### First Up (or entirely rebuild)

```
$ docker compose up --build --watch
```

Now you can access to `http://localhost:8080` for booQ

And you can access booQ MariaDB by executing commands
`docker compose exec db bash` and `mysql -uroot -ppassword -Dbooq`

#### test

You can test this project

```
$ docker compose -f docker/test/docker-compose.yml up --abort-on-container-exit
```

#### Rebuild

`docker compose up --no-deps --build`

#### Destroy Containers and Volumes

`docker compose down -v`

### Setup VSCode

write it down in your `.vscode/settings.json`

```json
{
  "go.testEnvVars": {
    "MYSQL_DATABASE": "test"
  }
}
```
