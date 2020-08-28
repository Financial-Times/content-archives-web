const AWS = require('aws-sdk');

AWS.config.update({
  region: process.env.DYNAMODB_REGION,
  endpoint: process.env.DYNAMODB_ENDPOINT
});

const getUserAllowlist = uuid => new Promise((resolve, reject) => {
  const docClient = new AWS.DynamoDB.DocumentClient();
  docClient.get({
    TableName: process.env.DYNAMODB_TABLE_NAME,
    Key: { uuid }
  }, (err, data) => {
    if (err) {
      return reject(err);
    }
    const item = (typeof data.Item !== 'undefined') ? data.Item : {};
    const allowList = (typeof item.data !== 'undefined') ? item.data : [];
    resolve(allowList.split(',').map(i => i.trim()));
  });
});

const dynamoDBHealth = () => new Promise((resolve, reject) => {
  const ddb = new AWS.DynamoDB({ apiVersion: '2012-08-10' });
  ddb.listTables({}, (err, data) => {
    if (err) {
      return reject(err);
    }
    if (!data.TableNames || data.TableNames.indexOf(process.env.DYNAMODB_TABLE_NAME) < 0) {
      return reject(new Error('user settings table does not exist'));
    }
    resolve();
  });
});

module.exports = {
  getUserAllowlist,
  dynamoDBHealth
};
