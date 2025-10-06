# Contributing to Yappin

Thank you for your interest in contributing to Yappin! We welcome contributions from the community. This document provides guidelines and information to help you get started.

## Hacktoberfest Participation

This project participates in **Hacktoberfest**! Contributions made during October will count towards your Hacktoberfest progress. Please ensure your pull requests follow the guidelines below to be eligible.

**Hacktoberfest Rules:**
- Pull requests must be made to public repositories
- PRs must be merged or approved by maintainers
- PRs labeled as `hacktoberfest-accepted` count automatically
- Spam PRs will be marked as `invalid`

## Ways to Contribute

- **Code Contributions**: Fix bugs, add features, improve performance
- **Documentation**: Improve README, add code comments, create tutorials
- **Testing**: Write unit tests, integration tests, or help with QA
- **Design**: UI/UX improvements, accessibility enhancements
- **Bug Reports**: Report issues with detailed information
- **Feature Requests**: Suggest new features or improvements

## Development Setup

### Prerequisites

- Go 1.24 or later
- Node.js 18 or later
- Docker and Docker Compose
- Git

### Local Development

**Note for Contributors**: Since this is an open-source project, you'll need to set up your own local development environment. The `.env` file in the repository contains sensitive information and is not included. You can create your own local database setup.

1. **Fork and Clone**
   ```bash
   git clone https://github.com/your-username/yappin.git
   cd yappin
   ```

2. **Set up Environment**
   ```bash
   cd server
   # Copy the example environment file and modify it
   cp .env.example .env
   # Edit .env with your own database credentials (don't commit this file)
   ```

3. **Start Database**
   ```bash
   # Use Docker to run your own PostgreSQL instance
   docker run --name yappin-postgres -e POSTGRES_DB=yappin_dev -e POSTGRES_USER=your_db_user -e POSTGRES_PASSWORD=your_db_password -p 5432:5432 -d postgres:16-alpine

   # Or use docker-compose for both DB and Adminer
   docker-compose up -d db adminer
   ```

4. **Run Migrations**
   ```bash
   cd server
   go run db/migrations/migrate.go up
   cd ..
   ```

5. **Backend Development**
   ```bash
   cd server
   go mod tidy
   go run main.go
   ```

6. **Frontend Development**
   ```bash
   cd client
   npm install
   npm run dev
   ```

7. **Access Points**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8081 (or configured port)
   - Adminer: http://localhost:8080

### Testing Changes Without Full Setup

If you want to contribute code changes but don't want to set up the full database:

- **Frontend-only changes**: You can work on the client code and test UI changes
- **Backend logic**: You can modify Go code and ensure it compiles with `go build`
- **API contracts**: Changes to API endpoints can be reviewed in code review
- **Unit tests**: Write and run unit tests for your changes

For full integration testing, you'll need the database setup above.

## Coding Standards

### Go Backend

- Follow standard Go formatting: `gofmt -w .`
- Use `go vet` and `golint` for code quality
- Follow Go naming conventions
- Write comprehensive error handling
- Add comments for exported functions and types
- Use meaningful variable and function names

### Frontend (Svelte/TypeScript)

- Use TypeScript for all new code
- Follow Prettier formatting: `npm run format`
- Lint code: `npm run lint`
- Use TailwindCSS classes consistently
- Follow SvelteKit best practices
- Write accessible HTML with proper ARIA attributes

### General

- Write clear, concise commit messages
- Keep PRs focused on a single feature or fix
- Add tests for new functionality
- Update documentation as needed

## Git Workflow

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-number-description
   ```

2. **Make Changes**
   - Write code following the standards above
   - Test your changes thoroughly
   - Update documentation if needed

3. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   # Use conventional commit format:
   # feat: new feature
   # fix: bug fix
   # docs: documentation
   # style: formatting
   # refactor: code restructuring
   # test: adding tests
   # chore: maintenance
   ```

4. **Push and Create PR**
   ```bash
   git push origin your-branch-name
   ```
   Then create a pull request on GitHub.

## Pull Request Guidelines

- **Title**: Use clear, descriptive titles
- **Description**: Explain what changes were made and why
- **Testing**: Describe how you tested the changes
- **Screenshots**: Include screenshots for UI changes
- **Related Issues**: Reference any related issues with `#issue-number`

### PR Checklist

- [ ] Code follows project standards
- [ ] Tests pass (if applicable)
- [ ] Documentation updated
- [ ] Commit messages are clear
- [ ] PR description is detailed
- [ ] No merge conflicts

## Testing

### Backend Tests
```bash
cd server
go test ./...
```

### Frontend Tests
```bash
cd client
npm run test
```

## Reporting Issues

When reporting bugs, please include:

- **Steps to reproduce**
- **Expected behavior**
- **Actual behavior**
- **Environment** (OS, browser, Go/Node versions)
- **Screenshots** (if applicable)
- **Error messages** or logs

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming environment for all contributors.

## Getting Help

- Check existing issues and documentation first
- Join our community discussions
- Ask questions in GitHub issues

Thank you for contributing to Yappin! 🎉