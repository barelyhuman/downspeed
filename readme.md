# downspeed 
> fast.com speedtest in your terminal

It's more for getting an approx download speed than the 
accurate one and is done by using fast.com as a file source 

## Install 
Head over to the [releases](https://github.com/barelyhuman/downspeed/releases) page and download 
one for your system or compile one for your system.

## From source 

### Prerequisites
1. Make sure you have golang v1.18 installed 

### Compile
1. Clone 

```sh
git clone https://github.com/barelyhuman/downspeed
```

2. change into the directory, make sure the deps are installed locally, build the codebase

```sh
cd downspeed
go mod tidy
go build 
```

## Usage 

```
$ downspeed
```

Yeah, that's it, what more do you want?!

