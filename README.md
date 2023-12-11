# yajson
Yet Another JSON unmarshaller. It is built at jscan, blazing fast JSON iterator, and uses own unsafe inserting into structures. Main difference between
this and other JSON unmarshallers is, that it compiles a JSON parser for each model. So by that, we eliminate type-unsafety of default standard approach
with `json.Unmarshal(any)`, but also gain performance - instead of checking the passed object every time via reflect, we do this just once - when the parser
is being initialized, after that we do fill attributes by their offsets. Moreover, instead of returning a pointer, we return a new instance of a model by
value - this also brings us to the fair zero-allocations!
