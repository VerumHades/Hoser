// mongo-init.js

// The database you want to create
const dbName = "apiDatabase";
const dbUser = "apiUser";
const dbPassword = "apiUserPassword";

// Get a handle to the database
db = db.getSiblingDB(dbName);

// Create collections
db.createCollection("users");
db.createCollection("listings");
db.createCollection("libraries");

// Optional: create indexes for the collections if you need
// db.users.createIndex({ email: 1 }, { unique: true });
// db.listings.createIndex({ title: 1 });
// db.libraries.createIndex({ userId: 1, listingId: 1 });

// Create a scoped MongoDB user for this database
db.createUser({
  user: dbUser,
  pwd: dbPassword,
  roles: [
    { role: "readWrite", db: dbName },
  ]
});