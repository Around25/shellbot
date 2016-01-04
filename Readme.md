Shellbot
========

Shellbot allows the user to connect to any server through ssh and setup a certain infrastructure configuration
by running various commands on the remote server or by copying data from the host through SCP.


Installation
------------

Download the appropriate release for your architecture and start using it.

Usage
-----

__Copy Command__

Using the copy command you can copy a file or directory from the local host to a server and vice versa.

In order to user the command you must first have a configuration file that defines the list of managed servers.

To copy from local host to a server named dev-1 use this command: `$> shellbot --config ./shellbot/devops.yaml copy ./hosts.txt dev:/etc/hosts`
To download from a server use this command: `$> shellbot --config ./shellbot/devops.yaml copy dev:/etc/hosts ./hosts.txt`

__Shell Command__
Connect to a particular server using ssh use this command: `$> shellbot --config ./shellbot/devops.yaml shell dev-1`

__Setup Command__

To execute all tasks for the "dev" environment use this command: `$> shellbot --config ./shellbot/devops.yaml setup dev`

__Check Command -- NOT YET IMPLEMENTED__

The check command allows you to see if all servers for an environment are in their correct state.

Contributing
------------

In order to install all dependencies for the project run `make deps`.
To compile and install the `shellbot` tool run `make`.
In order to execute the tests run `make test` or `make test-coverage` in case you want to see the code coverage as well.  