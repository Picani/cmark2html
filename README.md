cmark2html
==========

A small utility to compile a set of CommonMark/Markdown files into HTML.

Usage
-----

    usage: cmark2html [<flags>] [<infile.md>...]

    Flags:
          --help               Show context-sensitive help (also try --help-long and --help-man).
      -t, --template=TEMPLATE  the template to use
      -l, --list               list available templates and exit
          --version            Show application version.

    Args:
      [<infile.md>]  input CommonMark file(s)


The templates are HTML files with a `{{content}}` mustache tag. This tag
will be expanded to the result of the infiles [CommonMark](http://commonmark.org/)
compilation.

The templates are looked for in the folder `$XDG_DATA_HOME/cmark2html/`.
See the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
for more information.

The default template compiles to a simple HTML file styled using
[WYSIWYG.css](https://jgthms.com/wysiwyg.css/)


Compilation
-----------

This tool is written in [Go](https://golang.org/) and depends on the following libraries:

* [Kingpin](https://github.com/alecthomas/kingpin) for command-line arguments parsing
* [raymond](https://github.com/aymerick/raymond) for handlebars.js/Mustache templating
* [Blackfriday](https://github.com/russross/blackfriday) for Markdown processing

To install the dependencies, run the following commands:

    $ go get gopkg.in/alecthomas/kingpin.v2
    $ go get github.com/aymerick/raymond
    $ go get gopkg.in/russross/blackfriday.v2

Then, compile with the following:

    $ go build cmark2html.go


Licence
-------

Copyright © 2018 Sylvain PULICANI <picani@laposte.net>

This work is free. You can redistribute it and/or modify it under the
terms of the Do What The Fuck You Want To Public License, Version 2,
as published by Sam Hocevar. See the COPYING file for more details.
