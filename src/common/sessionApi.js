const https = require('https');

const options = {
  hostname: 'api.ft.com',
  port: 443,
  path: '',
  method: 'GET',
  headers: {
    'x-api-key': process.env.SESSION_API_KEY,
  },
};

const getSession = (sessionId) => new Promise((resolve, reject) => {
  const params = { ...options, path: `/sessions/s/${sessionId}` };
  const req = https.request(params, (res) => {
    res.on('data', (d) => {
      const result = JSON.parse(d);
      if (res.statusCode >= 200 && res.statusCode < 300) {
        resolve(result);
        return;
      }

      let error;

      if (Object.prototype.hasOwnProperty.call(result, 'sessionInvalidReason')) {
        error = new Error('Corrupted session token.');
      } else if (Object.prototype.hasOwnProperty.call(result, 'message')) {
        error = new Error('Unauthorized api call.');
      } else {
        error = new Error('Unknown response.');
      }

      reject(error);
    });
  });

  req.on('error', (error) => {
    reject(error);
  });

  req.end();
});

module.exports = {
  getSession,
};
