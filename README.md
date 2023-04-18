# errors

Golang errors library.

## About The Project

Library defines common mechanism for business error instantiation and handling.

Features:
- Stack tracing
- JSON/XML serialization/deserialization
- Custom fields

Also, package defines most popular error templates:
- ValidationError
- BlockingLink
- ChecksumError
- Unauthorized
- Forbidden
- NotFound
- Timeout
- AlreadyExists
- AlreadyInProgress
- IllegalState
- PreconditionFailed
- PreconditionRequired
- ToManyRequests
- InternalError
- NotImplemented

## Usage

Define error template:

```
var CustomError = errors.Template{
	Code:    "MyCustomError",
	Message: Message("Some user-friendly description"),
	Params: errors.Params{
		errors.Param{Name: "Param1": Value: "param1"},
	},
}
```

Instantiate error:

```
...
if condition {
	return errors.New(CustomError, errors.Param{Name: "Param2": Value: "param2"})
}
...
```

## Similar projects

- [pkg/errors](https://github.com/pkg/errors)
- [cockroachdb/errors](https://github.com/cockroachdb/errors)

## License

Errors is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENCE.md)
for the full license text.

## Contact

- Email: `cherkashin.evgeny.viktorovich@gmail.com`
- Telegram: `@evgeny_cherkashin`