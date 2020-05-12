const { readFile } = require('fs');
const { join } = require('path');

const getPolicyByUser = (userId) => new Promise((resolve, reject) => {
  readFile(join(process.cwd(), 'policy.json'), (err, data) => {
    if (err) reject(err);
    else {
      const userPolicies = JSON.parse(data);
      if (userPolicies.uuid === userId) {
        resolve(userPolicies);
      } else {
        resolve(null);
      }
    }
  });
});

module.exports = {
  getPolicyByUser,
};
