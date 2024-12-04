
## Introduction | Bee-Aware
Bee-Aware is a Project built with Encore.dev and Svelte.
Its target is to create an easy to setup a camera / motion-detection-system. Its main purpose (as you mightve already guessed from the name) is to detect the activity from bees. But it also can serve as a cheap and easy to use security system.


## Prerequisites 

**Install Encore:**
- **macOS:** `brew install encoredev/tap/encore`
- **Linux:** `curl -L https://encore.dev/install.sh | bash`
- **Windows:** `iwr https://encore.dev/install.ps1 | iex`
  
**Docker:**
1. Install [Docker](https://docker.com)
2. Start Docker


## Run app locally

Before running your application, make sure you have Docker installed and running. Then run this command from your application's root folder:

```bash
encore run
```
To use the Slack integration, set the Slack Webhook URL (see tutorial above):

```bash
encore secret set --type local,dev,pr,prod SlackWebhookURL
```

**Note:** Cron Jobs do not execute when running locally.

## View the frontend

While `encore run` is running, head over to [http://localhost:4000/frontend/](http://localhost:4000/frontend/) to view the frontend for your uptime monitor.

## Using the API

Check if a given site is up (defaults to 'https://' if left out):
```bash
curl 'http://localhost:4000/ping/google.com'
```

Add a site to be automatically pinged every 1 hour:
```bash
curl 'http://localhost:4000/site' -d '{"url":"google.com"}'
```

Check all tracked sites immediately:
```bash
curl -X POST 'http://localhost:4000/check-all'
```

Get the current status of all tracked sites:
```bash
curl 'http://localhost:4000/status'
```

## Local Development Dashboard

While `encore run` is running, open [http://localhost:9400/](http://localhost:9400/) to access Encore's [local developer dashboard](https://encore.dev/docs/go/observability/dev-dash).

Here you can see traces for all requests you've made, see the application architecture diagram, and see API documentation in the Service Catalog.

## Connecting to databases

You can connect to your databases via psql shell:

```bash
encore db shell <database-name> --env=local --superuser
```

Learn more in the [CLI docs](https://encore.dev/docs/go/cli/cli-reference#database-management).

## Deployment

### Self-hosting

See the [self-hosting instructions](https://encore.dev/docs/go/self-host/docker-build) for how to use `encore build docker` to create a Docker image and configure it.

## Testing

```bash
encore test ./...
```
