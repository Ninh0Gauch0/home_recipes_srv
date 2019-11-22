# Logical Folder SpringBoot services Changelog

## Version 0.5.0-beta

Initial Home Recipes Service version:

* First approach cli implementation. Terminal signals handling.

* DTO definition.

* Server orchestration with channels and go routines.

* Server endpoints definition using mux library.

* Worker implementation with dummy responses.

* Error handling.

## Version 1.0.0

First stable version.

* All CRUD operations for recipes and ingredients are implemented.
* Mux library controls the http requests.
* Urface cli gives us a configuration cli tool for our application.
* We use channeling to orchestate our application.