const { join } = require('path');
const express = require('express');
const session = require('cookie-session');
const OktaMiddleware = require('@financial-times/okta-express-middleware');
const logger = require('@financial-times/n-logger').default;
const { listArchives, downloadArchive } = require('../aws/s3-service');
const healthCheckMiddleware = require('../middleware/health-checks');
const {
  downloadArchiveErrorLog,
  listArchivesErrorLog,
  errorMessage,
  privateContact,
} = require('../common/messages.json');

const placeholderRegEx = /{placeholder}/i;

const register = (cb) => {
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
  app.use('/static', express.static(join(process.cwd(), 'static')));
  app.use(healthCheckMiddleware);
  app.use(okta.router);
  app.use(okta.ensureAuthenticated());
  app.use(okta.verifyJwts());

  const error = (res, err, msg) => {
    logger.error(msg, err);
    res.status(500).send(errorMessage.replace(placeholderRegEx, privateContact));
  };

  app.get('/', (_, res) => {
    listArchives
      .then((archives) => res.render('index', { archives }))
      .catch((err) => error(res, err, listArchivesErrorLog));
  });

  app.get('/download/:prefix/:name', (req, res) => {
    const { prefix, name } = req.params;
    downloadArchive(join(prefix, name), res, (err) => error(res, err, downloadArchiveErrorLog.replace(placeholderRegEx, privateContact)));
  });

  // Start the server
  cb(app, logger);
};

module.exports = {
  register,
};
