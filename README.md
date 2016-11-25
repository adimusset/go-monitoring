# go-monitoring

A Go command tool to monitor an actively written-to w3c-formatted HTTP access log


You can start it like this
```
$ go build
$ ./go-monitoring [threshold] [file path]
```

If you want to run the tests
```
$ glide up
$ go test -v
```

And to see it working you can create a fake log file
```
$ cd test
$ go build
$ ./test [file path]
```

You have 2 functionnalities:
* Reporting every 10 seconds about the traffic: most visited section and most made request
* Alerting when the threshold is gone through (up and down)

Feel free to edit the code and add routines to increase the performances, you can also add features as the design is very simple
