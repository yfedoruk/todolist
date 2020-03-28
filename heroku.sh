#!/usr/bin/env bash
#heroku container:login
heroku container:push --app nameless-brook-78889 web
heroku container:release --app nameless-brook-78889 web
heroku logs --tailgc --app nameless-brook-78889