# Erlang Rebar Config Parser Documentation

This directory contains the complete documentation for the Erlang Rebar Config Parser library, built with VitePress.

## 📚 Documentation Structure

### English Documentation
- **[Getting Started](./guide/getting-started.md)** - Quick start guide
- **[Installation](./guide/installation.md)** - Installation instructions
- **[Basic Usage](./guide/basic-usage.md)** - Common usage patterns
- **[Advanced Usage](./guide/advanced-usage.md)** - Advanced scenarios and best practices
- **[API Reference](./api/)** - Complete API documentation
- **[Examples](./examples/)** - Real-world examples

### Chinese Documentation (中文文档)
- **[快速开始](./zh/guide/getting-started.md)** - 快速入门指南
- **[安装](./zh/guide/installation.md)** - 安装说明
- **[基本用法](./zh/guide/basic-usage.md)** - 常见使用模式
- **[高级用法](./zh/guide/advanced-usage.md)** - 高级场景和最佳实践
- **[API 参考](./zh/api/)** - 完整的 API 文档
- **[示例](./zh/examples/)** - 实际应用示例

## 🚀 Development

### Prerequisites

- Node.js 18 or later
- npm

### Setup

```bash
# Install dependencies
npm install

# Start development server
npm run docs:dev

# Build for production
npm run docs:build

# Preview production build
npm run docs:preview
```

### Available Scripts

- `npm run docs:dev` - Start development server with hot reload
- `npm run docs:build` - Build static site for production
- `npm run docs:preview` - Preview production build locally

## 📖 Writing Documentation

### Adding New Pages

1. Create a new Markdown file in the appropriate directory
2. Add the page to the navigation in `.vitepress/config.js`
3. Follow the existing structure and style

### Markdown Features

The documentation supports:

- Standard Markdown syntax
- Code syntax highlighting
- Custom containers (tips, warnings, etc.)
- Math expressions
- Mermaid diagrams

### Code Examples

Use fenced code blocks with language specification:

````markdown
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```
````

### Custom Containers

```markdown
::: tip
This is a tip
:::

::: warning
This is a warning
:::

::: danger
This is a dangerous warning
:::
```

## 🌐 Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the main branch via GitHub Actions.

### Manual Deployment

```bash
# Build the documentation
npm run docs:build

# The built files will be in .vitepress/dist/
# Deploy these files to your web server
```

### GitHub Pages Setup

The repository includes a GitHub Actions workflow (`.github/workflows/docs.yml`) that:

1. Builds the documentation on every push to main
2. Deploys to GitHub Pages automatically
3. Makes the documentation available at: `https://scagogogo.github.io/erlang-rebar-config-parser/`

## 📝 Contributing

### Adding Examples

When adding new examples:

1. Create the example file in the appropriate language directory
2. Include complete, runnable code
3. Add explanations for each step
4. Update the examples index page

### Translating Content

To add a new language:

1. Create a new directory under `docs/` (e.g., `docs/fr/` for French)
2. Copy the structure from `docs/zh/` as a template
3. Translate all content
4. Update the VitePress config to include the new language

### Style Guidelines

- Use clear, concise language
- Include practical examples
- Provide both simple and complex use cases
- Keep code examples focused and relevant
- Use consistent formatting and structure

## 🔧 Configuration

The documentation is configured in `.vitepress/config.js`. Key settings:

- **Site metadata**: Title, description, base URL
- **Navigation**: Sidebar and navbar structure
- **Theme**: Colors, fonts, and styling
- **Features**: Search, social links, edit links

## 📊 Analytics

The documentation includes:

- Built-in search functionality
- Mobile-responsive design
- Fast loading with static site generation
- SEO optimization

## 🐛 Troubleshooting

### Common Issues

1. **Build fails**: Check Node.js version (requires 18+)
2. **Links broken**: Ensure all internal links use relative paths
3. **Images not loading**: Place images in `public/` directory
4. **Styling issues**: Check VitePress theme configuration

### Getting Help

- Check the [VitePress documentation](https://vitepress.dev/)
- Review existing examples in the repository
- Open an issue on GitHub for specific problems

## 📄 License

This documentation is part of the Erlang Rebar Config Parser project and is licensed under the MIT License.
