# Contributing to GoPortScanner

Thank you for your interest in contributing to GoPortScanner! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

- Use the GitHub issue tracker
- Include detailed steps to reproduce the bug
- Provide system information (OS, Go version, etc.)
- Include error messages and logs
- Check if the issue has already been reported

### Suggesting Enhancements

- Use the GitHub issue tracker with the "enhancement" label
- Clearly describe the proposed feature
- Explain why this feature would be useful
- Consider the impact on existing functionality

### Code Contributions

- Fork the repository
- Create a feature branch from `main`
- Make your changes following the coding standards
- Add tests for new functionality
- Update documentation as needed
- Submit a pull request

## Development Setup

### Prerequisites

- Go 1.24 or higher
- Git
- Make (optional, for using Makefile)

### Local Development

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/goportscanner.git
   cd goportscanner
   ```

2. **Add the upstream remote**
   ```bash
   git remote add upstream https://github.com/rancmd/goportscanner.git
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Build the project**
   ```bash
   go build -o goportscanner cmd/goportscanner/main.go
   ```

5. **Run tests**
   ```bash
   go test ./...
   ```

### Development Tools

We recommend installing these tools for development:

```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
go install golang.org/x/tools/cmd/goimports@latest
```

## Coding Standards

### Go Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` or `goimports` to format code
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Code Organization

- Keep functions small and focused
- Use meaningful variable and function names
- Add comments for complex logic
- Follow the existing package structure

### Error Handling

- Always check and handle errors appropriately
- Use meaningful error messages
- Return errors rather than panicking
- Use `fmt.Errorf` with context for errors

### Documentation

- Add comments for exported functions and types
- Update README.md for user-facing changes
- Include examples in documentation
- Keep documentation up to date with code changes

## Testing

### Writing Tests

- Write tests for all new functionality
- Use descriptive test names
- Test both success and error cases
- Use table-driven tests when appropriate
- Mock external dependencies

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestFunctionName ./...
```

### Test Coverage

- Aim for at least 80% test coverage
- Focus on critical paths and edge cases
- Test error conditions and boundary values

## Pull Request Process

### Before Submitting

1. **Ensure tests pass**
   ```bash
   go test ./...
   ```

2. **Run linting**
   ```bash
   golangci-lint run
   ```

3. **Run security checks**
   ```bash
   gosec ./...
   ```

4. **Update documentation** if needed

5. **Test your changes** manually

### Pull Request Guidelines

- Use a descriptive title
- Include a detailed description of changes
- Reference related issues
- Include test results
- Add screenshots for UI changes (if applicable)

### Review Process

- All PRs require at least one review
- Address review comments promptly
- Maintainers may request changes
- PRs are merged after approval

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality in a backward-compatible manner
- **PATCH**: Backward-compatible bug fixes

### Creating a Release

1. **Update version** in relevant files
2. **Update CHANGELOG.md** with changes
3. **Create a release tag**
4. **Build and test** release artifacts
5. **Publish release** on GitHub

### Release Checklist

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated
- [ ] Version is updated
- [ ] Release notes are written
- [ ] Binaries are built and tested

## Security

### Reporting Security Issues

- **Do not** report security issues in public GitHub issues
- Email security issues to the maintainers
- Include detailed information about the vulnerability
- Allow time for assessment and fix

### Security Best Practices

- Follow secure coding practices
- Validate all inputs
- Use secure defaults
- Avoid common vulnerabilities (SQL injection, XSS, etc.)

## Communication

### Getting Help

- Use GitHub Discussions for questions
- Check existing issues and discussions
- Be respectful and patient
- Provide context when asking questions

### Community Guidelines

- Be respectful and inclusive
- Help others learn and grow
- Share knowledge and experiences
- Follow the project's code of conduct

## Recognition

Contributors will be recognized in:

- The project's README.md
- Release notes
- GitHub contributors list
- Project documentation

Thank you for contributing to GoPortScanner! Your contributions help make this tool better for everyone.