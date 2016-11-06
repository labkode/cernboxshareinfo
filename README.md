# cernboxshareinfo
Tool to create a STAR model schema to perform MDX queries on CERNBox shares.

## Requirements

To use this tool, the following software packages need to be pre-installed:

* [CERN phonebook application (yum repo for Scientific Linux CERN 6)](http://linuxsoft.cern.ch/cern/slc6X/i386nonpae/yum/os-nonpae/Packages/phonebook-1.9.2-1.slc6.noarch.rpm)

Before using the tool you need to have the MySQL data exported to a file, you can achieve that by running:

```
mysql -u <your_user> -h <your_host> --password='<your_password>' owncloud -N -B -e "select * from oc_share" > /tmp/shares.txt
```

You will also need a database with a star model to be fulfilled.
In this repo, you can find the [cernboxshareinfoschema.sql](./cernboxshareinfoschema.sql) to create the necessary structure:

```
mysql -u <your_user> -h <your_host> -p<your_password> < cernboxshareinfoschema.sql
```

## Usage

Once you have the data exported you can run the *cernboxshareinfo* tool.
For filling the star model you need to use **fillstarmodel=true** with your MySQL configuration.

```
$ ./cernboxshareinfo -h
Usage of ./cernboxshareinfo:
  -fillstarmodel
    	fill the existing sql star model on the configured mysql database
  -mysqldatabase string
    	mysql database (default "mydb")
  -mysqlhost string
    	mysql hostname (default "localhost")
  -mysqlpassword string
    	mysql password (default "admin")
  -mysqlport int
    	mysql port (default 3306)
  -mysqlusername string
    	mysql username (default "admin")
  -outputasjson
    	if enabled outputs the data in JSON format instead of proto format
  -phonebookbinary string
    	the phonebook binary (default "/usr/bin/phonebook")
  -phonebookdb string
    	phonebook cache db file (default "phonebook.db")
  -showdata
    	if enabled outputs the correlated data (one line per record) (default true)
  -sqlfile string
    	file containing the mysql oc_shares dump (default "shares.sql")
  -testcases int
    	number of test cases (default -1)
```
## Example

```
./cernboxshareinfo -mysqlusername=root -mysqlpassword=labkode -fillstarmodel=true 
```