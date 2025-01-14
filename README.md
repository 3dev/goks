
A Keystore files system for Go


# Overview

goks (go keystore) is a simple keystore file library and command line tool.

It is suitable for use in any situation where you would consider having a keystore
that stores content securely using a password. The file is simple a key-value
store but its content is encrypted using the provided password. The present 
implementation can store up to 1024 key-value pair securely in the keystore file. 
Each key can have a maximum of 32 bytes and the value can be of any length.


## goks Features

* A simple secure key-value keystore
* Up to 1024 key-values
* Simple and fast access to content
* It can be used to store application config securely (replacing tradition key-value config files which may contain sensitive data)
* Auto file extension (.goks)


# Using goks

goks is easy to use and easier to adopt.

## Step 1: Install goks

First use go get to install the latest version of the library.

    $ go get github.com/3dev/goks

Next include goks in your application.
```go
import "github.com/3dev/goks"
```

## Step 2: Use it like you would any other package

Throughout your application use any function and method like you normally
would.

```go
//Create keystore
ks, err := goks.New("/tmp/appKeystore.goks","securepassword")
if err != nil {
	panic(err)
}

//use ks by calling its methods
keys,_ := ks.keys()

//Opening an existing keystore
ks, err := goks.Open("/tmp/appKeystore.goks", "securepassword")
if err != nil {
	panic(err)
}

data, err := ks.get("myKey1")
```


The above sample is self explanatory. A keystore is either created ```goks.New()``` or opened ```goks.Open()```
and its content can be accessed, modified or added to using the functions listed below


## List of available functions

```go
func New(filename string, passkey string) (*KeyStore, error)
func Open(filename string, passkey string) (*KeyStore, error)
func (ks *KeyStore) Close() error
func (ks *KeyStore) Count() int
func (ks *KeyStore) Put(key string, data []byte) error
func (ks *KeyStore) Delete(key string) error
func (ks *KeyStore) Get(key string) ([]byte, error)
func (ks *KeyStore) KeyInfo(key string) (file.TableOfContent, error)
func (ks *KeyStore) Compact() error
```

## Command line tool
There is a provided command line tool that aids in creating and viewing the keystore file's content.
```commandline
goks create --file sample1.goks --pass Z2l2ZW1lYWp3dGdvb2RzZWNyZXRrZX
```
The above ```create``` command will create a new keystore file named ```sample1.goks``` protected with the password ```Z2l2ZW1lYWp3dGdvb2RzZWNyZXRrZX```

```aiignore
goks stats --file sample1.goks --pass Z2l2ZW1lYWp3dGdvb2RzZWNyZXRrZX 
```
The above ```stats``` command will reveal the statistics of the file like below
```aiignore
go keystore file:	'sample1.goks'
number of items:	 2
first key:		 "systemData"
```
To get the stats of a specific key you add the ```--key key_1``` argument to the call like so
```aiignore
goks stats --file sample1.goks --pass Z2l2ZW1lYWp3dGdvb2RzZWNyZXRrZX --key key_1 
```
And the output will be like

```aiignore
key info:
  index available:	false
  key name:		key_1
  data length:		36 bytes 	[0.04 KB]
  allocated space:	36 bytes 	[0.04 KB]
  file position:	46084 bytes 	[45.00 KB]
```
```aiignore
goks hex --file sample1.goks --pass Z2l2ZW1lYWp3dGdvb2RzZWNyZXRrZX --key key_1
```
The above command will give a hexdump of the provided key like
```aiignore
key info:
  index available:	false
  key name:		key_1
  data length:		112 bytes 	[0.11 KB]
  allocated space:	112 bytes 	[0.11 KB]
  file position:	46084 bytes 	[45.00 KB]

data:
00000000  44 41 54 41 5f 31 20 44  41 54 41 5f 31 20 44 41  |DATA_1 DATA_1 DA|
00000010  54 41 5f 31 20 44 41 54  41 5f 31 20 44 41 54 41  |TA_1 DATA_1 DATA|
00000020  5f 31 20 44 41 54 41 5f  31 20 44 41 54 41 5f 31  |_1 DATA_1 DATA_1|
00000030  20 44 41 54 41 5f 31 20  44 41 54 41 5f 31 20 44  | DATA_1 DATA_1 D|
00000040  41 54 41 5f 31 20 44 41  54 41 5f 31 20 44 41 54  |ATA_1 DATA_1 DAT|
00000050  41 5f 31 20 44 41 54 41  5f 31 20 44 41 54 41 5f  |A_1 DATA_1 DATA_|
00000060  31 20 44 41 54 41 5f 31  20 44 41 54 41 5f 31 20  |1 DATA_1 DATA_1 |
```

## License

goKeyStore is released under the MIT license.