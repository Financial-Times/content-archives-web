const { join } = require('path');
const express = require('express');
const session = require('cookie-session');
const cookieParser = require('cookie-parser');
const logger = require('@financial-times/n-logger').default;
const { listArchives, downloadArchive } = require('./src/s3-service');
const { userLoggedInMiddleware, userAllowListMiddleware } = require('./src/user-service');
const healthCheckMiddleware = require('./src/health-checks');
const messages = require('./src/messages.json');

const app = express();

app.use(session({
  secret: process.env.SESSION_SECRET,
  maxAge: 12 * 3600 * 1000, // 12 hours is the required age from our cyber-security-team
  httpOnly: true
}));

app.set('view engine', 'ejs');
app.use('/static', express.static(join(__dirname, 'static')));
app.use(healthCheckMiddleware);
app.use(cookieParser());
app.use(userLoggedInMiddleware);
app.use(userAllowListMiddleware);

const error = (res, err, msg) => {
  logger.error('Error retrieving content from Amazon S3', err);
  res.status(500).send(msg);
};

app.get('/', (req, res) => {
  listArchives
    .then(archives => archives.filter(i => req.allowList.indexOf(i.name) >= 0))
    .then(archives => res.render('index', { archives }))
    .catch(err => error(res, err, messages.listArchivesError));
});

app.get('/download/:prefix/:name', (req, res) => {
  const { prefix, name } = req.params;
  const fileName = join(prefix, name);
  if (req.allowList.indexOf(fileName) < 0) {
    return res.status(401).send('You are not authorized to access this file');
  }
  downloadArchive(fileName, res, err => error(res, err, messages.listArchivesError));
});

module.exports = { app };
