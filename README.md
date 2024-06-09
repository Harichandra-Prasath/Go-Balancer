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
    "MEDIA_ROOT":"",
    "Poll_Period":
}
```
1. **Port : INT**  
Port for the Go-Balancer
2. **Backends : [ ]String**  
List of Url strings of the backends  
3. **ALGO : String**  
Algorithm used to balance the requests  
4. **STATIC_ROOT : String**    
Path string of the directory used to server Static Content  
5. **MEDIA_ROOT : String**  
Path string of the directory used to server Media Content   
6. **Poll_Period : UINT**  
Time in seconds that should be used for Passive Health checking of the backends  

**It is must for the config.json to have all the fields, exlusion of any field will result in failure**  
  
Currently Supported Algorithms
```bash
"RR" : Round Robin
"LC" : Least Connections
"RANDOM" : Randomly choosen backend
```
Path **MUST** be Absolute from the root of the System  

## Future-Work

1. Including additional Balancing Algorithms like Spill Over , etc..  