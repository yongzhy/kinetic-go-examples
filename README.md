# Kinetic Go Examples

## Overview 

This repository contains a number of examples showing how to use Kinetic Go client in different way. 

## Instructions

kinetic-go client library is needed to build the examples. 

To install kinetic-go library : `go get github.com/Kinetic/kinetic-go`

Clone this examples repository : `git clone https://github.com/yongzhy/kinetic-go-examples`

Build examples : `cd kinetic-go-examples && ./buildall.sh`

Executable binary files will be in `bin` sub folder.

## Execute Examples 

All examples hard coded Kinetic device IP to be `127.0.0.1`, 
which assume [Kinetic Simulator](https://github.com/Kinetic/kinetic-java) runs on the same host that build those examples. 

To run example with actual Kinetic devices, please change the IP to your Kinetic device IP.

## License 

This project is licensed under Mozilla Public License, v. 2.0
* [Original](LICENSE) version