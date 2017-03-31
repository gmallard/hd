# A hex dump utility #

Similar to the Unix 'od' utility, but focused on hexadecimal
output.

## Parameters ##

-edgeMark string

    single character at edges if the right hand side. (default "|")

-goDump

    if true, use standard go encoding/hex/Dump.

-h	print usage message.

-hexUpper

    if true, print upper case hex.

-inFile string

    input file name.  Argument 0 may also be used.

-innerLen int

    dump line inner area byte count. (default 4)

-lineLen int

    dump line byte count. (default 16)

-minOffLen int

	minimum lenght of the offset field. (default -1)

-offBegin int

    begin dump at file offset.

-offEnd int

    end dump at file offset. (default -1)

-quiet

    if true, suppress header/trailer/informational messages.
    

## Examples ##

hd -inFile /etc/hosts

hd /etc/hosts

cat /etc/hosts | hd

hd -lineLen 8 -innerLen 2 /etc/hosts
