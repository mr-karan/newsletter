# Newsletter

**This repo is WIP.**

> This repo contains the source code I use to self host my personal newsletter. The subscription form integrates with [listmonk](https://listmonk.app).


![img](screenshot.png)

## Setup

### Frontend

```sh
$ cd frontend
$ yarn # to add the required dependencies
$ yarn run start:dev # starts a local dev server
$ yarn watch # looks for changes in `css` directory and uses `tailwind` css to build non-minified css for development.
$ yarn build:prod # uses with `PostCSS` plugins like `css-nano` and `purge-css` which remove unwanted CSS + minify.
```

### Backend

```sh
$ cd backend
$ docker run -p 9379:6379 redis:latest # to run redis
$ cp config.toml.sample config.toml # and replace with your config
$ make fresh # starts a webserver with static assets packed into binary. API is available at `/api` and static assets at `/static`.
```