const { join } = require('path');
const express = require('express');
const logger = require('@financial-times/n-logger').default;
const cookieParser = require('cookie-parser');
const { getSession } = require('../common/sessionApi');
const { listArchives, downloadArchive } = require('../aws/s3-service');
const { getPolicyByUser } = require('../aws/dynamo-service');
const { applyPolicies, hasPolicyAccess } = require('../common/policies');
const healthCheckMiddleware = require('../middleware/health-checks');
const {
  downloadArchiveErrorLog,
  listArchivesErrorLog,
  gettingUserSessionErrorLog,
  listingUserPoliciesErrorLog,
  noPoliciesForUserLog,
  errorMessage,
  forbiddenMessage,
  forbiddenResourceAccessLog,
  publicContact,
} = require('../common/messages.json');

const placeholderRegEx = /{placeholder}/i;
const resourcePlaceholderRegEx = /{resource_placeholder}/i;

const register = (cb) => {
  const app = express();

  app.set('view engine', 'ejs');
  app.use('/static', express.static(join(process.cwd(), 'static')));
  app.use(healthCheckMiddleware);
  app.use(cookieParser());

  const error = (res, err, msg) => {
    logger.error(msg, err);
    res.status(500).send(errorMessage.replace(placeholderRegEx, publicContact));
  };

  const forbidden = (res, msg) => {
    logger.warn(msg);
    res.status(403).send(forbiddenMessage);
  };

  const redirectToLogin = (res) => {
    res.setHeader('Location', `https://accounts.ft.com/login?location=${process.env.APP_BASE_URL}`);
    res.setHeader('Cache-Control', 'private, max-age=0, no-cache, no-store, must-revalidate');
    res.setHeader('Pragma', 'no-cache');
    res.setHeader('Expires', '0');
    res.status(302).end();
  };

  app.use((req, res, next) => {
    if (req.cookies.FTSession_s === undefined) {
      redirectToLogin(res);
    } else {
      next();
    }
  });

  app.use((req, res, next) => {
    const sessionId = req.cookies.FTSession_s;

    getSession(sessionId)
      .then((userSession) => {
        req.userId = userSession.uuid;
        next();
      })
      .catch((err) => error(res, err, gettingUserSessionErrorLog.replace(placeholderRegEx, sessionId)));
  });

  app.use((req, res, next) => {
    const { userId } = req;

    getPolicyByUser(userId)
      .then((userPolicies) => {
        if (userPolicies) {
          req.userPolicies = userPolicies;
          next();
        } else {
          forbidden(res, noPoliciesForUserLog.replace(placeholderRegEx, userId));
        }
      })
      .catch((err) => error(res, err, listingUserPoliciesErrorLog.replace(placeholderRegEx, userId)));
  });

  app.get('/', (req, res) => {
    const { userPolicies } = req;

    listArchives
      .then((archives) => res.render('index', { archives: applyPolicies(archives, userPolicies) }))
      .catch((err) => error(res, err, listArchivesErrorLog.replace(placeholderRegEx, publicContact)));
  });

  app.get('/download/:prefix/:name', (req, res) => {
    const { prefix, name } = req.params;
    const { userPolicies, userId } = req;
    const fullName = join(prefix, name);

    if (hasPolicyAccess(fullName, userPolicies)) {
      downloadArchive(fullName, res, (err) => error(res, err, downloadArchiveErrorLog.replace(placeholderRegEx, publicContact)));
    } else {
      forbidden(res, forbiddenResourceAccessLog.replace(placeholderRegEx, userId).replace(resourcePlaceholderRegEx, name));
    }
  });

  // Start the server
  cb(app, logger);
};

module.exports = {
  register,
};
