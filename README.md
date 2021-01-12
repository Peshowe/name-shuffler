# Name shuffler

A simple package that shuffles an array of provided names and sends out an email with those names in a random order and a random quote of the day.
Can be scheduled and used for the order in daily team meetings.

## Usage

You could either clone the repo and build the module yourself, 
or you could get and install it with

```
go get github.com/Peshowe/name-shuffler
``` 

Depending on your choise of getting the package you can either run the main file [shuffle.go](shuffle.go) or the installed name-shuffler module. 
It accepts a single argument passed via the command line:

`--yamlPath`
Path to the YAML file with the data for the script.

```
./name-shuffler [--yamlPath yamlFile]
```

### With Docker 

Alternatively, a Dockerfile is available in the repo, that takes the yaml file in the current directory and builds an image that can run the name shuffler script. 

You can use the following command to create an image called "name-shuffler":
```
docker build -t name-shuffler .
```

You can then run a container like this (--rm flag will remove the container once it's done):
```
docker run --rm name-shuffler
```


### Structure of YAML

An example YAML file is given in [email_details.yaml](email_details.yaml). 

## Example output

Contents of an example email that will be sent by the script: 

```
Today's daily order will be the following:

Petar 
Ivan 
Bob 
Alice 


Random quote of the day:
"Many a false step was made by standing still."
- Fortune Cookie

Regards,
The random generator

```