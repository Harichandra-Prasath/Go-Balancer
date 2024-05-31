# GO BALANCER

This is a load Balancer built with go that can balance the requests across Backends.

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Future Work](#future-work)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/Harichandra-Prasath/Go-Balancer.git
cd Go-Balancer
```

2. Create config.json 
[Configuration](#configuration)

3. Run the application

```bash 
make run
```

## Configuration

Below is the template for config.json

```bash
{
    "Port":,
    "Backends":["","",...],
    "ALGO":"",
    "STATIC_ROOT":"",
    "MEDIA_ROOT":""
}
```
1. **Port**  
Port for the Go-Balancer
2. **Backends**  
List of Url strings of the backends  
3. **ALGO**  
Algorithm used to balance the requests  
4. **STATIC_ROOT** 
Path string of the directory used to server Static Content  
5. **MEDIA_ROOT**  
Path string of the directory used to server Media Content  

  
Currently Supported Algorithms
```bash
"RR" : Round Robin
"LC" : Least Connections
```
Path **MUST** be Absolute from the root of the System  

## Future-Work

1. Including additional Balancing Algorithms like Spill Over , etc..  
2. More Configurable from the file including Health check Period, etc..