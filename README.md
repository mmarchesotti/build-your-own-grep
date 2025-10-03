# Go Grep: A Custom `grep` Implementation

This project is a from-scratch implementation of the `grep` command-line utility in Go. It is designed to explore the principles of regular expression engines, including parsing, compilation, and execution using a Non-deterministic Finite Automaton (NFA).

## Features

  * **Pattern Matching**: Search for regex patterns in files or standard input.
  * **File & Stdin Support**: Accepts a list of files to search or reads from `stdin` when no files are provided.
  * **Recursive Search**: Use the `-r` flag to recursively search for patterns within a directory.
  * **Compiler-based Engine**: The regex pattern is compiled into an efficient NFA for matching, avoiding the overhead of backtracking for most patterns.

## Supported Regex Syntax

The engine supports a solid subset of common ERE (Extended Regular Expression) features:

| Feature | Syntax | Example | Description |
| :--- | :--- | :--- | :--- |
| Literals | `a`, `b`, `1` | `cat` | Matches the exact character sequence. |
| Character Classes | `\d`, `\w` | `\d{3}` | Matches digits or word characters. |
| Character Sets | `[...]` | `[abc]` | Matches any character in the set. |
| Negated Sets | `[^...]` | `[^0-9]` | Matches any character not in the set. |
| Wildcard | `.` | `a.c` | Matches any character except newline. |
| Quantifiers | `*`, `+`, `?` | `a*`, `b+`, `c?` | Match zero-or-more, one-or-more, or zero-or-one times. |
| Alternation | `|` | `cat\|dog` | Matches either "cat" or "dog". |
| Grouping | `(...)` | `(ab)+` | Groups expressions for quantifiers or alternation. |
| Positional Anchors | `^`, `$` | `^start`, `end$` | Matches the beginning or end of a line. |

## Architecture

This project is built using a multi-stage, compiler-inspired pipeline to process and execute regular expressions. This design is robust, modular, and easy to extend.

The flow is as follows:

1.  **Lexer (`lexer.go`)**: The raw regex string is fed into the lexer, which breaks it down into a flat sequence of tokens (e.g., `LITERAL`, `KLEENE_CLOSURE`, `GROUPING_OPENER`).

2.  **Parser (`parser.go`)**: The stream of tokens is organized into a hierarchical **Abstract Syntax Tree (AST)**. The AST represents the grammatical structure and precedence of the regex operators.

3.  **NFA Compiler (`build_nfa.go`)**: The AST is traversed and compiled into a **Non-deterministic Finite Automaton (NFA)** using Thompson's construction algorithm. Each node of the AST is converted into a corresponding NFA fragment, which are then linked together to form the complete state machine.

4.  **NFA Simulator (`nfa_simulator.go`)**: The final NFA is executed against each line of input text. The simulator steps through the input character by character, keeping track of all possible active states. If an accepting state is reached, the line is considered a match.

This NFA-based approach is highly efficient for most patterns as it avoids the exponential complexity that can arise from backtracking engines.

## Usage

### Building

To build the executable, run the following command from the project's root directory:

```sh
go build -o mygrep ./cmd/mygrep
```

### Examples

**Search for a pattern in files:**

```sh
./mygrep 'pattern' file1.txt file2.txt
```

**Search from standard input (piping):**

```sh
cat data.log | ./mygrep 'ERROR'
```

**Recursive search within a directory:**

```sh
./mygrep -r 'TODO' ./project_directory
```

## Future Work

The current NFA engine is fast and correct for the features it supports. However, it cannot handle advanced features like **backreferences** (`\1`). The next major development goal is to implement an optional, secondary **backtracking engine**. This engine will reuse the existing Lexer and Parser but will walk the AST directly to enable the stateful matching required for backreferences.