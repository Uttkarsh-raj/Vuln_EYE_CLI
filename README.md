# Vuln-EYE CLI Tool

The **Vuln-EYE CLI Tool** is a command-line utility designed to **detect vulnerabilities** in Android applications. It analyzes various files to identify security issues, specifically focusing on the **Android manifest** file, **Gradle files** in **Java/Kotlin** projects, and `pubspec.yaml` files in **Flutter** projects and checks for the presence for vulnerable third-party dependencies. The tool is also designed for easy CI integration, allowing it to be seamlessly incorporated into your continuous integration pipeline to automatically scan and assess vulnerabilities in your codebase as part of your build process.

## Features

- **Manifest Analysis**: Scans the Android manifest file for security vulnerabilities such as exposed activities, debuggable flags, and backup settings.
- **Gradle File Analysis**: Checks Gradle files for potential security issues and configurations.
- **Flutter Dependency Checking**: Analyzes the `pubspec.yaml` file for vulnerable dependencies.
- **Dependency Checking**: Identifies known vulnerable third-party dependencies in your project.
- **Fix Suggestions**: Provides the fixed version if available to address detected vulnerabilities.
- **Verbose Output**: Optionally provides detailed output for debugging and understanding the analysis.

## Why Choose Us?

The Vuln-EYE CLI Tool offers a streamlined, command-line-based approach to vulnerability detection for Android applications. Unlike traditional tools, it provides a focused analysis of key files and dependencies and easily integrates with your CI pipeline to analyze your project, making it a seamless addition to your development workflow.

## Setup

### Locally

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/Uttkarsh-raj/Vuln_EYE_CLI.git
   ```

2. **Install Dependencies:**

   ```bash
   cd veye-cli
   go mod tidy
   ```

3. **Build the Tool:**

   ```bash
   go build
   ```

4. **Install the Tool:**

   ```bash
   go install
   ```

5. **Run the Command:**

   Use the `veye` command to scan your project:

   ```bash
   veye scan [flags]
   ```

   **Flags:**

   - `--fix`       Provides the fixed version if available.
   - `--flutter`   Scans for Flutter dependencies.
   - `--verbose`   Provides detailed output.
   - `-h, --help`  Displays help information.

### GitHub Actions Setup

You can integrate the veye CLI Tool into your GitHub Actions workflow. Here’s a sample workflow file:

```yaml
name: Run CLI Tool

on:
  push:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v1

    - name: Pull Docker image
      run: |
        docker pull uttkarshraj/vuln-eye-cli:latest

    # List files in the mounted directory (for debugging)
    - name: List files in mounted directory
      run: |
        docker run -v ${{ github.workspace }}:/repo -w /repo uttkarshraj/vuln-eye-cli:latest ls -alh /repo

    # Run CLI tool with `veye scan` command in the root directory - For Java/Kotlin projects
    - name: Run `veye scan` in the root directory
      run: |
        docker run -v ${{ github.workspace }}:/repo -w /repo uttkarshraj/vuln-eye-cli:latest veye scan

    # Run CLI tool with flags like `veye scan --flutter` in the root directory - For Flutter projects
    - name: Run `veye scan --flutter` in the root directory
      run: |
        docker run -v ${{ github.workspace }}:/repo -w /repo uttkarshraj/vuln-eye-cli:latest veye scan --flutter
```
You can customize the workflow by adding any necessary flags to the veye scan command based on your requirements.<br>
**Flags:**

   - `--fix`       Provides the fixed version if available.
   - `--flutter`   Scans for Flutter dependencies.
   - `--verbose`   Provides detailed output.
   - `-h, --help`  Displays help information.

## Known Limitations

- Currently optimized for Android and Flutter projects.
- The CLI tool focuses on detecting known vulnerable dependencies. The dependency data is sourced from [OSV](https://osv.dev/).
- The accuracy of the vulnerability detection depends on the completeness and currency of the OSV database. Vulnerabilities may not be detected if they are not listed in OSV or if the database is not updated.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
