# Math Service

There are 6 application containers. Note that "Order of Operations" won't be dealt with. The system will take a string (`"5+2*10/5-6"`) and return a number (`8`).

## Services

### Gateway

This service fronts the system and exposes ability to do math to the world

### Tokenization

There is one service for tokenization. This returns an ordering of which services to call

### Operators

These are the four operator: +, -, *, /

## Building

This following will build 6 Dockerfiles

```
./build-services.sh
```

## Testing

[This link](https://www.w3schools.com/tags/ref_urlencode.asp) will help with URL encoding your requests.

```
curl "http://172.17.4.101:32111/compute?equation=22%2B99"
```
