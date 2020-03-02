## scripting-go

This package is used to evaluate expressions written in NCScript format. In the current state it is only meant to evaluate simple expressions
and not able to parse full NCScripts since it is not necessary for the current application.

### Usage

#### parsing expressions

You can parse expressions by calling

```
expression,err:=scripts.Parse(<expression as string>)
```

**err** then will contain any parsing errors. If err is anything but *nil* you can assume expression not to contain anything useful.
When **err** is *nil* however you can then evaluate the expression by calling

```
result,err:=expression.Execute(<variables>)
```

Again **err** will contain any evaluation errors. If everything went well with evaluating the expression, **result** contains the evaluation result.
You can specify **<variables>** if you used any variables in the given expression. If you used variables in the expression and don't specify them
here the evaluation automatically results in an error since the variable value can't get resolved.