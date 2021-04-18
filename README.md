# Containermon
Simple Docker container monitoring tool and save metrics to file

### Install with go:
```
go get -u github.com/aaitaazizi/containermon
```

### Container Name or ID is the only required argument
```
$ containermon --container jovial_euler
{"ts":"2020-05-04T18:00:51Z","timeElapsed":13.19,"cpuTimeElapsed":0.67,"percentCPUSinceStart":5.07,"percentCPUThisInterval":5.07}
{"ts":"2020-05-04T18:01:00Z","timeElapsed":22.88,"cpuTimeElapsed":0.91,"percentCPUSinceStart":3.97,"percentCPUThisInterval":2.48}
{"ts":"2020-05-04T18:01:10Z","timeElapsed":32.24,"cpuTimeElapsed":1.18,"percentCPUSinceStart":3.67,"percentCPUThisInterval":2.93}
```

### You can also output to csv and change the collection interval
```
$ containermon --container jovial_euler --interval 5 --output-format csv
ts,timeElapsed,cpuTimeElapsed,percentCPUSinceStart,percentCPUThisInterval
2020-05-04T18:21:41Z,7.80,0.26,3.39,3.39
2020-05-04T18:21:46Z,12.75,0.68,5.31,8.34
2020-05-04T18:21:50Z,17.20,0.80,4.68,2.85
```
