# ascii character frequencies in english

This repo provides files ([json](https://raw.githubusercontent.com/piersy/ascii-char-frequency-english/main/ascii_freq.json),
[txt](https://raw.githubusercontent.com/piersy/ascii-char-frequency-english/main/ascii_freq.txt)) detailing the frequencies of ascii chararcters
ocurring in english text. The characters are represented as the decimal value
of the byte and the frequencies represented as 64 bit floating point values.

The files were derived from the [Reuters21578
corpus](http://www.daviddlewis.com/resources/testcollections/reuters21578/)
which is included along with the code used to derive the frequencies.

## Running the code

You will need golang version 1.13 or higher.

To build and run the code, run:
```
go run main.go
```

