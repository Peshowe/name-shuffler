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

`--yamlFile`
Path to the YAML file with the data for the script.

```
./name-shuffler [--yamlFile yamlFile]
```

### Structure of YAML

An example YAML file is given in [email_details.yaml](email_details.yaml). 
