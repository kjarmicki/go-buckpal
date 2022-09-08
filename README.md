# Go Buckpal

This is mostly faithful Go port of the application accompanying Tom Homberg's book Get Your Hands Dirty on Clean Architecture.

# Notes on technology

- In the original application, repositories were more like DAOs, returning ORM-specific JPA objects. I found that problematic because it made creating generic interfaces for these repositories impossible. I diverged from the original design by making repositories return domain objects instead and thus making it possible to have multiple implementations of the repository.
- I tried using [go-mysql-server](https://github.com/dolthub/go-mysql-server) as an in-memory engine for repository integration tests. It didn't work out - the library had its own quirks to the point that made tests unreliable (most notably with mapping data types), because its behaviour was signigicantly different than the real database. Additionally, MySQL versions tend to differ in capability, so I ended up using Docker SDK with the exact same database version instead. Slower, yes, but way more dependable.
- GORM is like any other ORM I have ever used - shiny and lovely at first glance, but an absolute bitch to debug. `db.ToSQL` tends to return empty strings with no further explanation. I'd rather go with a simple query builder.

