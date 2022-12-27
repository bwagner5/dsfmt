# dsfmt

dsfmt (Dist Stats Formatter) is a simple CLI utility to format disk stats (aka cat /proc/diskstats). 

## Why dsfmt?

Do you love the output of `/proc/diskstats`? 

```
cat /proc/diskstats
 259       0 nvme0n1 9678 7 1010217 42256 5163 1567 1052815 12339 0 21780 37572 0 0 0 0
 259       1 nvme0n1p1 9540 7 1005313 42123 5161 1559 1052735 12337 0 21692 37520 0 0 0 0
 259       2 nvme0n1p128 45 0 360 24 0 0 0 0 0 52 0 0 0 0 0
```

Do you have the man page memorized? 

Do you love counting columns to figure out what a number might be telling you?

Do you love miscounting columns and realizing that everything is fine, or everything is on fire? 

Maybe you love eye tests where you try to read large numbers with no commas?  

Or you just love doing millisecond conversions.

I don't like any of these things, so I made a tool that makes looking at disk stats a lot easier.


## Usage:

```
> cat /proc/diskstats | dsfmt
       DEVICE       | READS / MERGED | SECTORS READ | READ TIME | WRITES / MERGED | SECTORS WRITTEN | WRITE TIME | IOS NOW | IOS TIME | WEIGHTED IOS TIME | DISCARDS / MERGED | SECTORS DISCARDED | DISCARD TIME
--------------------+----------------+--------------+-----------+-----------------+-----------------+------------+---------+----------+-------------------+-------------------+-------------------+---------------
  nvme0n1 259/0     | 9677 / 7       |      1010113 | 42.255s   | 5050 / 1551     |         1051069 | 12.233s    |       0 | 21.708s  | 37.572s           | 0 / 0             |                 0 | 0s
+-------------------+----------------+--------------+-----------+-----------------+-----------------+------------+---------+----------+-------------------+-------------------+-------------------+--------------+
  nvme0n1p1 259/1   | 9539 / 7       |      1005209 | 42.122s   | 5048 / 1543     |         1050989 | 12.23s     |       0 | 21.62s   | 37.52s            | 0 / 0             |                 0 | 0s
+-------------------+----------------+--------------+-----------+-----------------+-----------------+------------+---------+----------+-------------------+-------------------+-------------------+--------------+
  nvme0n1p128 259/2 | 45 / 0         |          360 | 24ms      | 0 / 0           |               0 | 0s         |       0 | 52ms     | 0s                | 0 / 0             |                 0 | 0s
--------------------+----------------+--------------+-----------+-----------------+-----------------+------------+---------+----------+-------------------+-------------------+-------------------+---------------
```

Or use a shorter format which excludes Discard and Flush stats:

```
> cat /proc/diskstats | dsfmt --short
       DEVICE       | READS / MERGED | SECTORS READ | READ TIME | WRITES / MERGED | SECTORS WRITTEN | WRITE TIME | IOS NOW | IOS TIME | WEIGHTED IOS TIME
--------------------+----------------+--------------+-----------+-----------------+-----------------+------------+---------+----------+--------------------
  nvme0n1 259/0     | 9677 / 7       |      1010113 | 42.255s   | 5117 / 1567     |         1052261 | 12.302s    |       0 | 21.732s  | 37.572s
+-------------------+----------------+--------------+-----------+-----------------+-----------------+------------+---------+----------+-------------------+
  nvme0n1p1 259/1   | 9539 / 7       |      1005209 | 42.122s   | 5115 / 1559     |         1052181 | 12.3s      |       0 | 21.644s  | 37.52s
+-------------------+----------------+--------------+-----------+-----------------+-----------------+------------+---------+----------+-------------------+
  nvme0n1p128 259/2 | 45 / 0         |          360 | 24ms      | 0 / 0           |               0 | 0s         |       0 | 52ms     | 0s
--------------------+----------------+--------------+-----------+-----------------+-----------------+------------+---------+----------+--------------------
```

## Installation

Packages, binaries, and archives are published for all major platforms (Mac amd64/arm64 & Linux amd64/arm64):

Debian / Ubuntu:

```
[[ `uname -m` == "aarch64" ]] && ARCH="arm64" || ARCH="amd64"
OS=`uname | tr '[:upper:]' '[:lower:]'`
wget https://github.com/bwagner5/dsfmt/releases/download/v0.0.7/dsfmt_0.0.7_${OS}_${ARCH}.deb
dpkg --install dsfmt_0.0.7_linux_amd64.deb
cat /proc/diskstats | dsfmt
```

RedHat:

```
[[ `uname -m` == "aarch64" ]] && ARCH="arm64" || ARCH="amd64"
OS=`uname | tr '[:upper:]' '[:lower:]'`
rpm -i https://github.com/bwagner5/dsfmt/releases/download/v0.0.7/dsfmt_0.0.7_${OS}_${ARCH}.rpm
```

Download Binary Directly:

```
[[ `uname -m` == "aarch64" ]] && ARCH="arm64" || ARCH="amd64"
OS=`uname | tr '[:upper:]' '[:lower:]'`
wget -qO- https://github.com/bwagner5/dsfmt/releases/download/v0.0.7/dsfmt_0.0.7_${OS}_${ARCH}.tar.gz | tar xvz
chmod +x dsfmt
```


