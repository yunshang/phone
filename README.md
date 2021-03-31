# phone

Parsing, validating and creating phone numbers

## Description

Golanag library for phone number parsing, validation, and formatting.

## Install

    $ go get -u github/yunshang/phone

And then `go mod download` from your command line.

### Automatic country and area code detection

Phone does its best to automatically detect the country and area code while parsing. To do this, phone uses data stored in `data/phone/countries.yml`.

Each country code can have a regular expression named `area_code` that describes what the area code for that particular country looks like.

If an `area_code` regular expression isn't specified, a default value which is considered correct for the US will be used.

If your country has phone numbers longer that 8 digits - exluding country and area code - you can specify that within the country's configuration in `data/phone/countries.yml`

### Validating

Validating is very relaxed, basically it strips out everything that's not a number or '+' character:

```go
Phoner.valid("blabla 091/512-5486 blabla")
```

### Formatting

Formating is done via the `#format` method. The method accepts a `Symbol` or a `String`.

When given a string, it interpolates the string with the following fields:

* %c - country_code (385)
* %a - area_code (91)
* %A - area_code with leading zero (091)
* %n - number (5125486)
* %f - first @@n1_length characters of number (configured through country n1_length), default is 3 (512)
* %l - last characters of number (5486)
* %x - the extension number

```ruby
pn = Phoner.Parse("+385915125486")
pn.ToS # => "+385915125486"
pn.format("%A/%f-%l") # => "091/512-5486"
pn.format("+ %c (%a) %n") # => "+ 385 (91) 5125486"
```

When given a symbol it is used as a lookup for the format in the <tt>Phoner::Phone.named_formats</tt> hash.

```ruby
pn.format("europe") # => "+385 (0) 91 512 5486"
pn.format("us") # => "(234) 123-4567"
pn.format("default_with_extension") # => "+3851234567x143"
```

### Finding countries by their isocode

If you don't have the country code, but you know from other sources what country a phone is from, you can retrieve the country using the country isocode (such as 'de', 'es', 'us', ...). Remember to call `Phoner.load` before using this lookup.

```go
if country = Phoner.FindByCountryIsocode(user_country_isocode)
  phone_number = Phoner.parse(user_input, :country_code => country.country_code)
end
```

## Examples

```golang
import "github/yunshang/phone"
```

### Initializing

Initialize a new phone object with the number, area code, country code and extension number:

```go
args := []string{"5125486", "91", "385"}
Phoner.new(args)
```

```go
args := []string{"5125486", "91", "385", "143"}
Phoner.new(args)
```

### Parsing

Create a new phone object by parsing from a string. Phoner::Phone does it's best to detect the country and area codes:

```go
Phoner.Parse("+385915125486")
Phoner.Parse("00385915125486")
```

If the country or area code isn't given in the string, you must set it, otherwise it doesn't work:

```go
args := []string{"091/512-5486", "385"}
Phoner::Phone.parse(args)
```

If you feel that it's tedious, set the default country code once:

```go
Phoner.SetDefaultCountryCode = "385"
args := []string{"091/512-5486"}
Phoner.Parse(args)
```

Same goes for the area code:

```go
args := []string{"5125486", "91"}
Phoner.Parse(args)
```
or

```go
Phoner.SetDefaultCountryCode = "385"
Phoner.SetDefaultAreaCode = "47"
args := []string{"451-588"}

Phoner.Parse(args)
```

## Adding and maintaining countries

From time to time, the specifics about your countries information may change. You can add or update your countries configuration by editing `data/phone/countries.yml`

The following are the available attributes for configuration:

* `country_code`: Required. A string representing your country's international dialling code. e.g. "123"
* `national_dialing_prefix`: Required. A string representing your default dialling prefix for national calls. e.g. "0"
* `char_3_code`: Required. A string representing a country's ISO code. e.g. "US"
* `name`: Required. The name of the country. e.g. "Denmark"
* `international_dialing_prefix`: Required. The dialling prefix a country typically uses when making international calls. e.g. "0"
* `area_code`: Optional. A regular expression detailing valid area codes. Default: "\d{3}" i.e. any 3 digits.
* `max_num_length`: Optional. The maximum length of a phone number after country and area codes have been removed. Default: 8
