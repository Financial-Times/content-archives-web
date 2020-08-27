var AWS = require("aws-sdk");

AWS.config.update({
    region: process.env.DYNAMODB_REGION,
    endpoint: process.env.DYNAMODB_ENDPOINT
});

var docClient = new AWS.DynamoDB.DocumentClient();

const getUserAllowlist = uuid => new Promise((resolve, reject) => {
  docClient.get({
    TableName: process.env.DYNAMODB_TABLE_NAME,
    Key: { uuid }
  }, (err, data) => {
    if (err) {
      return reject(err);
    }
    const item = (typeof data.Item != "undefined") ? data.Item : {}
    const allowList = (typeof item.data != "undefined") ? item.data.split(",").map(i => i.trim()) : [];
    resolve(allowList);
  });
});

module.exports = {
  getUserAllowlist
}
