require('dotenv').config();

const https = require('https');
const fs = require('fs');
const privateKey = fs.readFileSync('sslcert/server.key', 'utf8');
const certificate = fs.readFileSync('sslcert/server.crt', 'utf8');
const credentials = { key: privateKey, cert: certificate };

const logger = require('@financial-times/n-logger').default;

const { app } = require('./app');

const server = https.createServer(credentials, app);
const bind = server.listen(process.env.PORT, () => {
  logger.info(`Listening on port ${bind.address().port}!`);
});
