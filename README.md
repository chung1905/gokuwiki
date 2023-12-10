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

### Changelogs
- 0.6.2: `noindex` in `/edit/` pages
- 0.6.1: change from dark mode to light mode, add markdown editor
- 0.6.0: Success/failed message
- 0.5.2: Upgrade dependencies
- 0.5.1: Doesn't create empty commit while start
- 0.5.0: Add `robots.txt`
- 0.4.3: Fix bug can't create wiki in directory
- 0.4.2: Remove `InsecureSkipTLS: true`
- 0.4.0: Add Cloudflare Turnstile captcha
- 0.3.1: `InsecureSkipTLS: true`
- 0.3: Auto push to repo by access token
- 0.2.3: Add commit message
- 0.2.2-1: Add buttons
- 0.2.1: Add header
- 0.2: commit to git after edit/delete
- 0.1: view/edit wiki
