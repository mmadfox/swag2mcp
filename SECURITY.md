# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the latest minor version of `swag2mcp`.

| Version | Supported          |
| ------- | ------------------ |
| latest  | ✅                 |
| < latest| ❌                 |

We highly recommend all users to always use the latest version.

## Reporting a Vulnerability

We take the security of `swag2mcp` seriously. If you believe you have found a security vulnerability, please report it to us as soon as possible.

**Please do not report security vulnerabilities through public GitHub issues, discussions, or pull requests.**

Instead, please send a detailed report to our private security email: **sergey.liskonog@gmail.com**.

### What to Include
To help us triage and address the issue effectively, please include as much of the following information as possible:
*   **Type of issue** (e.g., buffer overflow, denial of service, etc.).
*   **Full paths** of source file(s) related to the manifestation of the issue.
*   The location of the affected source code (tag/branch/commit or direct URL).
*   Any special configuration required to reproduce the issue.
*   Step-by-step instructions to reproduce the issue.
*   Proof-of-concept or exploit code (if possible).
*   Impact of the issue, including how an attacker might exploit it.

This information will help us to quickly verify and address the vulnerability.

### Response Process
1.  You will receive an acknowledgment of your report within **48 hours**.
2.  Our team will investigate the issue and provide you with a status update within **5 business days**.
3.  Once the vulnerability is confirmed, we will work on a fix and release a patch for the latest supported version.
4.  We will publicly disclose the vulnerability after the patch has been released, usually with a mention of the reporter (unless you prefer to stay anonymous).

## Security Measures and Best Practices for Users

When using `swag2mcp` in your environment, we recommend you follow these security best practices:

*   **Keep it Updated:** Always use the latest version of `swag2mcp` to benefit from the latest security patches.
*   **Secure your Specifications:** Be mindful of the OpenAPI/Swagger/Postman specifications you load. If they contain sensitive information (e.g., authentication tokens, internal endpoint details), ensure the files are stored securely and access to them is restricted.
*   **Limit Permissions:** When using `swag2mcp` with OIDC or other authentication methods, grant the MCP server the **minimum required permissions** for your LLM agents to operate. Avoid using overly permissive service accounts.
*   **Review MCP Tools:** The project exposes 16 MCP tools. Regularly review which tools are enabled and the actions they can perform to prevent unintended or malicious use.
*   **Network Security:** Ensure the environment where `swag2mcp` runs is secured and that network access to the APIs it manages is properly controlled.
*   **Go Dependency Management:** We use Go Modules. Regularly run `go get -u` and review your `go.mod` and `go.sum` files to include the latest, secure versions of our dependencies.
