# Software-Defined-Infrastructure

## DNS 

### Use dig to query A/CNAME/MX/NS records from various machines/domains of your choice. 
```
dig jakobwaibel.com
```
```
; <<>> DiG 9.16.8-Ubuntu <<>> jakobwaibel.com
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 28487
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 65494
;; QUESTION SECTION:
;jakobwaibel.com.		IN	A

;; ANSWER SECTION:
jakobwaibel.com.	300	IN	A	172.67.143.46
jakobwaibel.com.	300	IN	A	104.21.46.240

;; Query time: 23 msec
;; SERVER: 127.0.0.53#53(127.0.0.53)
;; WHEN: Mon Oct 18 20:22:30 CEST 2021
;; MSG SIZE  rcvd: 76
```

```
dig +noall +answer www.hdm-stuttgart.de
```
```
www.hdm-stuttgart.de.	3499	IN	A	141.62.1.59
www.hdm-stuttgart.de.	3499	IN	A	141.62.1.53
```

```
dig -x 141.62.1.53
```
```
; <<>> DiG 9.16.8-Ubuntu <<>> -x 141.62.1.53
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 27731
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 65494
;; QUESTION SECTION:
;53.1.62.141.in-addr.arpa.	IN	PTR

;; ANSWER SECTION:
53.1.62.141.in-addr.arpa. 3600	IN	PTR	iz-www-2.hdm-stuttgart.de.

;; Query time: 167 msec
;; SERVER: 127.0.0.53#53(127.0.0.53)
;; WHEN: Mon Oct 18 20:25:00 CEST 2021
;; MSG SIZE  rcvd: 92
```

### Setup BIND

Install BIND
```
sudo apt update
sudo apt install bind9 bind9utils bind9-doc
```

Check if the bind9 service is running 
```
systemctl status bind9
```

If not, start BIND
```
sudo systemctl start bind9
sudo systemctl enable named
```

Check if BIND (named) is listening on port 53
```
sudo ss -tlnp | grep named
```
```
# ss -tlnp | grep named
LISTEN 0      10                        141.62.75.101:53        0.0.0.0:*    users:(("named",pid=2976,fd=19))
LISTEN 0      10                            127.0.0.1:53        0.0.0.0:*    users:(("named",pid=2976,fd=16))
LISTEN 0      4096                          127.0.0.1:953       0.0.0.0:*    users:(("named",pid=2976,fd=25))
LISTEN 0      10                                [::1]:53           [::]:*    users:(("named",pid=2976,fd=22))
LISTEN 0      10     [fe80::d8f8:4cff:feae:b184]%eth0:53           [::]:*    users:(("named",pid=2976,fd=24))
LISTEN 0      4096                              [::1]:953          [::]:*    users:(("named",pid=2976,fd=26))root@sdi1a:~# ss -tlnp | grep named
LISTEN     0      10     141.62.75.101:53                       *:*                   users:(("named",pid=16759,fd=22))
LISTEN     0      10     127.0.0.1:53                       *:*                   users:(("named",pid=16759,fd=21))
LISTEN     0      128    127.0.0.1:953                      *:*                   users:(("named",pid=16759,fd=23))
```

Check the status of the bind name server
```
# rndc status
version: BIND 9.16.15-Debian (Stable Release) <id:4469e3e>
running on sdi1a: Linux x86_64 5.4.128-1-pve #1 SMP PVE 5.4.128-2 (Wed, 18 Aug 2021 16:20:02 +0200)
boot time: Thu, 21 Oct 2021 12:03:01 GMT
last configured: Thu, 21 Oct 2021 12:03:13 GMT
configuration file: /etc/bind/named.conf
CPUs found: 1
worker threads: 1
UDP listeners per interface: 1
number of zones: 102 (97 automatic)
debug level: 0
xfers running: 0
xfers deferred: 0
soa queries in progress: 0
query logging is OFF
recursive clients: 0/900/1000
tcp clients: 0/150
TCP high-water: 0
server is up and running
```

## Configure BIND

Edit the `/etc/bind/named.conf.options` file.

```
vim /etc/bind/named.conf.options
```
This file should have the following content:
```
options {
	directory "/var/cache/bind";

	// If there is a firewall between you and nameservers you want
	// to talk to, you may need to fix the firewall to allow multiple
	// ports to talk.  See http://www.kb.cert.org/vuls/id/800113

	// If your ISP provided one or more IP addresses for stable 
	// nameservers, you probably want to use them as forwarders.  
	// Uncomment the following block, and insert the addresses replacing 
	// the all-0's placeholder.

	forwarders {
		1.1.1.1; // Cloudflare
		8.8.8.8; // Google
	};

	//========================================================================
	// If BIND logs error messages about the root key being expired,
	// you will need to update your keys.  See https://www.isc.org/bind-keys
	//========================================================================
	dnssec-validation no;
	recursion yes;
	
};
```

Next, we need to add zones for our forward and reverse lookups. Achieve this by editing the `/etc/bind/named.conf.local` file.

```
vim /etc/bind/named.conf.local
```

This file should have the following content:
```
//
// Do any local configuration here
//

// Consider adding the 1918 zones here, if they are not used in your
// organization
//include "/etc/bind/zones.rfc1918";

zone "mi.hdm-stuttgart.de" {
	type master;
	file "/etc/bind/db.mi.hdm-stuttgart.de";
};

zone "75.62.141.in-addr.arpa" {
	type master;
	file "/etc/bind/db.141";
};
```

Create the template for our zone file by copying the template.

```
cp /etc/bind/db.empty /etc/bind/db.mi.hdm-stuttgart.de
```

```
vim /etc/bind/db.mi.hdm-stuttgart.de
```
The file should have the following content

```
; BIND reverse data file for empty rfc1918 zone
;
; DO NOT EDIT THIS FILE - it is used for multiple zones.
; Instead, copy it, edit named.conf, and use that copy.
;
$TTL	86400
; Start of Authority record defining the key characteristics of this zone
@	IN	SOA	ns1.mi.hdm-stuttgart.de. hostmaster.mi.hdm-stuttgart.de. (
			      1		; Serial
			 604800		; Refresh
			  86400		; Retry
			2419200		; Expire
			  86400 )	; Negative Cache TTL

@	IN	NS	ns1.mi.hdm-stuttgart.de.
@	IN	A	141.62.75.101
@	IN	MX	10	mx1.hdm-stuttgart.de.

; A records
www				IN	A	141.62.75.101
sdi1a.mi.hdm-stuttgart.de.	IN	A	141.62.75.101
sdi1b.mi.hdm-stuttgart.de.	IN	A	141.62.75.101
ns1.mi.hdm-stuttgart.de. 	IN	A	141.62.75.101

; CNAME records
www1-1	IN	CNAME	www	
www1-2	IN	CNAME	www	
info	IN	CNAME	www

```

Enable IPv4 in `/etc/default/named`

The content should be:
```
#
# run resolvconf?
RESOLVCONF=no

# startup options for the server
OPTIONS="-4 -u bind"
```

Now create the `/etc/bind/db.141` file and insert the following:

```
;
; BIND reverse data file for local loopback interface
;
$TTL    604800
@       IN      SOA     ns1.mi.hdm-stuttgart.de. hostmaster.mi.hdm-stuttgart.de. (
                              2         ; Serial
                         604800         ; Refresh
                          86400         ; Retry
                        2419200         ; Expire
                         604800 )       ; Negative Cache TTL
;
@	IN	NS	ns1.mi.hdm-stuttgart.de.
101	IN	PTR	sdi1a.mi.hdm-stuttgart.de.

```

Now restart the bind9 service

``` 
systemctl restart bind9
```

Now everything should work accordingly. Test your configurations with the following commands:

```
# dig @141.62.75.101 sdi1a.mi.hdm-stuttgart.de

; <<>> DiG 9.16.15-Debian <<>> @141.62.75.101 sdi1a.mi.hdm-stuttgart.de
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 15934
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: 38476f14e5071de20100000061752cde10f1ae98e7409184 (good)
;; QUESTION SECTION:
;sdi1a.mi.hdm-stuttgart.de.	IN	A

;; ANSWER SECTION:
sdi1a.mi.hdm-stuttgart.de. 86400 IN	A	141.62.75.101

;; Query time: 0 msec
;; SERVER: 141.62.75.101#53(141.62.75.101)
;; WHEN: Sun Oct 24 11:52:30 CEST 2021
;; MSG SIZE  rcvd: 98

```

```
 dig @141.62.75.101 ns1.mi.hdm-stuttgart.de

; <<>> DiG 9.16.15-Debian <<>> @141.62.75.101 ns1.mi.hdm-stuttgart.de
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 5720
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: 40f7e22511c8f42e0100000061752d2d56c4468f01c6c800 (good)
;; QUESTION SECTION:
;ns1.mi.hdm-stuttgart.de.	IN	A

;; ANSWER SECTION:
ns1.mi.hdm-stuttgart.de. 86400	IN	A	141.62.75.101

;; Query time: 0 msec
;; SERVER: 141.62.75.101#53(141.62.75.101)
;; WHEN: Sun Oct 24 11:53:49 CEST 2021
;; MSG SIZE  rcvd: 96

```

```
# dig @141.62.75.101 www1-1.mi.hdm-stuttgart.de

; <<>> DiG 9.16.15-Debian <<>> @141.62.75.101 www1-1.mi.hdm-stuttgart.de
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 21939
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: 0143132ee843744e0100000061752d45fd31a73abc08c9a3 (good)
;; QUESTION SECTION:
;www1-1.mi.hdm-stuttgart.de.	IN	A

;; ANSWER SECTION:
www1-1.mi.hdm-stuttgart.de. 86400 IN	CNAME	www.mi.hdm-stuttgart.de.
www.mi.hdm-stuttgart.de. 86400	IN	A	141.62.75.101

;; Query time: 0 msec
;; SERVER: 141.62.75.101#53(141.62.75.101)
;; WHEN: Sun Oct 24 11:54:13 CEST 2021
;; MSG SIZE  rcvd: 117

root@sdi1a:/etc/bind# 
```

```
dig @141.62.75.101 www1-2.mi.hdm-stuttgart.de

; <<>> DiG 9.16.15-Debian <<>> @141.62.75.101 www1-2.mi.hdm-stuttgart.de
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 18823
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: b4ce868197631a250100000061752d4ffa262ce084071761 (good)
;; QUESTION SECTION:
;www1-2.mi.hdm-stuttgart.de.	IN	A

;; ANSWER SECTION:
www1-2.mi.hdm-stuttgart.de. 86400 IN	CNAME	www.mi.hdm-stuttgart.de.
www.mi.hdm-stuttgart.de. 86400	IN	A	141.62.75.101

;; Query time: 0 msec
;; SERVER: 141.62.75.101#53(141.62.75.101)
;; WHEN: Sun Oct 24 11:54:23 CEST 2021
;; MSG SIZE  rcvd: 117

```

```
# dig @141.62.75.101 info.mi.hdm-stuttgart.de

; <<>> DiG 9.16.15-Debian <<>> @141.62.75.101 info.mi.hdm-stuttgart.de
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 5033
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: 294653251b1b58c20100000061752d6a07e56e6fc50e545c (good)
;; QUESTION SECTION:
;info.mi.hdm-stuttgart.de.	IN	A

;; ANSWER SECTION:
info.mi.hdm-stuttgart.de. 86400	IN	CNAME	www.mi.hdm-stuttgart.de.
www.mi.hdm-stuttgart.de. 86400	IN	A	141.62.75.101

;; Query time: 0 msec
;; SERVER: 141.62.75.101#53(141.62.75.101)
;; WHEN: Sun Oct 24 11:54:50 CEST 2021
;; MSG SIZE  rcvd: 115

```

```
# dig @141.62.75.101 mx1.hdm-stuttgart.de

; <<>> DiG 9.16.15-Debian <<>> @141.62.75.101 mx1.hdm-stuttgart.de
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 45334
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: 250912354598ef4f0100000061752db4dcd3a65394a4e05a (good)
;; QUESTION SECTION:
;mx1.hdm-stuttgart.de.		IN	A

;; ANSWER SECTION:
mx1.hdm-stuttgart.de.	2390	IN	A	141.62.1.22

;; Query time: 0 msec
;; SERVER: 141.62.75.101#53(141.62.75.101)
;; WHEN: Sun Oct 24 11:56:04 CEST 2021
;; MSG SIZE  rcvd: 93

```

```
# dig @141.62.75.101 www.google.com

; <<>> DiG 9.16.15-Debian <<>> @141.62.75.101 www.google.com
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 51453
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: e2c3965c0b2a669a0100000061752dc21bed21ef598dd7f6 (good)
;; QUESTION SECTION:
;www.google.com.			IN	A

;; ANSWER SECTION:
www.google.com.		137	IN	A	142.250.186.68

;; Query time: 11 msec
;; SERVER: 141.62.75.101#53(141.62.75.101)
;; WHEN: Sun Oct 24 11:56:18 CEST 2021
;; MSG SIZE  rcvd: 87

```

```
# dig @141.62.75.101 -x 141.62.75.101

; <<>> DiG 9.16.15-Debian <<>> @141.62.75.101 -x 141.62.75.101
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 36454
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1232
; COOKIE: b1a16000bbd7c0420100000061752dd7df9f14643255bb09 (good)
;; QUESTION SECTION:
;101.75.62.141.in-addr.arpa.	IN	PTR

;; ANSWER SECTION:
101.75.62.141.in-addr.arpa. 604800 IN	PTR	sdi1a.mi.hdm-stuttgart.de.

;; Query time: 0 msec
;; SERVER: 141.62.75.101#53(141.62.75.101)
;; WHEN: Sun Oct 24 11:56:39 CEST 2021
;; MSG SIZE  rcvd: 122

```

## LDAP

### Browsing an existing LDAP Server using `Apache Directory Studio`

#### Setup Apache Directory Studio to anonymously connect to ldap1.hdm-stuttgart.de using TLS.

To setup Apache Directory Studio to connect to the LDAP server, create a new connection with the following properties: 
![Connect to LDAP server](./static/ldap_connection.png)

If your connection was successful, your user interface should look somewhat like this: 
![Connection successful](./static/ldap_connection_success.png)

#### Use a filter like (uid-xy234) to find your personal entry beneath ou=userlist,dc=hdm-stuttgart,dc=de, Use the corresponding DN e.g. uid=xy234, ou=userlist,dc=hdm-stuttgart,dc=de to reconnect using password authentication. Then browse your own entry again. Can you spot any difference?

Therefore, you right-click on the userlist and apply a filter on the children. This should work like displayed in the following two pictures.

![Userlist](./static/ldap_filter_children_on.png)

![Filter Children](./static/ldap_filter_children.png)

As one can see in the picture below, the entry appears and can be examined further by clicking on it.

![Filter Children Result](./static/ldap_filter_children_result.png)

The other option to find ones own entry is to create a search. The content should look somewhat like in the following picture.

![Create Search](./static/ldap_create_search.png)

Then ones user should be found and can be examined. 

![Search Results](./static/ldap_search_results.png)

After following the same steps while being logged in, the user entry should look more like in the following picture. As one can see, some more personal information like the hash of the users password or his student-identification number is provided.

![Search Results while being logged in](./static/ldap_jw_information_logged_in.png)
### Browsing an existing LDAP Server using `ldapsearch`

#### Setup ldapsearch to anonymously connect to ldap1.hdm-stuttgart.de using TLS.

The following command connects to the LDAP server and displays the information below. The output is a shortened version, as there are too many user entries to display here.

```shell
$ ldapsearch -x -b "ou=userlist,dc=hdm-stuttgart,dc=de" -H ldap://ldap1.hdm-stuttgart.de
```

```shell
# gast39, userlist, hdm-stuttgart.de
dn: uid=gast39,ou=userlist,dc=hdm-stuttgart,dc=de
hdmCategory: 4
sn: fixme
loginShell: /bin/sh
uidNumber: 46139
gidNumber: 35102
uid: gast39
objectClass: inetOrgPerson
objectClass: posixAccount
objectClass: shadowAccount
objectClass: hdmAccount
objectClass: hdmSambaDomain
objectClass: eduPerson
cn:: dW5rbm93biA=
homeDirectory: /home/stud/XX/gast39
givenName: Gast
eduPersonAffiliation: faculty
eduPersonAffiliation: library-walk-in

# gast38, userlist, hdm-stuttgart.de
dn: uid=gast38,ou=userlist,dc=hdm-stuttgart,dc=de
hdmCategory: 4
sn: fixme
loginShell: /bin/sh
uidNumber: 46138
gidNumber: 35102
uid: gast38
objectClass: inetOrgPerson
objectClass: posixAccount
objectClass: shadowAccount
objectClass: hdmAccount
objectClass: hdmSambaDomain
objectClass: eduPerson
cn:: dW5rbm93biA=
homeDirectory: /home/stud/XX/gast38
givenName: Gast
eduPersonAffiliation: faculty
eduPersonAffiliation: library-walk-in
```

#### Use a filter like (uid-xy234) to find your personal entry beneath ou=userlist,dc=hdm-stuttgart,dc=de, Use the corresponding DN e.g. uid=xy234, ou=userlist,dc=hdm-stuttgart,dc=de to reconnect using password authentication. Then browse your own entry again. Can you spot any difference?	

```shell
$ ldapsearch -x -b "uid=jw163, ou=userlist,dc=hdm-stuttgart,dc=de" -H ldap://ldap1.hdm-stuttgart.de
# extended LDIF
#
# LDAPv3
# base <uid=jw163, ou=userlist,dc=hdm-stuttgart,dc=de> with scope subtree
# filter: (objectclass=*)
# requesting: ALL
#

# jw163, userlist, hdm-stuttgart.de
dn: uid=jw163,ou=userlist,dc=hdm-stuttgart,dc=de
displayName: Waibel Jakob Elias
employeeType: student
objectClass: hdmAccount
objectClass: hdmStudent
objectClass: inetOrgPerson
objectClass: posixAccount
objectClass: shadowAccount
objectClass: eduPerson
eduPersonAffiliation: member
eduPersonAffiliation: student
eduPersonAffiliation: library-walk-in
uid: jw163
mail: jw163@hdm-Stuttgart.de
uidNumber: 67828
cn: Waibel Jakob Elias
loginShell: /bin/sh
hdmCategory: 1
gidNumber: 100
givenName: Jakob Elias
homeDirectory: /home/stud/j/jw163
sn: Waibel

# search result
search: 2
result: 0 Success

# numResponses: 2
# numEntries: 1
```

After being logged in, we can spot some more information:

```shell
$ ldapsearch -x -D "uid=jw163, ou=userlist,dc=hdm-stuttgart,dc=de"  -W -H ldap://ldap1.hdm-stuttgart.de -b " uid=jw163, ou=userlist,dc=hdm-stuttgart,dc=de"  -s sub 'uid=jw163'
Enter LDAP Password: 
# extended LDIF
#
# LDAPv3
# base < uid=jw163, ou=userlist,dc=hdm-stuttgart,dc=de> with scope subtree
# filter: uid=jw163
# requesting: ALL
#

# jw163, userlist, hdm-stuttgart.de
dn: uid=jw163,ou=userlist,dc=hdm-stuttgart,dc=de
businessCategory: 1
businessCategory: {11112-3}
employeeType: student
postOfficeBox: 2G
objectClass: hdmAccount
objectClass: hdmStudent
objectClass: inetOrgPerson
objectClass: posixAccount
objectClass: shadowAccount
objectClass: eduPerson
eduPersonAffiliation: member
eduPersonAffiliation: student
eduPersonAffiliation: library-walk-in
uid: jw163
mail: jw163@hdm-Stuttgart.de
uidNumber: 67828
cn: Waibel Jakob Elias
loginShell: /bin/sh
hdmCategory: 1
gidNumber: 100
employeeNumber: CENSORED
givenName: Jakob Elias
homeDirectory: /home/stud/j/jw163
sn: Waibel
matrikelNr: CENSORED 
userPassword:: CENSORED
shadowLastChange: 18316
sambaNTPassword: CENSORED

# search result
search: 2
result: 0 Success

# numResponses: 2
# numEntries: 1
```

### Setup an OpenLDAP server

First of all, one might use the following commands to install some useful utilities.

```shell
# apt install dialog
```

```shell
# apt install slapd
```

In the following, the admin password is set to `password`. Please make sure you are using a more secure password than we are when configuring your LDAP server.	

With `dpkg-reconfigure slapd`, you can modify your base configuration. Following the dialog accordingly should lead to  a successful result.

![First dialog page](./static/ldap_1.png)

![Second dialog page](./static/ldap_2.png)

![Third dialog page](./static/ldap_3.png)

![Forth dialog page](./static/ldap_4.png)

![Fifth dialog page](./static/ldap_5.png)

![Sixth dialog page](./static/ldap_6.png)

After setting up slapd, we can use `ss -tlnp` to verify that a server is running.

```shell
# ss -tlnp
State                    Recv-Q                   Send-Q                                     Local Address:Port                                       Peer Address:Port                   Process                                             
LISTEN                   0                        10                                         141.62.75.101:53                                              0.0.0.0:*                       users:(("named",pid=448,fd=19))                    
LISTEN                   0                        10                                             127.0.0.1:53                                              0.0.0.0:*                       users:(("named",pid=448,fd=16))                    
LISTEN                   0                        128                                              0.0.0.0:22                                              0.0.0.0:*                       users:(("sshd",pid=127,fd=3))                      
LISTEN                   0                        4096                                           127.0.0.1:953                                             0.0.0.0:*                       users:(("named",pid=448,fd=21))                    
LISTEN                   0                        100                                            127.0.0.1:25                                              0.0.0.0:*                       users:(("master",pid=292,fd=13))                   
LISTEN                   0                        1024                                             0.0.0.0:389                                             0.0.0.0:*                       users:(("slapd",pid=7245,fd=8))                    
LISTEN                   0                        128                                                 [::]:22                                                 [::]:*                       users:(("sshd",pid=127,fd=4))                      
LISTEN                   0                        100                                                [::1]:25                                                 [::]:*                       users:(("master",pid=292,fd=14))                   
LISTEN                   0                        1024                                                [::]:389                                                [::]:*                       users:(("slapd",pid=7245,fd=9))
```

```shell
ldapsearch -Q -LLL -Y EXTERNAL -H ldapi:/// -b cn=config dn
dn: cn=config

dn: cn=module{0},cn=config

dn: cn=schema,cn=config

dn: cn={0}core,cn=schema,cn=config

dn: cn={1}cosine,cn=schema,cn=config

dn: cn={2}nis,cn=schema,cn=config

dn: cn={3}inetorgperson,cn=schema,cn=config

dn: olcDatabase={-1}frontend,cn=config

dn: olcDatabase={0}config,cn=config

dn: olcDatabase={1}mdb,cn=config
```

```shell
# ldapwhoami -x 
anonymous
```

As mentioned in the task, we need to rename the `dc` to `betrayer.com`. So we just do `dpkg-reconfigure slapd` another time using the required information.

![Rename dc to destroyer.com](./static/ldap_destroyer.png)

We can now use `ldapwhoami` with our configuration:

```shell
# ldapwhoami -x -D cn=admin,dc=betrayer,dc=com -W
Enter LDAP Password: 
dn:cn=admin,dc=betrayer,dc=com
```

We can connect to our server using Apache Directory Studio with the following configuration.

![Connect to LDAP server](./static/ldap_connect_openldap.png)

To authorize as an administrator, we just use our admin credentials.

![Connect as Administrator](./static/openldap_admin.png)

### Populating your DIT

After creating our LDAP tree, it looks like this: 

![LDAP tree](./static/organizational_structure.png)

Our export dump looks like this: 

```ldif
version: 1

dn: dc=betrayer,dc=com
objectClass: dcObject
objectClass: organization
objectClass: top
dc: betrayer
o: betrayer.com

dn: cn=admin,dc=betrayer,dc=com
objectClass: organizationalRole
objectClass: simpleSecurityObject
cn: admin
userPassword:: e1NTSEF9cEhFK0VQT0cyZ3lSeU9nanZGcXNXT2I1ekdzR2w5Q0Q=
description: LDAP administrator

dn: ou=departments,dc=betrayer,dc=com
objectClass: organizationalUnit
objectClass: top
ou: departments

dn: ou=software,ou=departments,dc=betrayer,dc=com
objectClass: organizationalUnit
objectClass: top
ou: software

dn: ou=financial,ou=departments,dc=betrayer,dc=com
objectClass: organizationalUnit
objectClass: top
ou: financial

dn: ou=devel,ou=software,ou=departments,dc=betrayer,dc=com
objectClass: organizationalUnit
objectClass: top
ou: devel

dn: ou=testing,ou=software,ou=departments,dc=betrayer,dc=com
objectClass: organizationalUnit
objectClass: top
ou: testing

dn: uid=bean,ou=devel,ou=software,ou=departments,dc=betrayer,dc=com
objectClass: inetOrgPerson
objectClass: organizationalPerson
objectClass: person
objectClass: top
cn: Audrey Bean
sn: Bean
givenName: Audrey
mail: bean@betrayer.com
uid: bean
userPassword:: e3NtZDV9YVhKL2JlVkF2TDRENk9pMFRLcDhjM3ovYTZQZzBXeHA=

dn: uid=smith,ou=devel,ou=software,ou=departments,dc=betrayer,dc=com
objectClass: inetOrgPerson
objectClass: organizationalPerson
objectClass: person
objectClass: top
cn: Jane Smith
sn: Smith
givenName: Jane
mail: smith@betrayer.com
uid: smith
userPassword:: e3NtZDV9YVhKL2JlVkF2TDRENk9pMFRLcDhjM3ovYTZQZzBXeHA=

dn: uid=waibel,ou=financial,ou=departments,dc=betrayer,dc=com
objectClass: inetOrgPerson
objectClass: organizationalPerson
objectClass: person
objectClass: top
cn: Jakob Waibel
sn: Waibel
givenName: Jakob
mail: waibel@betrayer.com
uid: waibel
userPassword:: e3NtZDV9YVhKL2JlVkF2TDRENk9pMFRLcDhjM3ovYTZQZzBXeHA=

dn: uid=simpson,ou=financial,ou=departments,dc=betrayer,dc=com
objectClass: inetOrgPerson
objectClass: organizationalPerson
objectClass: person
objectClass: top
cn: Homer Simpson
sn: Simpson
givenName: Homer
mail: simpson@betrayer.com
uid: simpson
userPassword:: e3NtZDV9YVhKL2JlVkF2TDRENk9pMFRLcDhjM3ovYTZQZzBXeHA=

dn: uid=pojtinger,ou=testing,ou=software,ou=departments,dc=betrayer,dc=com
objectClass: inetOrgPerson
objectClass: organizationalPerson
objectClass: person
objectClass: top
cn: Felix Pojtinger
sn: Pojtinger
givenName: Felix
mail: pojtinger@betrayer.com
uid: pojtinger
userPassword:: e3NtZDV9YVhKL2JlVkF2TDRENk9pMFRLcDhjM3ovYTZQZzBXeHA=

dn: uid=simpson,ou=testing,ou=software,ou=departments,dc=betrayer,dc=com
objectClass: inetOrgPerson
objectClass: organizationalPerson
objectClass: person
objectClass: top
cn: Maggie Simpson
sn: Simpson
givenName: Maggie
mail: simpson@betrayer.com
uid: simpson
userPassword:: e3NtZDV9YVhKL2JlVkF2TDRENk9pMFRLcDhjM3ovYTZQZzBXeHA=
```

### Testing a bind operation as non - admin user

First we set a password for `bean`. In our case, we set the password to `password`.

![Set password for bean](./static/bean_pw.png)

Afterwards we create a new connection with the following properties for the authentication section:

![Bean authentication](./static/bean.png)

After that, our new connection looks like this: 

![Bean view](./static/bean_auth.png)

### Extending an existing entry

To extend an existing entry and to add a `posixAccount`, one needs to right-click on the `objectClass` attribute and click on `New Value`.  

![posixAccount](./static/posix.png)

After that look for the `posixAccount` like displayed in the screenshot and `Add` it so it appears in the right column. After clicking `Next`, one needs to specify a valid `gidNumber`, `homeDirectory` and `uidNumber`. Now complete the configuration by clicking on `Finish`.

### Filter based search

#### All users with a uid attribute value starting with the letter “b”.

![Search properties](./static/users_start_with_b.png)

![Search results](./static/result_starts_with_b.png)


#### All entries either with either a defined uid attribute or a ou attribute starting with letter “d”.

![Search properties](./static/filter_defined_uid_ou_starting_with_d.png)

![Search results](./static/uid_or_starts_with_d.png)

#### All users entries within the whole DIT having a gidNumber value of 100.

![Search properties](./static/gid100.png)

![Search results](./static/gid100_result.png)
#### All user entries belonging to the billing department having a uid value greater than 1023.

![Search properties](./static/uid1023.png)

![Search results](./static/uid1023_result.png)
#### All user entries within the whole DIT having a commonName containing the substring “ei”.

![Search properties](./static/ei.png)

![Search results](./static/ei_result.png)
#### All user entries within the whole DIT belonging to gidNumber == 100 or having a uid value starting with letter “t”.

![Search properties](./static/100_t.png)

![Search results](./static/100_t_result.png)

#### Accessing LDAP data by a mail client

To use LDAP data with a mail client, it is necessary to add a new address book. It's configuration is displayed in the picture below:

![Address book configuration](./static/mail_configuration.png)

After configuring the address book, addresses are only available when searching. As shown in the following picture, we can now find the fictional characters we created when populating our tree. 

![Example of LDAP Address](./static/mail_example.png)

#### LDAP configuration

```bash
# ldapsearch -Y EXTERNAL -H ldapi:/// -b cn=config
SASL/EXTERNAL authentication started
SASL username: gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth
SASL SSF: 0
# extended LDIF
#
# LDAPv3
# base <cn=config> with scope subtree
# filter: (objectclass=*)
# requesting: ALL
#

# config
dn: cn=config
objectClass: olcGlobal
cn: config
olcArgsFile: /var/run/slapd/slapd.args
olcLogLevel: none
olcPidFile: /var/run/slapd/slapd.pid
olcToolThreads: 1

# module{0}, config
dn: cn=module{0},cn=config
objectClass: olcModuleList
cn: module{0}
olcModulePath: /usr/lib/ldap
olcModuleLoad: {0}back_mdb

# schema, config
dn: cn=schema,cn=config
objectClass: olcSchemaConfig
cn: schema
olcObjectIdentifier: OLcfg 1.3.6.1.4.1.4203.1.12.2
olcObjectIdentifier: OLcfgAt OLcfg:3
olcObjectIdentifier: OLcfgGlAt OLcfgAt:0
...
olcDatabase: {-1}frontend
olcAccess: {0}to * by dn.exact=gidNumber=0+uidNumber=0,cn=peercred,cn=external
 ,cn=auth manage by * break
olcAccess: {1}to dn.exact="" by * read
olcAccess: {2}to dn.base="cn=Subschema" by * read
olcSizeLimit: 500

# {0}config, config
dn: olcDatabase={0}config,cn=config
objectClass: olcDatabaseConfig
olcDatabase: {0}config
olcAccess: {0}to * by dn.exact=gidNumber=0+uidNumber=0,cn=peercred,cn=external
 ,cn=auth manage by * break
olcRootDN: cn=admin,cn=config

# {1}mdb, config
dn: olcDatabase={1}mdb,cn=config
objectClass: olcDatabaseConfig
objectClass: olcMdbConfig
olcDatabase: {1}mdb
olcDbDirectory: /var/lib/ldap
olcSuffix: dc=betrayer,dc=com
olcAccess: {0}to attrs=userPassword by self write by anonymous auth by * none
olcAccess: {1}to attrs=shadowLastChange by self write by * read
olcAccess: {2}to * by * read
olcLastMod: TRUE
olcRootDN: cn=admin,dc=betrayer,dc=com
olcRootPW: {SSHA}0mF0mX3FUNd0zZ0u3bKPQQA1ubycnHvh
olcDbCheckpoint: 512 30
olcDbIndex: objectClass eq
olcDbIndex: cn,uid eq
olcDbIndex: uidNumber,gidNumber eq
olcDbIndex: member,memberUid eq
olcDbMaxSize: 1073741824

# search result
search: 2
result: 0 Success

# numResponses: 11
# numEntries: 10
```

As mentioned in the task, we find the required information near the end of the output. In the first occurence of `oldRootDN`, where `olcRootDN: cn=admin,cn=config`, no passwort is set afterwards. In the second occurence `olcRootDN=cn=admin,dc=betrayer,dc=com`, a password is set with `oldRootPW`. This password not being specified implies that configuration access is limited to localhost.

To change these settings, we must create a LDIF file which sets the required properties. A password hash can be generated using a website like: https://projects.marsching.org/weave4j/util/genpassword.php. We chose to use the password `password` again. 

```bash
# cat add_olcRootPW.ldif 
dn: olcDatabase={0}config,cn=config
replace: olcRootPW
olcRootPW: {ssha}U1oLsngiAsXhFQrav5UzHxn2/9qPNcL+
```

```bash
# ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f ~/add_olcRootPW.ldif
modifying entry "olcDatabase={0}config,cn=config"

```

![Network parameters for configuration connection](./static/1.png)

![Authentication for configuration connection](./static/2.png)

![Browser Options for configuration connection](./static/3.png)

After connecting, we can set the `olcLogLevel`. There are multiple options to set the log level to. We decided to get the most information possible, which is achieved by using the `olcLogLevel` of `-1`. 

![Set Log Level](./static/set_loglevel.png)

Now we can look at the syslog with `cat /var/log/syslog`. 

![Displaying Syslog](./static/syslog.png)

TODO: Is it required to write to a logfile?

