goNetViz
========

This is `goNetViz` and it is written in [golang](https://golang.org/).
`goNetViz` visualizes your network traffic, either by reading a file or
attaching to a network interface.

[![Build Status](https://travis-ci.org/florianl/goNetViz.svg?branch=master)](https://travis-ci.org/florianl/goNetViz) [![Coverage Status](https://coveralls.io/repos/github/florianl/goNetViz/badge.svg?branch=master)](https://coveralls.io/github/florianl/goNetViz?branch=master)


        $ ./goNetViz -help
          ./goNetViz [-bits ...] [-count ...] [-file ... | -interface ...] [-filter ...] [-list_interfaces] [-help] [-prefix ...] [-size ... | -timeslize ... | -terminal] [-version]
        -bits uint
                Number of bits per pixel. It must be divisible by three and smaller than 25 or 1.
                To get black/white results, choose 1 as input. (default 24)
        -count uint
                Number of packets to process.
                If argument is 0 the limit is removed. (default 25)
        -file string
                Choose a file for offline processing.
        -filter string
                Set a specific filter.
        -help
                Show this help.
        -interface string
                Choose an interface for online processing.
        -list_interfaces
                List available interfaces.
        -prefix string
                Prefix of the resulting image. (default "image")
        -size uint
                Number of packets per image. (default 25)
        -terminal
                Visualize output on terminal.
        -timeslize uint
                Number of microseconds per resulting image.
                So each pixel of the height of the resulting image represents one microsecond.
        -version
                Show version.
        -version
                Show version

Building
--------

        $ git clone git@github.com:florianl/goNetViz.git
          Cloning into 'goNetViz'...
          [...]
        $ cd goNetViz/
        $ export GOPATH=$HOME/go
        $ go get -u github.com/google/gopacket
          [...]
        $ go build
          [...]
        $ ./goNetViz
          [...]

Or you can get it directly via [golang](https://golang.org/):

        $ export GOPATH=$HOME/go
        $ go get -u github.com/florianl/goNetViz
          [...]
        $ $GOPATH/bin/goNetViz
          [...]

pcap-filter
----------

See: https://www.tcpdump.org/manpages/pcap-filter.7.html


Examples
--------

The images below are based on the very same IP traffic. The differences are
based on the number of bits per pixel. In the first image one bit is used
per pixel. Then, in the second image, 3 bits of the payload are used per pixel.
This is followed by 9 and 12 bits per pixel and finally 24 bits per pixel.

![1 Payloadbits per Pixel](img/ping1.png)

![3 Payloadbits per Pixel](img/ping3.png)

![9 Payloadbits per Pixel](img/ping9.png)

![12 Payloadbits per Pixel](img/ping12.png)

![24 Payloadbits per Pixel](img/ping24.png)


Run
---

./goNetViz -interface en0 -scale 20 -filter "port 53"
