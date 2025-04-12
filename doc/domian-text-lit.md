Domain Text Literal
=====

The journey of **Domain Text Literals** in Go+ began with a proposal in early 2024 when a community member suggested adding JSX syntax support to Go+:

* https://github.com/goplus/gop/issues/1770

While JSX has gained widespread adoption in frontend development, particularly in React-based applications, the immediate benefits of building JSX syntax directly into Go+ weren't immediately clear, causing the proposal to be temporarily shelved.

The turning point came when Go+ needed to support [TPL (Text Processing Language)](../tpl/README.md) syntax for the [Go+ Mini Spec](spec-mini.md) project. This necessity prompted a reconsideration of how Go+ should handle domain-specific notations more broadly.

A common understanding in programming language design suggests that **Domain-Specific Languages (DSLs)** often struggle to compete with general-purpose languages. However, this perspective overlooks the fact that numerous domain languages exist and thrive in specialized contexts:

* **Interface description**: HTML, JSX
* **Configuration and data representation**: JSON, YAML, CSV
* **Text syntax representation**: EBNF-like grammar (including TPL syntax), regular expressions
* **Document formats**: Markdown, DOCX, HTML

What distinguishes these domain languages is that they aren't Turing-complete. They lack the full capabilities of general-purpose languages, such as I/O operations, function definitions, and comprehensive flow control structures.

Rather than competing with general-purpose languages, these domain languages typically complement them. Most mainstream programming languages either officially support or have community-built libraries to interact with these domain languages.

This complementary relationship led to the term "**Domain Text Literal**" rather than "**Domain-Specific Language**", emphasizing their role as specialized text formats that can be embedded within general-purpose code.

## Go+'s Approach to Domain Text Literals

After considerable deliberation on how Go+ should support domain text literals, inspiration came from Markdown's code block syntax:

<img src=images/dtl/image-1.png width=960>

Initially, there was consideration to make Go+'s domain text syntax identical to Markdown's. However, this would have prevented Go+ code from being embedded as a domain text within Markdown documents, potentially reducing interoperability between Go+ and Markdown. After careful consideration, the current syntax was chosen to ensure optimal compatibility.

## Built-in Domain Text Literals in Go+

Go+ currently supports several domain text literals natively:

* [TPL](../tpl/README.md)
* JSON, XML, CSV
* Regular expressions (regexp, regexposix)
* HTML (requiring import of `"golang.org/x/net/html"`)

Here are some examples demonstrating their usage:

<img src=images/dtl/image-2.png width=960>

## Extensibility and Implementation

One of the powerful aspects of Domain Text Literals in Go+ is their extensibility. Users can add support for new domain text formats. The `domainTag` represents a package that must have a global `func New(string)` function (with any return type). The domain text is essentially just a call to this function, making the underlying mechanism remarkably simple.

## Beyond Syntactic Sugar

Domain Text Literals offer more than just convenient syntax. They enable Go+ tooling to understand the semantics of these embedded texts rather than treating them as ordinary strings. This semantic understanding enables:

* Code formatters like `gop fmt` to format both Go+ code and supported domain texts simultaneously
* IDE plugins to provide syntax highlighting and advanced features for recognized domain texts
