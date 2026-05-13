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
   #h1 [Hello World!]
```

It uses an indented syntax.

## Partial
A partial will be named something like `_jtml.jtml` and can be included into a template using `#jtml` as shown above. This particular partial is the document root, and can look like this:

```
@open
<!DOCTYPE html>
<html lang="en">

@close
</html>
```

The `@open` and `@close` strings are `directives` and right now only open and close are supported. Due to the indented nature of the templates, these directives tell the template processor that when this partial processes, it will output everything included inside of it wrapped in the `@open` and `@close` content.

## Parameters
Some includes can have parameters. Parameters are passed in the same line, each parameter wrapped with []. They are 1-based when referring to them in a template ($1, $2, ... $N). Currently, only index based parameters are supported. I definitely see a use case for named parameters, and being able to check for parameters.

Example:

_h1.jtml
```
<h1>$1</h1>
```

```
#h1 [Hello, World!]
```

Will output `<h1>Hello, World!</h1>` as expected

Another real world example, if you wanted to have a `<body>` tag be able to take a list of css classes, this is how you would do that.

_body.jtml

```
@open
<body class="$1">

@close
</body>
```

Usage:

```
#body [class1 class2 class3]
```

(This used to be 3 separate parameters before this update on 2026/05/13)

## Template Organization

(Added 2026/05/08) Templates can now be organized into folders. This is only done for partials, so now you can group them into folders and reference them by the path

```
#jtml
 #layout/head
```

Will reference a template in the `layout` folder named `_head.jtml`.

There is no limit on nesting. As I'm starting to use this in real world situations, it became obvious that this was necessary. More updates will happen with the same preface in the future :)

# Future Enhancements

- Named parameters
- Conditionals

## Note

I'm sure something like this exists in the wild but I write software and this was a fun little project. Check out the stack that I needed within the `parser.go` and my quick implementation of [stack](https://github.com/jasontconnell/collections).