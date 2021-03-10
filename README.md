
# 🍧 Sherbet
High performance image manipulation server

### Architecture
```
                         |-------------|     |--------|
 Ruby HTTP Parser    --> |             |     |        |
 thumbor HTTP Parser --> | Op Genrator |     | Engine |
 Parser 1            --> |             | --> |        |
 Parser 2            --> |             |     |        |
                         |-------------|     |--------|   
```
![Blank diagram - Page 10](https://user-images.githubusercontent.com/661693/110565676-d5257000-811c-11eb-9e34-4d6be8266cb4.png)

### Supported parsers

#### Ruby style query params
Resize endpoint

```
/resize

```
Query Params
 + **url**
 + **width** - absolute or relative to the image size(650 or 0.5)
 + **height** - absolute or relative to the image size

```
/crop

```
```
 /process
 
 ```

### Operations
+ resize
+ crop
+ resize

#### Filters
+ smart




