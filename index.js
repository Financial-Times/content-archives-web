require('dotenv').config();
const http = require('http');
const { register } = process.env.AUTH === 'external' ? require('./src/external') : require('./src/internal');

const cb = (app, logger) => {
  const server = http.createServer(app);
  const bind = server.listen(process.env.PORT, () => logger.info(`Listening on port ${bind.address().port}!`));
};

register(cb);
