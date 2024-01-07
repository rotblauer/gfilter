# gfilter

Filter JSON lines with [GJSON](https://github.com/tidwall/gjson) [query syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md).

```shell
$ go install .

$ gfilter --help
Usage of gfilter:
  -match-all string
        match all of these properties (gjson syntax, comma separated queries)
  -match-any string
        match any of these properties (gjson syntax, comma separated queries)
  -match-none string
        match none of these properties (gjson syntax, comma separated queries)

$ cat <<EOF | gfilter --match-all '#(name.first=="Janet")'
{"name":{"first":"Janet","last":"Prichard"},"age":47}
{"name":{"first":"Carol","last":"Smith"},"age":49}
{"name":{"first":"John","last":"Smith"},"age":42}
{"name":{"first":"Lisa","last":"Smith"},"age":49}
EOF

{"name":{"first":"Janet","last":"Prichard"},"age":47}
```
