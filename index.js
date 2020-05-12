require('dotenv').config();
const http = require('http');

const cb = (app, logger) => {
  const server = http.createServer(app);
  const bind = server.listen(process.env.PORT, () => logger.info(`Listening on port ${bind.address().port}!`));
}

if (process.env.AUTH === 'external') {
  const { register } = require('./src/external')
  register(cb);
} else {
  const { register } = require('./src/internal')
  register(cb);
}
