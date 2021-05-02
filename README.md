# Webhooks Receiver for Sonatype Nexus IQ

Generic webhook receiver for [Nexus IQ Server - Webhooks](https://help.sonatype.com/iqserver/automating/iq-server-webhooks).

Binary `iq-webhook-receiver` listens for IQ Server POST messages at /callback then logs the event payload to an arbitrary file (default: `events.log`).

This utility was built with Cloudwatch Agents in mind. An example of the log event produced is as follows:
```
2021/05/03 06:47:22 7f4a6dde-5c68-4999-bcc0-a62f3fb8ae48 iq:policyManagement 687f3719b87232cf1c11b3ef7ea10c49218b6df1 HmacSHA1 {"owner":{"access":[ ...

```

### A. Configure IQ Server
1. Click the *System Preferences* icon in the IQ Server toolbar and then click *Webhooks*.
2. Click the  *Create Webhook*  button. The dialog opens:

![Webhooks Dialog](https://help.sonatype.com/iqserver/files/330210/330211/5/1590007757377/webhooks.png)
3. Enter the webhook URL. (default: http://localhost:3000/callback)
4. Leave the *webhook secret key* blank. __Note__: This assumes you are running this utility in the same host as the IQ Server.
5. Select one or more Event Types.
6. Click *Create*.

### B. Build the Binary
```
GOOS=${OS} GOARCH=${ARCH} go build -o iq-webhook-receiver
```
You can find all valid GOOS and GOARCH combinations [here](https://golang.org/doc/install/source#environment).

### C. Running & Usage
```
$ ./iq-webhook-receiver -h
Usage of ./iq-webhook-receiver:
  -path string
        log file path (default "./events.log")
  -port int
        listen port (default 3000)
```