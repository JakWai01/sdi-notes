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

	// forwarders {
	// 	0.0.0.0;
	// };

	//========================================================================
	// If BIND logs error messages about the root key being expired,
	// you will need to update your keys.  See https://www.isc.org/bind-keys
	//========================================================================
	dnssec-validation auto;

	listen-on-v6 { any; };

	// hide version number from clients for security reasons.
	version "not currently available";

	// disable recursion on authoritative DNS server.
	recursion no;

	// enable the query log
	querylog yes;

	// disallow zone transfer
	allow-transfer { none; };
};
```

Next, we need to add a zone for our domain name. Achieve this by editing the `/etc/bind/named.conf.local` file.

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
	allow-query { any; };
	allow-transfer { 141.62.75.101; };
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

; Name servers for this domain
	IN	NS	ns1.mi.hdm-stuttgart.de.

; A records
ns1	IN	A	141.62.75.101
sdi1a	IN	A	141.62.75.101
www	IN	A	141.62.75.101
@	IN	A	141.62.75.101

; CNAME records
www1-1	IN	CNAME	141.62.75.101	
www1-2	IN	CNAME	141.62.75.101	
info	IN	CNAME	141.62.75.101
```

Enable IPv4 in `/etc/default/bind9`

The content should be:
```
#
# run resolvconf?
RESOLVCONF=no

# startup options for the server
OPTIONS="-4 -u bind"
```

Now, add the reverse zone. Therefore, edit the `/etc/bind/named.conf.local` file so it looks like this:

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
	allow-query { any; };
	allow-transfer { 141.62.75.101; };
};

zone "101.75.62.in-addr.arpa" {
	type master;
	notify no;
	file "/etc/bind/db.141";
};
```

Now create the `/etc/bind/db.141` file and insert the following:

```

```