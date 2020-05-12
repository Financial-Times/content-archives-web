const { join } = require('path');
const HealthCheck = require('@financial-times/health-check');
const expressWebService = require('@financial-times/express-web-service');
const { bucketHealth } = require('../aws/s3-service');

class S3BucketCheck extends HealthCheck.Check {
  run() {
    return bucketHealth
      .then(() => {
        this.ok = true;
        this.checkOutput = '';
        this.lastUpdated = new Date();
      })
      .catch((error) => {
        this.ok = false;
        this.checkOutput = error.message;
        this.lastUpdated = new Date();
      });
  }
}

const health = new HealthCheck({
  checks: [
    new S3BucketCheck({
      id: 'check-connectivity-to-s3',
      name: 'Check connectivity to AWS S3',
      panicGuide: 'https://runbooks.in.ft.com/upp-exports',
      technicalSummary: 'The service is unable to connect to AWS S3',
      businessImpact: 'Content and Concept archives won\'t be available for download',
      severity: 1,
      interval: 60000 * 10, // ten minute interval
    }),
  ],
});

const healthCheckMiddleware = expressWebService({
  about: {
    systemCode: 'upp-exports',
    name: 'UPP Daily Exports',
    description: 'Downloadable Content and Concept archives.',
    serviceTier: 'bronze',
  },
  healthCheck: health.checks(),
  manifestPath: join(__dirname, '../..', 'package.json'),
});

module.exports = healthCheckMiddleware;
