## Gokuwiki

### Get Started
- Download `build/docker-compose.yml` then run `docker compose up -d`

### Todo
- ~~Message when saving (error/success)~~
- Auto push commit (partial done)
  - by ssh key
  - (in case remote existed)
  - custom branch
  - ~~by access token~~ (done)
  - ~~[fix CA in docker image](https://stackoverflow.com/questions/64462922/docker-multi-stage-build-go-image-x509-certificate-signed-by-unknown-authorit)~~
- Image upload
- Show history
- Unit Test
- Allow omit comment
- Create config object to store configuration, currently reading config value from os env everytime
- Re-style darkmode
- ~~fix bug saving while saving `/` path still success~~
- add `sitemap.xml`

### Development
- Run `go build`

### Release
- Run `make`
