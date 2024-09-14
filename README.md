# Gitup CLI Tool

The **Gitup CLI Tool** is a command-line utility designed to **detect vulnerabilities** in Android applications. It analyzes various files to identify security issues, specifically focusing on the **Android manifest** file, **Gradle files** in **Java/Kotlin** projects, and `pubspec.yaml` files in **Flutter** projects and checks for the presence for vulnerable third-party dependencies.

## Features

- **Manifest Analysis**: Scans the Android manifest file for security vulnerabilities such as exposed activities, debuggable flags, and backup settings.
- **Gradle File Analysis**: Checks Gradle files for potential security issues and configurations.
- **Flutter Dependency Checking**: Analyzes the `pubspec.yaml` file for vulnerable dependencies.
- **Dependency Checking**: Identifies known vulnerable third-party dependencies in your project.
- **Fix Suggestions**: Provides the fixed version if available to address detected vulnerabilities.
- **Verbose Output**: Optionally provides detailed output for debugging and understanding the analysis.

## Why Choose Us?

Unlike traditional static analysis tools, the **Gitup CLI Tool** offers a streamlined, context-specific approach tailored to Android and Flutter projects. With a focus on comprehensive vulnerability detection and actionable fix suggestions, it provides a practical and efficient solution for maintaining application security.

## Setup

### Locally

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/gitup-cli.git
   ```

2. **Install Dependencies:**

   ```bash
   cd gitup-cli
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

   Use the `gitup` command to scan your project:

   ```bash
   gitup scan [flags]
   ```

   **Flags:**

   - `--fix`       Provides the fixed version if available.
   - `--flutter`   Scans for Flutter dependencies.
   - `--verbose`   Provides detailed output.
   - `-h, --help`  Displays help information.

### GitHub Actions Setup

You can integrate the Gitup CLI Tool into your GitHub Actions workflow. Here’s a sample workflow file:

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
        docker pull uttkarshraj/gitup-cli:latest

    # List files in the mounted directory (for debugging)
    - name: List files in mounted directory
      run: |
        docker run -v ${{ github.workspace }}:/repo -w /repo uttkarshraj/gitup-cli:latest ls -alh /repo

    # Run CLI tool with `gitup scan` command in the root directory - For Java/Kotlin projects
    - name: Run `gitup scan` in the root directory
      run: |
        docker run -v ${{ github.workspace }}:/repo -w /repo uttkarshraj/gitup-cli:latest gitup scan

    # Run CLI tool with flags like `gitup scan --flutter` in the root directory - For Flutter projects
    - name: Run `gitup scan --flutter` in the root directory
      run: |
        docker run -v ${{ github.workspace }}:/repo -w /repo uttkarshraj/gitup-cli:latest gitup scan --flutter
```
You can customize the workflow by adding any necessary flags to the gitup scan command based on your requirements.<br>
**Flags:**

   - `--fix`       Provides the fixed version if available.
   - `--flutter`   Scans for Flutter dependencies.
   - `--verbose`   Provides detailed output.
   - `-h, --help`  Displays help information.

## Known Limitations

- Currently optimized for Android and Flutter projects.
- The CLI tool may not cover all types of vulnerabilities and configurations.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
