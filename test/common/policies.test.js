const { expect } = require('chai');
const { applyPolicies } = require('../../src/common/policies');
const archive = require('../__mocks__/archieve.json');

describe('Filtering archives with', () => {
  describe('`applyPolicies` returns correct count when', () => {
    const defaultPolicies = {
      withConcept: false,
      withAllVersions: true,
    };
    const testCases = [{
      policy: {
        since: 1900,
      },
      expectCount: 7,
      testName: '`since` date before first archive',
      message: 'does not include last-30-days.',
    }, {
      policy: {
        since: 2500,
      },
      expectCount: 0,
      testName: '`since` date after last archive',
      message: 'does not include any',
    }, {
      policy: {
        since: 1900,
        withAllVersions: false,
      },
      expectCount: 4,
      testName: '`since` date before first archive and without all versions included',
      message: 'does not return all-versions archive',
    },
    {
      policy: {
        since: 1900,
        withAllVersions: false,
        withConcept: true,
      },
      expectCount: 5,
      testName: '`since` date before first archive, without all versions, but concepts included',
      message: 'does not return all-versions archive',
    },
    {
      policy: {
        since: 1900,
        withAllVersions: true,
        withConcept: true,
      },
      expectCount: 8,
      testName: '`since` date before first archive with all versions and concepts included',
      message: 'returns all content',
    },
    {
      policy: {
        since: 2018,
        withAllVersions: true,
        withConcept: true,
      },
      expectCount: 7,
      testName: '`since` date after 2017 (2018 included) with all versions and concepts included',
      message: 'returns all content without 2017',
    }];

    testCases.forEach((testCase) => {
      it(`${testCase.testName}`, () => {
        const results = applyPolicies(archive, { ...defaultPolicies, ...testCase.policy });
        expect(results.length).to.equal(testCase.expectCount, testCase.message);
      });
    });
  });
});
