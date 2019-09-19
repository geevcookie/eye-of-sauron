# Eye of Sauron

This is a demo Go app that shows how you can collect system metrics. The metrics
can be viewed via a UI or simple API.

## Building

1. Clone this repo
2. Run `make build`

## Running

To run the app start the container with:

```bash
docker run -p 8080:8080 eye-of-sauron
```

By default the application will report on container stats. To get host stats
start the app with the following:

```bash
docker run -p 8080:8080 -v /proc:/host/proc:ro -e HOST_PROC=/host/proc eye-of-sauron 
```

## Usage

Once the application is running you can access the UI at `http://localhost:8080/`

The API docs can be accessed at `http://localhost:8080/swagger/index.html`

### Metrics

You can change the frequency metrics are collected by editing `eye-of-sauron.yml`.

You will need to rebuild after config changes.