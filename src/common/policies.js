const conceptsMatch = /^(archives)?(.*)-(concepts)/i;
const conceptsFilter = (archive) => (conceptsMatch.test(archive.name));

const sincePattern = /^(archives)?(.*)-(\d{4})/i;
const sinceFilter = (archives, since) => archives
  .filter((archive) => (sincePattern.test(archive.name) && archive.name.match(sincePattern)[3] >= since));

const allVersionsPattern = /^(archives)?(.*)-(all-versions)/i;
const allVersionsFilter = (archive) => (!allVersionsPattern.test(archive.name));

const applyPolicies = (archives, policy) => {
  const {
    withAllVersions, withConcept, since,
  } = policy;

  let allowedArchives = Array.prototype.concat.apply([], sinceFilter(archives, since));

  if (withConcept) {
    allowedArchives = allowedArchives.concat(archives.filter(conceptsFilter));
  }

  if (!withAllVersions) {
    allowedArchives = allowedArchives.filter(allVersionsFilter);
  }

  return allowedArchives;
};

const hasPolicyAccess = (item, policy) => {
  const sampleArray = [{ name: item }];
  const filteredArray = applyPolicies(sampleArray, policy);

  return filteredArray.length > 0;
};

module.exports = {
  applyPolicies,
  hasPolicyAccess,
};
