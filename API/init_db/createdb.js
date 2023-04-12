/////
// Create a new Database
// Only invoked by an administrative process, when a new customer
// subscribes
//////

try {
  host;
} catch(e) {
  host = "localhost:27017"
}
var m = new Mongo(host)

try {
  isTest;
} catch(e) {
  isTest = false
}

if (isTest) {
  try {
    DB_NAME;
  } catch(e) {
    DB_NAME = "AutoTest"
  }
} else {
  DB_NAME;
  CUSTOMER_RECORDS_DB;
  ADMIN_USER;
  ADMIN_PASS;

  //Authenticate first
  var authDB = m.getDB("test")
  authDB.auth(ADMIN_USER,ADMIN_PASS);

  //Update customer record table
  var odb = m.getDB(CUSTOMER_RECORDS_DB)
  odb.customer.insertOne({"name": DB_NAME});
}

var db = m.getDB("ogree"+DB_NAME)
db.createCollection('account');
db.createCollection('domain');
db.createCollection('site');
db.createCollection('building');
db.createCollection('room');
db.createCollection('rack');
db.createCollection('device');

//Template Collections
db.createCollection('room_template');
db.createCollection('obj_template');
db.createCollection('bldg_template');

//Group Collections
db.createCollection('group');

//Nonhierarchal objects
db.createCollection('ac');
db.createCollection('panel');
db.createCollection('cabinet');
db.createCollection('corridor');
db.createCollection('sensor');

//Stray Objects
db.createCollection('stray_device');
db.createCollection('stray_sensor');

//Enfore unique Tenant Names
db.domain.createIndex( {parentId:1, name:1}, { unique: true } );

//Enforce unique children
db.site.createIndex({name:1}, { unique: true });
db.building.createIndex({parentId:1, name:1}, { unique: true });
db.room.createIndex({parentId:1, name:1}, { unique: true });
db.rack.createIndex({parentId:1, name:1}, { unique: true });
db.device.createIndex({parentId:1, name:1}, { unique: true });
//Enforcing that the Parent Exists is done at the ORM Level for now

//Make slugs unique identifiers for templates
db.room_template.createIndex({slug:1}, { unique: true });
db.obj_template.createIndex({slug:1}, { unique: true });
db.bldg_template.createIndex({slug:1}, { unique: true });

//Unique children restriction for nonhierarchal objects and sensors
db.ac.createIndex({parentId:1, name:1}, { unique: true });
db.panel.createIndex({parentId:1, name:1}, { unique: true });
db.cabinet.createIndex({parentId:1, name:1}, { unique: true });
db.corridor.createIndex({parentId:1, name:1}, { unique: true });

//Enforce unique children sensors
db.sensor.createIndex({parentId:1, type:1, name:1}, { unique: true });

//Enforce unique Group names 
db.group.createIndex({parentId:1, name:1}, { unique: true });

//Enforce unique stray objects
db.stray_device.createIndex({parentId:1,name:1}, { unique: true });
db.stray_sensor.createIndex({name:1}, { unique: true });

//Create a default domain
db.domain.insertOne({name: DB_NAME, hierarchyName: DB_NAME, category: "domain", attributes:{color:"ffffff"}, description:[], createdData: new Date(), lastUpdated: new Date()})