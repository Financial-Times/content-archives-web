const { join } = require('path');
const express = require('express');
const logger = require('@financial-times/n-logger').default;
const cookieParser = require('cookie-parser');
const { get } = require('../common/https');
const { listArchives, downloadArchive } = require('../common/s3-service');
const healthCheckMiddleware = require('../common/health-checks');
const messages = require('../common/messages.json');


const register = (cb) => {
  const app = express();

  app.set('view engine', 'ejs');
  app.use('/static', express.static(join(process.cwd(), 'static')));
  app.use(healthCheckMiddleware);
  app.use(cookieParser())

  const error = (res, err, msg) => {
    logger.error('Error retrieving content from Amazon S3', err);
    res.status(500).send(msg);
  };

  const redirectToLogin = (res) => {
    res.setHeader('Location', `https://accounts.ft.com/login?location=${process.env.APP_BASE_URL}`);
    res.setHeader('Cache-Control', 'private, max-age=0, no-cache, no-store, must-revalidate');
    res.setHeader('Pragma', 'no-cache');
    res.setHeader('Expires', '0');
    res.status(302).end();
  };

  app.use((req, res, next) => {
      if (req.cookies['FTSession_s'] === undefined) {
        redirectToLogin(res);
      } else {
        next();
      }
  });

  app.use((req, res, next) => {
    const sessionId = req.cookies['FTSession_s'];

    get(sessionId)
      .then((userSession) => {
        req.user = userSession.uuid;
        next();
      })
      .catch((err) => error(res, err, 'Failed to load membership session'))
  });

  app.get('/', (req, res) => {
    console.log(req.user);

    listArchives
      .then((archives) => res.render('index', { archives }))
      .catch((err) => error(res, err, messages.listArchivesError));
  });

  app.get('/download/:prefix/:name', (req, res) => {
    const { prefix, name } = req.params;
    downloadArchive(join(prefix, name), res, (err) => error(res, err, messages.listArchivesError));
  });

  // Start the server
  cb(app, logger);
};

module.exports = {
  register
}
