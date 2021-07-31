// admin
db.auth("mongoadm", "mongoadm")

// user
userdb = db.getSiblingDB("notes_system")
userdb.createUser({
  "user": "nsuser",
  "pwd" : "nsuser",
  "roles": [
    { "role" : "readWrite", "db" : "notes_system"}
  ],
  "mechanisms": [ "SCRAM-SHA-1" ],
  "passwordDigestor": "client"
})
userdb.auth("nsuser", "nsuser")
