`impl` and `goimpl` are a library and a tool to generate method stubs for implementing interfaces in golang (IDE type functionality).

You type

`sort.Interface ml *musicList`

and it transforms it to 

```
func (ml *musicList) Len() int {
	panic("TODO: implement this method")
}

func (ml *musicList) Less(i int, j int) bool {
	panic("TODO: implement this method")
}

func (ml *musicList) Swap(i int, j int) {
	panic("TODO: implement this method")
}
```

It **can do this for any packages** (not only core libraries, but for any code you use/write).

# How To Specify The Interface?
The [tests](https://github.com/ajmesa9891/impl/blob/master/impl/impl_test.go) have many examples. In short, the path is expected to be in the format of  `<package>.<interface>`, where  `<package>` is exactly as you have to specify it in an import statement. Here are a few examples: 

```
io.Reader
sort.Interface
// Assuming you've written a package you import 
// with "impl/impl/test_data/panther"
impl/impl/test_data/panther.Clawable
// Assuming you've imported package "github.com/the/package/path"
// which has interface "Interface"
github.com/the/package/path.Interface
```

Optionally, you could specify only 1 method to scaffold:

```
sort.Interface::Len
```

# How to Setup?

## With Your Favorite Editor
Since `gomipl` is a command line tool that works well with [go generate](https://blog.golang.org/generate), it can easily integrate with any editor or IDE. Follow the steps below:

1. Get the CLI tool `goimpl` by running the following go command in your terminal

   `go get -u github.com/ajmesa9891/impl/goimpl`

2. Have your editor run `go generate` on file save. If you're using Sublime with [GoSublime](https://github.com/DisposaBoy/GoSublime):
   * Open GoSublime user settings (Preferences > Package Settings > GoSublime > User Settings).
   * Modify your `"on_save"` option to something similar to 

      ```
      "on_save": [
            {"cmd": "gs9o_open", "args": {"run": ["go", "generate"], "focus_view": false}}
      ],
      ```
   * Save the user settings.
3. Add `//go:generate goimpl $GOFILE ` as a snippet. If you're using Sublime with [GoSublime](https://github.com/DisposaBoy/GoSublime):
   * Open GoSublime user settings (Preferences > Package Settings > GoSublime > User Settings).
   * Modify your `"snippets"` option to something similar to 

      ```
      "snippets": [
         {
            "match": {"global": true},
            "snippets": [
               {
                  "text": "goimpl",
                  "title": "go generate impl",
                  "value": "//go:generate goimpl \\$GOFILE $0"
               }
            ]
         }
      ],
      ```
4. Run your snippet, add the interface you're trying to implement and the receiver, save your file, and see the interface scaffolding implemented. If you're using GoSublime press `Esc` to exit the gs9o terminal. For example, try 

   `//go:generate goimpl $GOFILE sort.Interface ml *musicList`

   Save and the comment should have transformed into the interface scaffolding.

## With Anything Else
Essentially, `goimpl` takes a file, an interface, and a receiver, and replaces a comment with the implementation of that interface. The [go generate](https://blog.golang.org/generate) tool allows us to easily integrate it into the golang ecosystem. Try using the tool with go generate alone to understand how to integrate it with anything else.

# Why 2? `impl` & `goimpl`?
`impl` is a library to create interface stubs and can only be used programmatically. `goimpl` is a command layered on top that makes it easy to use with `go generate`.