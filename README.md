# Flagmatic5000
🇨🇾 🇩🇪 🇬🇷 🇬🇧 🇺🇸 🇪🇸 🇳🇴 Tool to map city locations to country flags

This tool takes in a yaml list of people and uses the city attribute to determine the country they reside in.

It can find the country by checking the city against the csv of cities in the world. If the city is not in the csv or
there is a problem with the way the city is spelt in the csv, the country wont be found.

It then uses the ISO ALPHA-2 to grab the flag and appends it on the image.

Do what you wish with the code.

## Building

Make sure you have go installed.

**Install dependencies**
```sh
go get github.com/gocarina/gocsv
go get gopkg.in/yaml.v2
```
**Build it**

```sh
go build applyflags.go
```

**Run it**

```sh
./applyflags experts.yml
```
