# Design

Implement simple web-server application for storing stores items catalog.

* Implement storage in memory as array. Implement storage as package.\
Storage should provide three methods:\
Add (one or more items),\
Search (search for produce by partial match of whole words),\
Fetch (get information about specific item).\
Delete (delete item).\
Storage package should be thread-safe.
* Implemet REST API package. Define REST API Request/Payload formats. Error formats. Use JSON.\
REST API should implement Add, Search, Fetch methods wrapping storage methods.
* Code should be testable. Cover code with unit tests.
* Add fixtures for predefined items in storage.
* Add build/deployment instructions. Implement Makefile for build/deploy commands. Implement Dockerfile.

## Storage Format


## API Format

