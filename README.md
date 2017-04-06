collopi
=======

A tiny library to allow making http requests a bit easier using
fluent interfaces.

Main usage is to allow creating common http requests using most common
patterns. It is useful together with an IDE that supports auto-completion.

    client := collopi.NewClient("http://example.com")

    status, err := client.GET().Path("foo", "bar"). // GET request to http://example.com/foo/bar
        Accept(404).                                // Additionally also accept a 404 and do not return an error
        Target(&result).                            // Store the body in the given object (currently only json possible)
        Do()                                        // execute the request

Apart from the Body which is made available by the Target(interface{}) method
there are two return parameters. The int status code returned and an optional
error.