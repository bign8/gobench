# GoBench Web
Web interface for GoBench running on Appengine

## Develop
```
goapp serve
```

## Deploy
```
go generate
```

## Notes

- `static/plotly.js` is not a full plotly distribution, but is the basic partial bundle.  See https://github.com/plotly/plotly.js/blob/master/dist/README.md#plotlyjs-basic for details.
- Appengine automagically GZIPs content for supportive clients
