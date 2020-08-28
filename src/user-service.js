const graphql = require('graphql.js');
const logger = require('@financial-times/n-logger').default;
const { stringify } = require('querystring');
const { getUserAllowlist } = require('./user-settings-service');

const getUserId = session => {
  const url = `${process.env.MEMBERSHIP_API_URL}/users/profile`;

  const graph = graphql(url, {
    alwaysAutodeclare: true,
    headers: {
      'x-api-key': process.env.MEMBERSHIP_API_KEY
    }
  });

  const req = graph.query(`
    userBySession(session: $session) {
      id
    }
  `);

  return req({ session })
    .then(({ userBySession }) => userBySession.id);

};

const userAllowListMiddleware = (req, res, next) => {
  getUserAllowlist(req.userId)
    .then(allowList => {
      if (!allowList || allowList.length === 0) {
        res.status(401).send('You are not authorized to access this page');
      }
      req.allowList = allowList;
      next();
    })
    .catch(err => {
      logger.error('Error loading credentiasl information', err);
      res.status(500).send('Error loading credentials information; ');
    });
};

const userLoggedInMiddleware = (req, res, next) => {

  const getSession = () => new Promise((resolve, reject) => {
    if (req.cookies.FTSession_s) {
      return resolve(req.cookies.FTSession_s);
    }
    reject();
  });

  const generateRedirectUrl = () => {
    const hostname = req.hostname.replace(/\/$/, '');
    const params = stringify({
      location: encodeURI(`https://${hostname}${process.env.PORT}/`)
    });

    return `${process.env.FT_ACCOUNTS_URL}/login?${params}`;
  };

  const redirectToLogin = () => {
    res.setHeader('Location', generateRedirectUrl());
    res.setHeader('Cache-Control', 'private, max-age=0, no-cache, no-store, must-revalidate');
    res.setHeader('Pragma', 'no-cache');
    res.setHeader('Expires', '0');
    res.status(302).end();
  };

  getSession(req)
    .then(session => getUserId(session))
    .then(userId => req.userId = userId)
    .then(() => next())
    .catch(err => {
      logger.error('Error loading credentiasl information', err);
      redirectToLogin();
    });
};

module.exports = {
  userLoggedInMiddleware,
  userAllowListMiddleware
};
