Directory to store application source code files.

- Bootstrapper function, called by `goadmin`, has access to application's global configurations
(type `github.com/go-akka/configuration.Config`) and the Echo server instance.
- Bootstrapper function is responsible for registering request routing.
- Bootstrapper function can replace Echo's renderer with its own. However, calling
`EchoRegisterRenderer(namespace string, renderer echo.Renderer)` to register a namespace-scope renderer is recommended.
Then, when invoking `echo.Context.Render(code int, name string, data interface{})`, template name must be prefixed
with `namespace:` so that it routed to the correct renderer.
