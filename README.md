# content-archives-web

## Introduction
Simplistic UI for exposing the zip archives for automated content export.

## Running locally


1) Verify that you have nodejs and npm install locally:

```
node -v
npm -v
```

 If the above commands does not yeld results, install node using homebrew. There are plenty of example in the internet.

2) Run the following command:

```
npm install
```

3) Rename `.env.example` file to `.env` and populate the file using the information available [here](https://dashboard.heroku.com/apps/upp-exports/settings) under *Config Vars* section. For the Okta related secrets login into Vault, navigate to the proper folder based on the information from the following [section](https://github.com/Financial-Times/okta/wiki/Config-Guide#what-next). On how to get started with Vault refer to the following [documentation](https://github.com/Financial-Times/vault/wiki/Getting-Started-With-Vault).

4) Run the following command:

```
npm run start:dev
```

If all env values are correctly populated the application should start on localhost under the port defined in the `.env` file.

5) Run locally with heroku cli:

```
heroku local web
```