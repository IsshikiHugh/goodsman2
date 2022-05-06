# Goodsman Database config

- config replica set:
  
```cfg

# master/slave.conf

# Where and how to store data.
storage:
  dbPath: <DBPATH>
  journal:
    enabled: true

# where to write logging data.
systemLog:
  destination: file
  logAppend: true
  path: <LOGPATH>

# network interfaces
net:
  port: <PORT>
  bindIp: 0.0.0.0

# how the process runs
processManagement:
  fork: true

### THIS ONE ###
replication:
  replSetName: goodsman2

security:
  authorization: enabled
  keyFile: <KEYFILEPATH>

 ```

- If you run master and slave db on same server, they CAN NOT have same path and port.

- use ``openssl rand -base64 745 > mongodb-keyfile`` to generate keyfile. It must be and **ONLY** be read-available to current user. You can use ``chmod`` and ``chown`` to config it.
-  Use ``mongod --config <path>`` to start mongo, use ``rs`` to config PRIMARY db. Then create a user who have ``readWrite`` authority in ``goodsman2`` db.


---

- after config MongoDB replica set, connect DB and init data with below command: 

```shell

# <shell>

mongo mongodb://<user>:<pwd>@<host>:<port>,<host>:<port>/goodsman2?tls=false&readPreference=secondaryPreferred&authSource=goodsman2&replicaSet=goodsman2

# <mongo shell>

use goodsman2
db.createCollection("goods")
db.goods.insertOne({
    "_id":   "gid",
    "name":  "bottom",
    "lore":  "this is the bottom of box",
    "msg":   "the first item in this program",
    "num":   1,
    "price": 9.99,
    "auth":  3,
    "image": "data//image:nil"
})

db.createCollection("employees")
db.employees.insertOne({
    "_id":   "default_group_1",
    "name":  "employee",
    "auth":  -1,
    "money": 500,
})
db.employees.insertOne({
    "_id":   "default_group_2",
    "name":  "admin",
    "auth":  -1,
    "money": 1000,
})
db.employees.insertOne({
    "_id":   "default_group_3",
    "name":  "super_admin",
    "auth":  -1,
    "money": 5000,
})

db.createCollection("records_hang")
db.records_hang.createIndex({
    "eid": 1,
    "gid": 1,},{
    "unique": true
})

db.createCollection("records_done")
db.records_done.createIndex({
    "eid": 1,
    "gid": 1,},{
    "unique": true
})
```

- If your server still can't connect database, add ``--verbose`` and look over logs.