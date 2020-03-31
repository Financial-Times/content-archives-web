require('dotenv').config();

const http = require('http');
const { join } = require('path');
const express = require('express');
const session = require('cookie-session');
const OktaMiddleware = require('@financial-times/okta-express-middleware');
const logger = require('@financial-times/n-logger').default;
const { listArchives, downloadArchive } = require('./src/s3-service');
const healthCheckMiddleware = require('./src/health-checks');
const messages = require('./src/messages.json');

const app = express();

const okta = new OktaMiddleware({
  client_id: process.env.OKTA_CLIENT_ID,
  client_secret: process.env.OKTA_CLIENT_SECRET,
  issuer: process.env.ISSUER,
  appBaseUrl: process.env.APP_BASE_URL,
  scope: 'openid name offline_access',
});

app.use(session({
  secret: process.env.SESSION_SECRET,
  maxAge: 12 * 3600 * 1000, // 12 hours is the required age from our cyber-security-team
  httpOnly: true,
}));

app.set('view engine', 'ejs');
app.use('/static', express.static(join(__dirname, 'static')));
app.use(healthCheckMiddleware);
app.use(okta.router);
app.use(okta.ensureAuthenticated());
app.use(okta.verifyJwts());

const error = (res, err, msg) => {
  logger.error('Error retrieving content from Amazon S3', err);
  res.status(500).send(msg);
};

app.get('/', (_, res) => {
  listArchives
    .then((archives) => res.render('index', { archives }))
    .catch((err) => error(res, err, messages.listArchivesError));
});

app.get('/download/:prefix/:name', (req, res) => {
  const { prefix, name } = req.params;
  downloadArchive(join(prefix, name), res, (err) => error(res, err, messages.listArchivesError));
});

const server = http.createServer(app);
const bind = server.listen(process.env.PORT, () => logger.info(`Listening on port ${bind.address().port}!`));
