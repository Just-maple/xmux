# Contributing to xmux

Thank you for your interest in contributing to xmux! This document provides guidelines for contributing to the project.

## Development Philosophy

xmux follows these core principles:

1. **Framework Agnosticism**: xmux should remain independent of any specific web framework
2. **Type Safety**: Leverage Go's type system for compile-time validation
3. **Minimal Dependencies**: Keep the core library lightweight and focused
4. **Clean Separation**: Maintain clear separation between business logic and framework concerns

## Getting Started

### Prerequisites
- Go 1.18 or higher
- Basic understanding of Go generics and interfaces

### Setup
1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/xmux.git`
3. Navigate to the project: `cd xmux`
4. Build the project: `go build ./...`

## Contribution Areas

### 1. Core Library Improvements
- Enhance the `Router` interface
- Add new utility functions for route registration
- Improve dependency injection patterns
- Add type-safe middleware support

### 2. Adapter Development
- Create adapters for additional web frameworks
- Improve existing adapter implementations
- Add framework-specific optimizations
- Create example implementations for different use cases

### 3. Documentation
- Improve README and API documentation
- Add usage examples
- Create tutorials for common patterns
- Translate documentation to other languages

### 4. Testing
- Add unit tests for core functionality
- Create integration tests for adapters
- Add benchmark tests for performance-critical code
- Improve test coverage

## Development Guidelines

### Code Style
- Follow standard Go conventions
- Use meaningful variable and function names
- Keep functions focused and concise
- Add comments for public APIs and complex logic

### Testing
- Write tests for new functionality
- Ensure existing tests continue to pass
- Include edge cases and error scenarios
- Use table-driven tests where appropriate

### Documentation
- Update documentation for new features
- Include code examples in documentation
- Document breaking changes clearly
- Keep the README up to date

## Pull Request Process

1. **Create a Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**
   - Write clean, well-documented code
   - Add tests for new functionality
   - Update documentation as needed

3. **Run Tests**
   ```bash
   go test ./...
   go vet ./...
   ```

4. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "Add: brief description of changes"
   ```

5. **Push to Your Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request**
   - Provide a clear description of the changes
   - Reference any related issues
   - Include test results if applicable

## Issue Reporting

When reporting issues, please include:

1. **Description**: Clear description of the issue
2. **Reproduction Steps**: Steps to reproduce the issue
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**: Go version, OS, browser, etc.
6. **Code Example**: Minimal code to reproduce the issue

## Code of Conduct

### Our Pledge
We are committed to making participation in this project a harassment-free experience for everyone.

### Our Standards
- Use welcoming and inclusive language
- Be respectful of differing viewpoints
- Accept constructive criticism gracefully
- Focus on what is best for the community

### Unacceptable Behavior
- Harassment, discrimination, or personal attacks
- Inappropriate language or imagery
- Trolling, insulting comments, or personal/political attacks
- Publishing others' private information without permission

## Questions?

If you have questions about contributing, please:
1. Check the existing documentation
2. Search existing issues and discussions
3. Create a new issue for clarification

Thank you for contributing to xmux!