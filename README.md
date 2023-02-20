# SAM-API

## Purpose

The SAM-API is an interface to configure the mappings from BSCS
GL accounts to SAP OFI. It is supposed to be run in docker environment.
It is a REST server handling requests from frontend. It provides
access to the backend BSCS tables. The request and responses are
encoded in JSON format. The database is Oracle handled by GORP/GORACLE
so it uses the Oracle instant client which is a part of the release
of the docker image.

## Business Logical Model

The business models is based on 5 main entities:

 - **Account** -> **SAP_ACCOUNT** (create, read, update, delete)
 - **Order** -> **SAP_ACC_SEGM_ORDER_NUMBERS** (create, read, update, delete)
 - **DictionarySegment** -> **CUSTOMER_SEGMENT** (create, read, delete)
 - **DictionaryAccountBscs** -> **GLACCOUNTS** (BSCS base line setup, read only on view)
 - **DictionaryAccountSap** - **SAP_OFI_ACCOUNTS** (create, read, delete)

The backend tables are created in **BILLDB** database in **CGSYSADM** schema.

The database user is SAMAPI having all access rights to the backend tables.
The public synonyms give access to them.

The corresponding database tables are access in a CRUD model.
The operations on the tables are mainly bulk ones. The whole
segment is being selected and sent via API. Operation on the mapping
tables like update or delete are to be done by primary key access.
Operations on the dictionary tables are mostly bulk ones as the frontend
need the whole content of configuration fine tuned manipulation
is out of scope.

**Account** methods:

 - **/api/account POST**
 - **/api/account GET**
 - **/api/account DELETE** 
 - **/api/account/{status}/{release} GET** 
 - **/api/account/{status}/{release}/{account} PUT**
 - **/api/account/{status}/{release}/{account} PATCH**
 - **/api/account/{status}/{release}/{account} DELETE**
 - **/api/account/log GET**
 
**Order** methods:

 - **/api/order POST**
 - **/api/order GET**
 - **/api/order DELETE** 
 - **/api/order/{status}/{release} GET**  
 - **/api/order/{status}/{release}/{account}/{segment} PUT**
 - **/api/order/{status}/{release}/{account}/{segment} PATCH**
 - **/api/order/{status}/{release}/{account}/{segment} DELETE**
 - **/api/order/log GET**
 
**DictionarySegment** methods:

 - **/api/dictionary/segment POST**
 - **/api/dictionary/segment GET**
 - **/api/dictionary/segment DELETE**
 - **/api/dictionary/segment/{id} PUT**
 - **/api/dictionary/segment/{id} PATCH**
 - **/api/dictionary/segment/{id} DELETE**
 
**DictionaryAccountBscs** methods:

 - **/api/dictionary/account/sap GET** 

**DictionaryAccountSap** methods:

 - **/api/dictionary/account/sap POST application/xlsx**
 - **/api/dictionary/account/sap POST**
 - **/api/dictionary/account/sap GET** 
 - **/api/dictionary/account/sap DELETE**
 
**Release** methods:
 - **/api/release/new POST**
 - **/api/release/{release} POST** 
 - **/api/release/{release} DELETE**

The entities like **Account** and **Order** are versioned by status and
release id values. The status may be wither **W** like **Work** or **P** like
**Production**. An user can modify only entries in **W** status.
The release values for the entries in **W** status is 0. The entries
in **P **status may have release value > 0 linking them in versioned
packages. This makes versioning of the sets of mappings possible.

The status values are:

 - **W** - **Work**
 - **C** - **Control**
 - **P** - **Production**

The **GET** operation may use value "last" instead of explicit id thus
simplifying the access to the latest version in production.

A release from **W** to **P** is feasible only after the entries are
accepted by user with **Control** role. It means that all attributes to be modified
by **Control** are set. The release itself is an update of **W** -> **P**
and assignment of release value = max(release) of **Account**, **Order**.

The GET on /api/account and /api/order works in as special way as it fills
the main screen of online application. Only active entries are loaded, ie.
entries in W or C status with release 0 or P status with latest release
of all.

The most easy method used to perform the release is:

 - **/api/releas/new**

Is is changing the values of the attributes **STATUS** and **RELEASE_ID**.
The values **STATUS** is set according to the role of the user:

 - **Booker**: **W** -> **C**
 - **Control**: **C** -> **P**

The value **RELEASE_ID** of **Account** and **Order** is set to the
max **RELEASE_ID** of **Account** entity.

The **DELETE** on specific release may success only if all validDate 
values are in the future and the release was not used in mapping
report.

The release does a **POST** operation in two contexts. It can allocate
completely new release id or it can use existing one thus appending
the records being release to the exiting one.

The swagger schema is provided in the source code. It can be
started using make target swagger-ui provided that the swagger-ui
docker image is available. It uses the files in the directory:
**swagger**.

The event of release creation causes sending of a confirmation mail to
to the configured mail address from config. This can be overridden
by env variable ALERTMAILADDRESS. If ALERTMAILSERVERADDRESS is empty
the action is not performed, for example during system test.

The **/api/dictionary/account/sap POST** method may send the content 
as Excel file, compressed or not. The decoding of the payload 
is based on the values in the header:

 - **Content-Type**: **application/xlsx**
 - **Content-Type**: **application/json**

The decompression of the payload is triggered with the header value:

 - **Content-Encoding**: **gzip**

Currently only one format is available

The methods loading Order or Account may return the results of GET 
opration in various formats like json or csv or Excel. In order to 
trigger such a specific formating it is necessary to use in GET 
header values **Content-Type** set to:

 - **Content-Type**: **application/xlsx**
 - **Content-Type**: **application/csv** 
 - **Content-Type**: **application/json**
 
The addiitional attribute REC_VERSION is used to mark the records as
being handled by some session so other sessions should not touch them.

The operations on Orders or Account are logged in the LOG tables. 
The are having additional paramters like:

 - OPCODE: 'I', 'U', 'D'
 - OPDATE

The specific method may read the contents of the log tables.

## Authorization schema

The access to certain combinations of entities, methods and attributes
is based on the role. The hooker can create any object, may set most
of the attributes (upon creation or afterward), it may delete any
entity thus preparing work. The control can only set specific attributes
in **Account** or **Order**. 

Both of the Booker and Control may release a set of prepared entities
but only Control's release allocates new release id. The release
done by Booker just keeps the old 0 release id. 

The CRUD access rules for the entities are as follows:

 - **Account**
   - Only entries in status W and release id 0 can be created - default values
   - The release id must be unsigned integer or last
   - Control can not create accounts
   - Booker can not set attribute ofiSapWbsCode
   - Control can set only attributes ofiSapWbsCode, status, releaseId
   - The attribute validFromDate must be the 1st of month in future (proper cutoff date)
   - Everybody can read
   
 - **Order**
   - Only entries in status W and release id 0 can be created - default values
   - The release id must be unsigned integer or last
   - Booker can not set attribute orderNumber
   - Control can set only one attribute orderNumber
   - The attribute validFromDate must be the 1st of month in future
   - Booker can not delete
   - Everybody can read
 
 - **DictionarySegment**
 - **DictionaryAccountSap**
   - Everybody can create, read, update, delete, ie. no constraints
   
 - **DictionaryAccountBscs**   
   - Everybody can read
   - Nobody can create, update , delete as it is a view on GLACCOUNT
   
The users are to be authorized with LDAP server. This authorization
can be switched off by setting some of the LDAP env variables to empty
string.

## Authentication schema

The REST requests are authorized with JWT schema. The user login
request must be issued as the first one in order to obtain a JWT
token with encoded properties like user and role. In all subsequent
requests this token must be used in the header of the request.

The token contain encrypted claims like user and role so that
they are not sent in open format in the header of the request.
This private/public keys are stored in the directory: **keys**.

The token is invalidated with logoff request.

The credentials of the user are used to determine access rights.
The access rights are controlled on the level of entity, its attributes
and type of request.

The server must support **CORS** preflight requests as the online
is to be based on Java Script model. A swagger ui is using it
as well. No restrictions are provided in **CORS**.

The methods in this scope are:

 - **/api/user/login**
 
 Registers user with password and role. It uses LDAP server to authorize the user.
 The JWT token is created and returned to the client. The token has
 expiry date which is returned in the header or the response. The user
 is added to the current user list.
 
 - **/api/user/relogin**
 
 Enables relogin of the already registered user.
 
 - **/api/user/logoff**
 
 Enables explicit log off from the system. The user is removed from
 current user list.
 
 - **/api/user/info** 
 
 Returns the info about the user suing the token, giving details
 like user id and role id.

The 1st method is creating the token for give user and role which
can be one of:

 - **Booker**
 - **Control**
 - **Admin**

The 2nd method enables refresh of the token which is normally granted for 1 hour
this can be however overridden by config or env or command line options.
Expired toke must be replaced with a new one form explicit relogin or login
done once again. The expiry period is a configurable value so it can be extended
in the setup.

## API Service architecture

Each module is located in a separate directory of the source tree.

The modules are:

 - **models**

 It contains main entities common do all modules. They are the representation
 of the logical entities also known as business layer. 

 - **data**

 It contains structures of the representations of the database entities.
 They are used by repository layers and converted to the structures
 of the model entities handled subsequently internally.

 - **resources**

 It contains the structures used in API. It is adding the transport layers
 like envelopes used int the requests and responses. Normally, the main
 root of the structures here is data. All other structures are included
 from the data models.

 - **repository**

 It is the database communication layer. It used ORM whenever it can.
 Select operations can be implemented with dynamic queries but it is
 not recommended. 

 - **controllers**

 They server as handlers for CRUD operations on API. They receive
 the payloads from the requests, they use the repository layers
 and they provide the responses to the requests.

 - **routers**

 All endpoints, paths, methods are implemented here. It configures
 the mux/gorilla/negroni features and binds the handlers with
 paths. The authorization, log, etc. decorators are applied here.

 - **common**

 Common utilities used by all other sub-modules.

## System health check

The methods for health check of the system are:

 - **/api/system/version**
 - **/api/system/stat**
 - **/api/system/health**

The authorization is not needed in order to access them.

## Configuration and usage

The configuration of Oracle access is stored in the **config.json**
file. Normally it is to be placed in the current working
directory of teh process. The structure of this json file is like this:

```
{
	"ServerIPAddress"       : "0.0.0.0",
	"Port"                  : "8000",
	"RunPath"               : ".",
	"KeyPath"               : "keys",
	"Debug"                 : "0",
	"OracleDBUser"          : "SAMAPI",
	"OracleDBPassword"      : "SAMAPI",
	"OracleServiceName"     : "t17bill",
	"AlertMailAddress"      : "root@localhost",
	"AlertMailServerAddress": "localhost:25",
	"AlertMailSenderAddress": "sam@localhost",
	"LdapBase"              : "dc=corpo,dc=t-mobile,dc=pl",
	"LdapHost"              : "corpo.t-mobile.pl",
	"LdapPort"              : "389",
	"LdapBindDN"            : "ou=Uzytkownicy,ou=Standard,ou=Warszawa,dc=corpo,dc=t-mobile,dc=pl"
}
```

The Oracle service name string is to be configured in local TNS
file placed in the working directory. The **TNS_ADMIN** environment
variable must point to this location. The values from the setup file
may be overiden by invocation switches or by environemt variables.

The invocation options are the the following:

```
Usage of ./sam-api:
  -alertmailaddress string
    	Alert mail address (default "root@localhost")
  -alertmailsenderaddress string
    	Alert mail sender address (default "samapi@localhost")
  -alertmailserveraddress string
    	Alert SNMP server address
  -config string
    	Config json file if not in $RUNPATH/config.json (default "config.json")
  -debug string
    	Debug level
  -jwttokenvalidhours string
    	JWT token validity period in hours (default "1")
  -keypath string
    	Key path (default "keys")
  -ldapbase string
    	LDAP base 
  -ldapbinddn string
    	LDAP bind DN 
  -ldaphost string
    	LDAP host 
  -ldapport string
    	LDAP port
  -oracledbpassword string
    	Oracle DB password
  -oracledbuser string
    	Oracle DB user
  -oracleservicename string
    	Oracle service name
  -runpath string
    	Run path 
  -serveripaddress string
    	Server address 
  -serverport string
    	Server port 
  -v	Version check
```

The HTTP server is terminated gracefully by handling TERM UNIX signal.

As the API server is supposed to be run the Kubernetes enbvironment
the configuration values ma be overriden by environment variables.
The invocation options have however even higher pririty which can be usefull
in interractive testing.

The environment variables are:

 - **SERVERIPADDRESS**: on which interface addressis it to be run
 - **SERVERPORT**: port to override the one from config file 
 - **RUNPATH**: where the modules was started and where config.json file is located
 - **KEYPATH**: where the keys are located
 - **DEBUG**: start in verbose mode printing in log all values
 - **ORACLEDBUSER**: user id for ORacle connection
 - **ORACLEDBPASSWORD**: password of the user
 - **ORACLESERVICENAME**: TNS service names from ./oracle/tnsnames.ora file
 - **ALERTMAILADDRESS**: notification address about release events
 - **ALERTMAILSERVERADDRESS**: address of the mail server with SNM<P port id
 - **ALERTMAILSENDERADDRESS**: address used as the sender in notification mails
 - **JWTTOKENVALIDHOURS**: validity period for JWT token in hours
 - **LDAPBASE**: LDAP base string
 - **LDAPBINDN**: LDAP bind DN string
 - **LDAPHOST**: LDAP host
 - **LDAPPORT**: port for LDAP user check
 
The verride the values from config file.

## Error and status codes returned by API

The error codes are described in file **swagger/swagger.yaml**.

They are:

 - **500** - The errors in the method handlers are by default errored out as
 Internal Server Error.

 - **403** - acces to particular entity, attribute or method is forbudden for
 the users role.

 - **401** - Access denied is ussued when the token is missing or it is malformed
 or expored. The API client is sipposed to handle this situation by token refresh
 done by relogin action.

 - **201** - Object created after POST

 - **200** - Successfull operation

## Building and deployment

The GO version used is: go1.10.4 linux/amd64. It should be put in /usr/local/go
path. 

The repository must be places in the $GOPATH/src location. The recommended 
.bashrc should provide the following variables:

```
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
export CDPATH=".:$GOPATH/src"
```

The recommended build platform is Linux Ubuntu Disco Dingo. The necessary 
dependencies can be loaded with make target: deps.

The version of the server modules is treated by git tag in format x.y.z where
 
 - z is the feature release id
 - y is the bug fix release id
 - x is the major production release id

The tag is used in docker image build as well it is stored in VERSION file 
in local directory thus marking this version as releasable.

The Makefile contains the following useful targets:

 - **deps**
 
 Downloads go dependencies and installs them in the local system.
 
 - **keys**
 
 Builds private/public key files in keys directory to be used byAPI server
 to create and decrypt JWT tokens.
 
 - **sql**
 
 Installs tables in the Oracle database, currently T17 only.
 
 - **build**
 
 Compiles the main target and leaves it in the local directory.
 
 - **run**
 
 Runs the main target.
 
 - **clean**
 
 Cleans the artifacts.
 
 - **docker**
 
 Creates the docker image with version from git tag. It uses the oracle
 instant client install files from **oracle** directory and tnsnames.ora
 file with used oracle database connections.
 
 - **docker-run**
 
 Starts the API server in docker image.
 
 - **docker-push**
 
 Pushes the docker image to the repository.
 
 - **test**
 
 Runs local unit test with coverage check. The tested components are 
 mainly controllers. This gives coverage of the most of the components.
 The **DEBUG** env variable triggers on the show of the log on the screen.
  
## Database installation

The scripts provided in the sql directory allow to install the tables in t

 - local test Oracle XE environment running in docker
 - test BSCS environment T17
 - production BSCS environment
 
They can be installed with Makefile targets: 

 - make sql-xe
 - make sql-t17
 - make sql-prod

The objects like Account or Order are accessible via views returning last
production version. This is the feed for reports using the backend tables:

 - **SAP_ACCOUNT_P**
 - **SAP_ACC_SEGM_ORDER_NUMBERS_P**
