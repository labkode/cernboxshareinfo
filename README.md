# cernboxshareinfo
Tool to create a STAR model schema to perform MDX queries on CERNBox shares

To use this tool, the following software packages need to be pre-installed:
* [CERN phonebook application (yum repo for Scientific Linux CERN 6)](http://linuxsoft.cern.ch/cern/slc6X/i386nonpae/yum/os-nonpae/Packages/phonebook-1.9.2-1.slc6.noarch.rpm)

Before using the tool you need to have the MySQL data exported to a file, you can achieve that by running:

```
mysql -u <your_user> -h <your_host> --password='<your_password>' owncloud -N -B -e "select * from oc_share" > /tmp/shares.txt
```

Once you have the data exportes you can run the *cernboxshareinfo* tool: TODO
