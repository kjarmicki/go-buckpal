# Go Buckpal

This is mostly faithful Go port of [the application](https://github.com/thombergs/buckpal) accompanying [Tom Homberg's book Get Your Hands Dirty on Clean Architecture](https://reflectoring.io/book/).
It was done purely for educational purposes and doesn't necessarily reflect how I'd write an actual production system.

## Use cases, ports and adapters in a nutshell

Use cases are interfaces, descriptions of what users can do with the system. They can be implemented as concrete application service classes. 

Ports are interfaces, they're also part of the application layer. They define the contract with the world outside of an application. Ports can point in (like request handling actions) or out (like persistence actions).

Adapters are the outermost layer of the architecture. Input adapters call input ports to handle incoming traffic. Output adapters call output ports to handle outgoing traffic. Some adapters (for example, out/persistence adapters) are implementations of ports and others (for example, in/web) are not, they call port implementations instead.

## Notes on architecture

- In the original application, repositories were more like DAOs, returning ORM-specific JPA objects. I found that problematic because it made creating generic interfaces for these repositories impossible. I diverged from the original design by making repositories return domain objects instead and thus making it possible to have multiple implementations of the repository.

## Notes on technology

- Go with its design decision that interfaces are implemented implicitly makes it harder than Java to understand which adapters are port implementations and which are not.
- I tried using [go-mysql-server](https://github.com/dolthub/go-mysql-server) as an in-memory engine for repository integration tests. It didn't work out - the library had its own quirks to the point that made tests unreliable (most notably with mapping data types), because its behaviour was signigicantly different than the real database. Additionally, MySQL versions tend to differ in capability, so I ended up using Docker SDK with the exact same database version instead. Slower, yes, but way more dependable.
- GORM is like any other ORM I have ever used - shiny and lovely at first glance / happy path case, but an absolute bitch to debug when it misbehaves. `db.ToSQL` tends to return empty strings with no further explanation. I'd rather go with a simple query builder.
