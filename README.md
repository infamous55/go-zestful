# Go Zestful

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/infamous55/go-zestful)
![Lines of code](https://img.shields.io/tokei/lines/github/infamous55/go-zestful)
![GitHub](https://img.shields.io/github/license/infamous55/go-zestful?color=blue&logoColor=%20)

An in-memory key-value store for JSON data written in Golang.

## Key Features

- Support for both LRU and LFU eviction policies
- Authentication using a secret specified at initialization time and JSON Web Tokens (JWTs)
- Web server for interacting with the cached items

## Installation

Make sure you have Go installed on your system, then run the following command:

```
$ go install github.com/infamous55/go-zestful@v1.0.0-alpha
```

## Usage

After installing Go Zestful, you can check all the command-line options by executing the binary with no arguments, or by using the `-h` or `-help` flags:

```
$ go-zestful
Usage of go-zestful:
  -capacity uint
        set the capacity of the cache
  -default-ttl value
        set the default time-to-live
  -eviction-policy value
        set the eviction policy of the cache (LRU or LFU)
  -port value
        set the port number for the web server
  -secret string
        set the authorization secret
```

If you don't provide a required property, Go Zestful will also check for SCREAMING_SNAKE_CASE environment variables starting with `ZESTFUL_` (e.g. `ZESTFUL_DEFAULT_TTL`).

After initializing the cache, you can interact with it through the web server. The API supports the following routes:

- **POST** `/auth/token` for retrieving a JWT. The request body should contain the secret specified at initialization time.

Example request body:

```json
{
  "secret": "random_string"
}
```
- **POST** `/auth/refresh` for refreshing your JWT right before it expires.

- **GET** `/cache` for getting information about the cache.
- **DELETE** `/cache` for purging all the items in the cache.

- **GET** `/items/{key}` for getting the value of one item by its key.
- **POST** `/items` for creating an item. The request body should contain the key, an optional TTL (time-to-live), and the value.

Example request body:

```json
{
  "key": "example_key",
  "ttl": "1h", // optional
  "value": "example_value"
}
```
- **PUT** `/items/{key}` for updating an existing item's value.
- **DELETE** `/items/{key}` for deleting an item.

All requests that perform any kind of CRUD operations must provide a valid JWT in the Authorization header preceded by the string "Bearer ".

## To Do

- [ ] Limit how much memory can be used
- [ ] Add support for using a configuration file

## License

Distributed under the MIT License. See `LICENSE` for more information.