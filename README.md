# Name shuffler

A simple package that shuffles an array of provided names and sends out an email with those names in a random order and a random quote of the day.
Can be scheduled and used for the order in daily team meetings.

## Usage

The main file is [shuffle.go](shuffle.go). It accepts a single argument passed via the command line:

`--yamlFile`
Path to the YAML file with the data for the script.

```
./shuffle.go [--yamlFile yamlFile]
```

### Structure of YAML

An example YAML file is given in [email_details.yaml](email_details.yaml). 
