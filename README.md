Twitter Rules
============================

## Running on Heroku

1. [![Deploy on Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)
1. Open the *Resources* tab in your app's dashboard: `https://dashboard.heroku.com/apps/YOUR APP NAME/resources`
1. Click on the *Heroku Scheduler* Add-on
1. Add a new scheduled task with `bin/twitter-rules -config config.yml` as the command
1. Upload a `config.yml` with your configuration: see [https://stackoverflow.com/a/40411074](https://stackoverflow.com/a/40411074)
1. Check your app's logs to see the scheduled tick: `https://dashboard.heroku.com/apps/YOUR APP NAME/logs`
