## xk6-mllp

Simple MLLP sender for K6


## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

```bash
$ go install go.k6.io/xk6/cmd/xk6@latest
$ xk6 build --with github.com/gjergjsheldija/xk6-mllp=. 
$ ./k6 run --vus 60 --duration 1m test.js   
```

## Docker

```shell
docker run -i gjergjsheldija/xk6-mllp:latest --vus 60 --duration 1m run - < test.js
```

## Example

```javascript
import mllp from 'k6/x/mllp';

const client = new mllp.Client({
    host: '127.0.0.1',
    port: '5000'
});

export default function () {
    client.send('./sample.hl7');
}
```
