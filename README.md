# A hex dump utility #

Similar to the Unix 'od' utility, but focused on hexadecimal
output.

## Parameters ##

Double minus (--) versions of parameters are also accepted.

-edgeMark string

    single character at edges of the right hand character block. (default "|")

-goDump

    if true, use standard go encoding/hex/Dump.

-h	print usage message.

-hexUpper

    if true, print upper case hex.

-inFile string

    input file name.  Argument 0 may also be used.

-inString string

    input string to dump.

-innerLen int

    dump line inner area byte count. (default 4)

-lineLen int

    dump line total byte count. (default 16)

-minOffLen int

	minimum length of the offset field. (default -1)

-offBegin int

    begin dump at file offset.

-offEnd int

    end dump at file offset. (default -1)

-quiet

    if true, suppress header/trailer/informational messages.

-version

	if true, display program version and exit.

## Examples ##

```
#
hd -inFile /etc/hosts
#
hd /etc/hosts
#
cat /etc/hosts | hd
#
hd -lineLen 8 -innerLen 2 /etc/hosts
#
cat /etc/hosts | hd -lineLen 6 -innerLen 3
#
hd -inString "αβγδεζηθικλμνξοπρστυφχψω"
```
