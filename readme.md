# JTML
JTML is a project I created because I was sick of building HTML pages. While Go templates can make it a little easier, the syntax is still repeated ad nauseum.

## Template

A template can either be a full template or a partial. For now, the partial template filenames begin with underscore (_).

Example

```
#jtml
 #head
 #body
  #container
   #h1 Hello World!
```

It uses an indented syntax.

## Partial
A partial will be named something like `_jtml.txt` and can be included into a template using `#jtml` as shown above. This particular partial is the document root, and can look like this:

```
@open
<!DOCTYPE html>
<html lang="en">

@close
</html>
```

The `@open` and `@close` strings are `directives` and right now only open and close are supported. Due to the indented nature of the templates, these directives tell the template processor that when this partial processes, it will output everything included inside of it wrapped in the `@open` and `@close` content.

## Parameters
Some includes can have parameters. Currently, only one parameter is supported, and only on a single line. I can see the use to have more than one, or named parameters. This can be added in time, and also support for when a parameter isn't included or a case based on what the value of the parameter is.

Example:

_h1.txt
```
<h1>$1</h1>
```

```
#h1 Hello, World!
```

Will output `<h1>Hello, World!</h1>` as expected


# Future Enhancements

- Ability to organize templates into folders and include them like `#admin/login`

- Named parameters
- Conditionals

## Note

I'm sure something like this exists in the wild but I write software and this was a fun little project. Check out the stack that I needed within the `parser.go` and my quick implementation of [stack](https://github.com/jasontconnell/collections).