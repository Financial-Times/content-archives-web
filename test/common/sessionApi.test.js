const nock = require('nock');
const { expect } = require('chai');
const { getSession } = require('../../src/common/sessionApi');

const happyPathResponse = '{"uuid":"bar","creationTime":1589190178455,"rememberMe":true,"mfa":false}';
const forbiddenResponse = '{"message":"Forbidden"}';
const corruptResponse = '{"sessionInvalidReason":"CORRUPT"}';
const genericErrorResponse = '{"foo":"bar"}';

describe('Get session API', () => {
  describe('happy path', () => {
    beforeEach(() => {
      nock('https://api.ft.com')
        .get('/sessions/s/foo')
        .reply(200, happyPathResponse);
    });

    it('resolves promise', () => getSession('foo'));

    it('resolved promise with userInfo', async () => {
      const userInfo = await getSession('foo');

      expect(userInfo).to.have.property('uuid');
    });
  });

  describe('unhappy path', () => {
    let nockMock;

    beforeEach(() => {
      nockMock = nock('https://api.ft.com')
        .get('/sessions/s/foo');
    });

    it('reject promise when forbidden response', (done) => {
      nockMock.reply(401, forbiddenResponse);
      getSession('foo').catch((err) => {
        expect(err.message).to.not.be.undefined;
        done();
      });
    });

    it('reject promise when corrupt token response', (done) => {
      nockMock.reply(401, corruptResponse);
      getSession('foo').catch((err) => {
        expect(err.message).to.not.be.undefined;
        done();
      });
    });

    it('reject promise when generic error response', (done) => {
      nockMock.reply(500, genericErrorResponse);
      getSession('foo').catch((err) => {
        expect(err.message).to.not.be.undefined;
        done();
      });
    });

    it('reject promise when http error occur', (done) => {
      nockMock.replyWithError({
        message: 'something awful happened',
      });
      getSession('foo').catch((err) => {
        expect(err.message).to.not.undefined;
        done();
      });
    });
  });
});
