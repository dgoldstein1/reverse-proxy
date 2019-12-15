# Custom Reverse Proxy

A fast and reusable reverse proxy as an alternative to serving deployments through proxies with complicated configurations like NGINX or Envoy. Especially useful for multiple service deployments on heroku, since the port changes and is dynamically configured with the `PORT ` environment variable.

## Configuration

All configuration is done through environment variables. First define the names for each service you want to proxy to:

```bash
export services="service1,service2" #delimited by commas
```

Then define the incoming back off of the proxy and outgoing URL to the remote server for each service you defined above.

```bash
export service1_incoming_path="/service1/"
export service1_outgoing_url="http://google.com"
export service2_incoming_path="/service2/"
export service2_outgoing_url="http://wikipedia.org"
```

Now set the port on which you want the proxy to be served on.

```bash
export PORT=8443
```

A sample configuration can be found in [the dev env file](./devEnv.sh).

Start the service and you should see:

```bash
2019/12/15 14:04:31 read in services [service1 service2]
2019/12/15 14:04:31 read in config ([]main.proxyConfig) (len=2 cap=2) {
 (main.proxyConfig) {
  incomingPath: (string) (len=10) "/service1/",
  outgoingURL: (*url.URL)(0xc000100000)(http://google.com),
  name: (string) (len=8) "service1"
 },
 (main.proxyConfig) {
  incomingPath: (string) (len=10) "/service2/",
  outgoingURL: (*url.URL)(0xc000100080)(http://wikipedia.org),
  name: (string) (len=8) "service2"
 }
}
2019/12/15 14:04:31 Serving on port :8443
..
```

Make a request to where you've deployed the service and you should see the following logs saying which endpoint and service is being hit:

```
2019/12/15 14:06:02 service2 -- /
2019/12/15 14:06:06 service1 -- /
2019/12/15 14:06:36 service2 -- /this/is/a/sample/route
```

Regex and more advanced route matching rules can be applied with the rules for the golang `http.HandleFunc` routing rules.

## Development




## Authors

* **David Goldstein** - [DavidCharlesGoldstein.com](http://www.davidcharlesgoldstein.com/?custom-reverse-proxy) - [Decipher Technology Studios](http://deciphernow.com/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
