# Lox Interpreter in Go

A complete implementation of the Lox programming language in Go, based on the language design from Robert Nystrom's "Crafting Interpreters" book. This project is a feature-complete port of the original Java implementation (jlox) to Go.

## Project Motivation
This project served multiple learning objectives for me:

- As my first substantial project in Go, it provided an excellent opportunity to deeply learn the language through a complex, real-world implementation
- I've always been fascinated by the inner workings of programming languages, and this project allowed me to gain a much greater understanding (and appreciation!) for what is going on under the hood
- I used this project as an opportunity to try Vim (for better or for worse...)

## Features

The interpreter supports all core features of the Lox language, including:

- Complete lexical analysis and tokenization
- Full expression parsing (prefix and infix)
- Rich control flow statements
- First-class functions with closures
- Object-oriented programming with classes
- Single inheritance
- Static variable resolution
- Robust error handling and reporting

## Implementation Details

This implementation is built from the ground up using only Go's standard library. Key components include:

- **Scanner**: Converts source code into tokens
- **Parser**: Builds an Abstract Syntax Tree (AST) using recursive descent parsing
- **Resolver**: Performs static analysis and variable resolution
- **Interpreter**: Executes the parsed code using the Visitor pattern

## Getting Started

### Prerequisites

- Go 1.23 or higher (older versions will likely work but are untested)

### Installation

```bash
git clone https://github.com/AlexKyriacou/go-lox-interpreter.git
cd lox
go build
```

### Usage

To begin, run the program with the `-h` flag to see available options:

```bash
./lox -h
```

## Language Examples

Basic variable declaration and arithmetic:
```lox
var a = 1;
var b = 2;
print a + b;
```

Functions and closures:
```lox
fun makeCounter() {
    var i = 0;
    fun count() {
        i = i + 1;
        return i;
    }
    return count;
}
```

Classes and inheritance:
```lox
class Animal {
    init(name) {
        this.name = name;
    }
}

class Dog < Animal {
    bark() {
        print "Woof!";
    }
}
```

## Acknowledgments

- Robert Nystrom for the original Lox language design and "Crafting Interpreters" book
- The excellent github.com/chidiwilliams/glox – A similar Go port of Lox that I referenced a few times when working through Go-specific implementation details.