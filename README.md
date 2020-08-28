# content-public-archives-web

## Introduction

Simplistic UI for exposing the zip archives for automated content export.

## Running locally

1) Generate SSL certificate

```sh
cd ./sslcert/
openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout server.key -out server.pem -subj "/C=US/CN=Example-Root-CA"
openssl x509 -outform pem -in server.pem -out server.crt
```

1) Configure local host

Add the following line to your `/ets/hosts` file

```txt
127.0.0.1 public-upp-exports.ft.com
```

1) Verify that you have nodejs and npm install locally:

```sh
node -v
npm -v
```

If the above commands does not yeld results, install node.

1) Run the following command:

```sh
npm install
```

1) Prepare configuration

Rename `.env.example` file to `.env` and populate the file using the information available [here](https://dashboard.heroku.com/apps/upp-exports/settings) under *Config Vars* section.

1) Run the following command:

```sh
npm run start:dev
```

If all env values are correctly populated the application should start on https://public-upp-exports.ft.com:8443 under the port defined in the `.env` file.

1) Run locally with heroku cli:

```
heroku local web
```


