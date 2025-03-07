# Previous: A powerful web codebase.

Previous is a reference codebase that provides simple, and powerful facilities for building websites and web applications. Previous is not a framework, or library. This project is a showcase of what a web application can look like with simple procedural code, while avoiding needlessly overcomplicated abstractions.

At the end of the day, the job of a "web developer" is to dynamically generate HTML, and serve it back to a requester. Previous was built to remove the friction that is common with popular, overengineered solutions, so you can focus on getting the job done.

Previous includes all the tools you need to handle HTTP requests, interact with a database, and serve HTML.

## Why not another solution?

Other tools and frameworks sell the promise of a "good developer experience" without any efforts towards robustness, longevity, or "true" simplicity. Many of these frameworks deprecate features, change APIs regularly, and have many dependencies, forcing you to do lots of housekeeping to keep your website running.

The code that drives your website / application should be specialized to fit the problems you are trying to solve. By subscribing to a framework, you lock yourself into a set of rules that may not align with those problems. In contrast to a framework, a reference codebase provides all the tools you need to get up and running quickly, while proving the means to extend the code to fit your project.

### Primary Benefits

- Compile to a statically-linked executable - Extremely easy to deploy to your platform of choice; docker not required.
- Everything is statically type-checked by the compiler - This includes HTML components, and database queries.
- Procedural code - No esoteric template / configuration syntax required. Pages are written using pure, procedural Go functions.
- Vendored dependencies - All third party code is included in this codebase as source, and is built alongside your user-level code. This allows for robustness and futureproofing, as you are not dependent on library authors, or package manager repository maintainers.
- Simple and extendable - This project serves as a reference for your individual needs. You are expected to extend or reduce this codebase as needed. For example, if you are building a simple personal site, you may not need a database, or authentication system! Just delete those parts!

## Documentation

Documentation can be found [here](https://github.com/maxamundsen/Previous/wiki).