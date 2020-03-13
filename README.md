# golint-extra

## Ignore lint

Use `@ignore-golint-extra`

Example:

```go
h := H{
    // @ignore-golint-extra
    a: "foo",
}
```

## Rules

- `zerovalue_struct`: ban incomplete struct initialization 
