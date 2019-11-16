# go-web-dashboard

A simple dashboard showing system.

## What?

Simply run this somewhere, visit the URL it runs on from a device, name the
device, then you can control the page shown on the device from elsewhere.

The idea is mostly to make old tablets display some webpages easily. Use that as
a display around the house. This could work for displays in an office, but the
code isn't really hardened for such uses.

## Security

This does no auth -- wrap the service in something that does like one of the
many oauth proxies. Or run on a private network.

There's no XSRF protection, see note under "What?".

Some pages can't be embedded because they set `X-Frame-Options`. You can proxy
them to remove that (todo: some way of doing this automatically, e.g. proxy with
go proxy tools, inject `<base>` maybe...).

## Full usage

```
go run github.com/dgl/go-web-dashboard/cmd/go-web-dashboard
```

Then visit http://localhost:4000/show

Type in a name.

In a different tab visit http://localhost:4000/, hit send, put in URL.

That's it.

## How do I...

### ...refresh the page regularly?

Just put something like `curl http://localhost:4000/send?name=...&url=...` in cron.

todo: just do this in the JS code...
