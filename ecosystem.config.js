module.exports = {
  apps: [{
    name: 'content-archives-web',
    script: './index.js',
    instances: 1,
    autorestart: true,
    watch: false,
    max_memory_restart: '500M'
  }]
};
