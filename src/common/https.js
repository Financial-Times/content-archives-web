const https = require('https')
const options = {
  hostname: 'api.ft.com',
  port: 443,
  path: '',
  method: 'GET',
  headers: {
    'x-api-key': process.env.SESSION_API_KEY
  }
}

const get = (sessionId) => new Promise((resolve, reject) => {
  params = Object.assign({}, options, { path: `sessions/s/${sessionId}` })
  const req = https.request(params, res => {
    console.log(`statusCode: ${res.statusCode}`)

    res.on('data', d => {
      resolve(JSON.parse(d));
    })
  })

  req.on('error', error => {
    reject(error);
  })

  req.end()
});

module.exports = {
  get
}