# Drone Webhook Plugin
This plugin is listening to drone global webhook and perform termination to old build when there is a newer build created of the same branch.

## !Note: This does not apply to master branch, since we do not want to kill off any PR requests. 

## Prerequisite
1. Install go >= 1.18
2. [Create a secret](#secret)
3. [Create a env](#env)
4. Attached this secret into your drone server with the ENV 
  - DRONE_WEBHOOK_ENDPOINT=<http://... or https://...>
  - DRONE_WEBHOOK_SECRET=[Create a secret](#secret)

## Secret 
```
    openssl rand -hex 16
```

## .env
```diff
+  By default, this file should be stored in ./.env

DRONE_SECRET=<same secret key u put for DRONE_WEBHOOK_SECRET>
DRONE_BEARER=<Admin access user token>
DRONE_URL=<URL of your drone server>
MASTER_BRANCH=<Master branch to prevent killing>

```

## To run this script
- go run main.go

## To build the go project and run it 
- go build 
- ./drone-webhook