# Contributing to highlight

Thank you for considering contributing to the `highlight` project! This document provides guidelines and instructions for contributing to make the process smooth and effective.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Adding New Language Support](#adding-new-language-support)
- [Style Guidelines](#style-guidelines)
- [Testing](#testing)

## Code of Conduct

Please be respectful and considerate of others when contributing to this project. We welcome contributions from everyone, regardless of level of experience, gender, gender identity and expression, sexual orientation, disability, personal appearance, body size, race, ethnicity, age, religion, or nationality.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/hili-cat.git
   cd hili-cat
   ```
3. **Set up the remote**:
   ```bash
   git remote add upstream https://github.com/AmirMahdyJebreily/hili-cat.git
   ```
4. **Create a branch** for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

The `highlight` project follows a standard Go project layout:

- `cmd/highlight`: Main application entry point
- `internal`: Internal packages used by the application
- `pkg`: Reusable packages that could potentially be used by other projects
- `config`: Configuration files
- `examples`: Example files for demonstration purposes

### Building the Project

Use the included build script:

```bash
# For the current platform
./build.sh

# For Linux specifically
./build.sh --linux

# For all platforms
./build.sh --all
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/ansi
```

## Pull Request Process

1. **Update your fork** with the latest upstream changes:
   ```bash
   git fetch upstream
   git merge upstream/main
   ```

2. **Commit your changes** with clear, descriptive commit messages:
   ```bash
   git add .
   git commit -m "Add feature: your feature description"
   ```

3. **Push your changes** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Submit a Pull Request** to the main repository's `main` branch
   - Describe what your changes do and why they should be included
   - Reference any related issues using `#issue_number`

5. **Code Review**: Wait for maintainers to review your PR
   - Be responsive to feedback and make necessary changes
   - Keep your PR updated with the latest main branch

## Adding New Language Support

One of the most valuable contributions is adding support for new programming languages. Here's how to do it:

1. **Update the configuration file** in `config/config.json`
2. **Test your language** with sample files
3. **Create a pull request** with your changes

See the [Language Configuration Guide](CONFIG_GUIDE.md) for detailed instructions.

## Style Guidelines

- Follow standard Go formatting using `gofmt`
- Use meaningful variable and function names
- Write clear comments for non-obvious code
- Keep functions small and focused on a single responsibility
- Minimize dependencies

## Testing

- Write tests for new features and bug fixes
- Ensure existing tests pass before submitting a PR
- Test your changes on multiple platforms if possible

Thank you for contributing to `highlight`! Your time and expertise help make this tool better for everyone.
