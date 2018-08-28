# gogogo

A simple web framework compatible with `net/http`

```
func A(c *Context) error {
	c.JSON([]byte("a"))
	return nil
}

func main() {
    app.Handle("/", &app.Got{A})
    
    l, err := net.Listen("tcp", ":"+port)
    if err != nil {
        logger.Errorf("Listen: %v", err)
        return
    }
    app.Serve(l)
}
```

A example site: [Cyeam](http://cyeam.com).
