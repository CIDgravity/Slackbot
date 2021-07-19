
### Get started with TwinQuasar SlackBot

First you need to download all required libraries (listed in go module)

```
go mod download
```

If needed, you can change the Node Filecoin endpoint edit details inside ```backend/filecoin/filecoin.go```

* Change the addresse on line 26
* You can also add an auth token inside headers (if applicable)


### Run the watchdogs (check update from the chain and treat blocks)
This script is used to track height changed from the chain, treat blocks (decode message, get rewards, decode actor ...) and send all informations in a remote (or local) database

```
go run main_filecoin_watchdogs.go
```

### Run the notifier
This script is used as a database trigger: when new block is treated by the watchdog, it will get the database record to check if Slack notification need to be send to the miner

Today two notifications are available:
* When a miner mine a block, the associated reward will be send as a Slack message
* When a miner passed is PoST, he will receive a Slack message
* When a miner failed is PoST, he will also receive a Slack notification

```
go run main_trigger_notifiers.go
```

### Install and use local PostgreSQL engine / database

First you need to install PostgreSQL engine and setup users and databases

To install PostgreSQL

```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

By default, a new user is available with username postgres (and a default detabase with name postgres also). To make our server more secure, we will create a new user without the superuser rights. To make it, we can use the folowing commands :

```
sudo -u postgres createuser --interactive
```

When you have answers all the questions, we need to create a UNIX user with exactly same name of our new PostgreSQL user

```
sudo adduser twinquasar
```

Now, we need to create a database (this command means create the database slackbot under superuser postgres)

```
sudo -u postgres createdb slackbot
```

Change UNIX user and connect to the new created database

```
sudo -i -u twinquasar
psql -d slackbot
```

The last step to do, is to import the project database schema, for this, we can use the following commands
For now, the best way to achieve this, is using pgAdmin4, to install it, do :

```
curl https://www.pgadmin.org/static/packages_pgadmin_org.pub | sudo apt-key add
sudo sh -c 'echo "deb https://ftp.postgresql.org/pub/pgadmin/pgadmin4/apt/$(lsb_release -cs) pgadmin4 main" > /etc/apt/sources.list.d/pgadmin4.list && apt update'
sudo apt install pgadmin4
sudo /usr/pgadmin4/bin/setup-web.sh
```
