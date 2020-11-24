# Serverless Songs Application

This project uses serverless.com framework to manage two services, one resides at `services/tracks-api` (golang) and
holds a series of endpoints for managing tracks (e.g. songs). The other at `services/auth` (node), containing the required logic
to talk with Auth0 as a provider for user authentication.

The `client/` folder holds a simple React SPA for the UI bits.

As is, the application is able to handle audio file uploading to S3 as well as saving records at DynamoDB. Logs
for both services can be found at CloudWatch.

## Project goal

The goal of the project was to play with AWS services, mainly exploring the potential of FaaS.

Additionally, it lays here an interesting boilerplate for future use.

## Thoughts on FaaS

Although in its early stages, FaaS shows tremendous potential. Serverless.com simplifies the provisioning quite a bit, turning it
easier for developers to work with multiple lambda functions per service. I got a few pros & cons in the whole experience developing
two applications with this setup.

### Pros

- Managing one `serverless.yml` per service seems logical and draws the line well regarding responsibility of each service.
- FaaS "forces" a clean architecture given its functional nature, the system tends to grow cleaner.
- Horizontal-scalable by default, though not necessarily a setup and forget process.
- Great solution for quick computations and event-based procedures.
- For a small application, deploying and exploring after a few changes wasn't time-consuming, considering you can deploy just
the lambda function you're working at.

Consider that a service here is a group of lambda functions. Each service can be deployed separately, and every lambda function within
these services can also be deployed independently, which is a quite liberating experience out-of-the-box.

By event-based I mean that quite a few AWS services can emit events, for instance, after a successful S3 upload you could call another
lambda function, or append to a SQS queue, sky is the limit here.

### Cons

- It's definitely _not_ a silver bullet solution. Long running and/or high memory consumption 
procedures are not very welcome (look for the current lambda function limitations).
- Offline setup is still a bit clunky. Golang lambda functions at the time of writing cannot run locally, 
and still, local environment is just a simplified version of what's being hosted at AWS for now. Which leads to..
- Integration testing is a bit odd considering you'll be literally paying to run each test suite. And considering
everything sits in the cloud, one needs to think about different environments for each suite to run as well as different
development environments for each developer, which can also be a positive, but expensive point.

## Setup

Make sure you already have AWS credentials setup in your machine, then go through the following. AWS resources will be 
generated automatically through Cloud Formation.

### Backend

To deploy an application run the following commands:

```
cd services/tracks-api
make build && sls deploy -v
```

```
cd services/auth
npm install
sls deploy -v
```

### Frontend

To run a client application first edit the `client/src/config.ts` file to set correct parameters 
(you'll need to setup an Auth0 application), then run the following commands:

```
cd client
npm install
npm run start
```

That's what you should see if you succcessfully login with Auth0 (and upload songs):

![UI](https://github.com/oswaldoferreira/serverless-songs/blob/master/simple-ui.png)
