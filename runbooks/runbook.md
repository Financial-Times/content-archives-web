<!--
    Written in the format prescribed by https://github.com/Financial-Times/runbook.md.
    Any future edits should abide by this format.
-->
# UPP Exports

UI for exposing the zip archives for automated content/concept export.

## Code

upp-exports

## Primary URL

https://upp-exports.ft.com/

## Service Tier

Bronze

## Lifecycle Stage

Production

## Host Platform

Heroku

## Architecture

The application is a simple [express](https://expressjs.com/) server pulling the content of a specific s3 bucket and visualizing it in a static html page. Each file listed in the page is than available to download from the user. The application is considered healthy if access to the s3 bucket can be performed. The node server is started with [pm2](https://pm2.keymetrics.io/) with a default memory limit of 1G. If an unhandled exception occur and the node process is terminated it will be automatically restarted by `pm2`.

## Contains Personal Data

No

## Contains Sensitive Data

No

<!-- Placeholder - remove HTML comment markers to activate
## Can Download Personal Data
Choose Yes or No

...or delete this placeholder if not applicable to this system
-->

<!-- Placeholder - remove HTML comment markers to activate
## Can Contact Individuals
Choose Yes or No

...or delete this placeholder if not applicable to this system
-->

## Failover Architecture Type

None

## Failover Process Type

None

## Failback Process Type

None

<!-- Placeholder - remove HTML comment markers to activate
## Failover Details
Enter descriptive text satisfying the following:
The actions required to fail this system from one region to another. Either provide a set of numbered steps or a link to a detailed process that operations can follow.

...or delete this placeholder if not applicable to this system
-->

## Data Recovery Process Type

NotApplicable

## Data Recovery Details

The service does not store data, so it does not require any data recovery steps.

## Release Process Type

FullyAutomated

## Rollback Process Type

FullyAutomated

## Release Details

The application is automatically deployed to Heroku upon commit in the `mater` branch of the repository. Hence the rollback process will require revert commit in order to restore to a previous state of the application.

<!-- Placeholder - remove HTML comment markers to activate
## Heroku Pipeline Name
Enter descriptive text satisfying the following:
This is the name of the Heroku pipeline for this system. If you don't have a pipeline, this is the name of the app in Heroku. A pipeline is a group of Heroku apps that share the same codebase where each app in a pipeline represents the different stages in a continuous delivery workflow, i.e. staging, production.

...or delete this placeholder if not applicable to this system
-->

## Key Management Process Type

Manual

## Key Management Details

Once the key is the AWS account are rotated the new keys has to be added as configuration variables in Heroku.

## Monitoring

<p>Healthcheck:&nbsp;&nbsp;<a href="https://upp-exports.ft.com/__health" style="background-color: rgb(255, 255, 255);">https://upp-exports.ft.com/__health</a><br></p>

## First Line Troubleshooting

<p><a href="https://dashboard.heroku.com/apps/upp-exports" target="_blank" style="background-color: rgb(255, 255, 255);">Heroku link</a><br></p><p><span style="font-weight: 700;">Issue:&nbsp;service is unable to connect to AWS S3</span></p><ul><li>check that AWS credentials are correctly&nbsp;configured in Heroku&nbsp;</li></ul>

## Second Line Troubleshooting

Please contact the team via a message in the #upp-prod-incidents slack channel.