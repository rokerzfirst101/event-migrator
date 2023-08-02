# Event Migration Script

The following script is made to transfer the events from the OLD account to the NEW account for hiring motion.

## Environment Variables

```
APP_NAME: The app name to be used by Gofr
APP_VERSION: The app version 

DB_DIALECT: DB Dialect for the DB
DB_HOST: DB Host
DB_PORT: DB Port
DB_NAME: DB User Name
DB_PASSWORD: DB Password

OLD_PROJECT_CREDENTIALS: Credentials for the old gcp project
OLD_OAUTH_TOKENS: OAuth tokens for the old account

NEW_PROJECT_CREDENTIALS: Credentials for the new gcp project
NEW_OAUTH_TOKENS: OAuth tokens for the new account
```