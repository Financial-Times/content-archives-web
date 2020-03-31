const S3 = require('aws-sdk/clients/s3');
const moment = require('moment');

const s3Client = new S3();

function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'K', 'M', 'GB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / (k ** i)).toFixed(dm)) + sizes[i];
}

const listArchives = new Promise((resolve, reject) => s3Client.listObjectsV2({
  Bucket: process.env.AWS_BUCKET_NAME,
  Prefix: process.env.AWS_BUCKET_PREFIX,
}, (err, data) => {
  if (err) {
    reject(err);
  } else {
    const result = data.Contents.map((s3Object) => ({
      name: s3Object.Key,
      lastModified: moment(s3Object.LastModified).format('YYYY-MM-DD hh:mm:ss'),
      size: formatBytes(s3Object.Size),
    }));
    resolve(result);
  }
}));

const bucketHealth = new Promise((resolve, reject) => s3Client.headBucket({
  Bucket: process.env.AWS_BUCKET_NAME,
}, (err, data) => {
  if (err) {
    reject(err);
  } else {
    resolve(data);
  }
}));

const downloadArchive = (key, res, errorCb) => s3Client.getObject({
  Bucket: process.env.AWS_BUCKET_NAME,
  Key: key,
}).on('httpHeaders', function handler(_, headers) {
  res.set('Content-Length', headers['content-length']);
  res.set('Content-Type', headers['content-type']);
  this.response
    .httpResponse
    .createUnbufferedStream()
    .pipe(res);
})
  .on('error', (err) => errorCb(err))
  .send();

module.exports = {
  listArchives,
  downloadArchive,
  bucketHealth,
};
