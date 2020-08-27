require('dotenv').config();

const http = require('http');
const { join } = require('path');
const express = require('express');
const session = require('cookie-session');
const logger = require('@financial-times/n-logger').default;
const { listArchives, downloadArchive } = require('./src/s3-service');
const { getUserAllowlist } = require("./src/dynamodb-service");
const healthCheckMiddleware = require('./src/health-checks');
const messages = require('./src/messages.json');

const app = express();

app.use(session({
  secret: process.env.SESSION_SECRET,
  maxAge: 12 * 3600 * 1000, // 12 hours is the required age from our cyber-security-team
  httpOnly: true,
}));

app.set('view engine', 'ejs');
app.use('/static', express.static(join(__dirname, 'static')));
app.use(healthCheckMiddleware);

const error = (res, err, msg) => {
  logger.error('Error retrieving content from Amazon S3', err);
  res.status(500).send(msg);
};

app.get('/', (_, res) => {
  getUserAllowlist("2e9ae5fe-c02c-4e49-a704-617e871c82b8")
    .then(allowList => {
      listArchives
        .then(archives => archives.filter(i => allowList.indexOf(i.name) >= 0))
        .then((archives) => res.render('index', { archives }))
        .catch((err) => error(res, err, messages.listArchivesError));
    })
    .catch(err => error(res, err, "Unable to load a list of allowed files"));
});

app.get('/download/:prefix/:name', (req, res) => {
  const { prefix, name } = req.params;
  downloadArchive(join(prefix, name), res, (err) => error(res, err, messages.listArchivesError));
});

const server = http.createServer(app);
const bind = server.listen(process.env.PORT, () => logger.info(`Listening on port ${bind.address().port}!`));
