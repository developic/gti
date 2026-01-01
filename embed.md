# Go Embedding Guide

## What is Embedding?

Embedding in Go allows you to include static files (like text files, images, templates, etc.) directly into your compiled binary. This means your application becomes self-contained and doesn't need external asset files to function.

## How to Use `//go:embed`

### Basic Syntax

```go
import "embed"

//go:embed filename.txt
var embeddedFile embed.FS

//go:embed directory/*
var embeddedDir embed.FS
```

### Key Points

1. **Directive Placement**: The `//go:embed` comment must be placed immediately before a variable declaration
2. **Variable Type**: Use `embed.FS` for the variable type
3. **Path Resolution**: Paths are resolved relative to the Go source file containing the directive
4. **Patterns**: You can use glob patterns like `*` and `**` for matching files

### Reading Embedded Files

```go
// Read a single file
data, err := embeddedFile.ReadFile("filename.txt")
if err != nil {
    // handle error
}

// Read from embedded directory
entries, err := embeddedDir.ReadDir("directory")
if err != nil {
    // handle error
}

// Read a specific file from embedded directory
data, err := embeddedDir.ReadFile("directory/filename.txt")
```

## Real-World Example: GTI Project

In this project, we embedded themes and word lists to make the binary self-contained.

### Embedding Themes

```go
// In src/cmd/theme.go
//go:embed themes/*
var embeddedThemes embed.FS

func loadAvailableThemes() map[string]config.ThemeColorsConfig {
    themes := make(map[string]config.ThemeColorsConfig)

    entries, err := embeddedThemes.ReadDir("themes")
    if err != nil {
        return themes
    }

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        themeName := entry.Name()
        themePath := "themes/" + themeName

        if colors, err := loadThemeFromEmbeddedFile(themePath); err == nil {
            themes[themeName] = colors
        }
    }

    return themes
}
```

### Embedding Words

```go
// In src/internal/words/generator.go
//go:embed words/*
var embeddedWords embed.FS

func loadWords(language string) []string {
    // ... language file resolution logic ...

    filePath := "words/" + fileName

    data, err := embeddedWords.ReadFile(filePath)
    if err != nil {
        return defaultWords
    }

    // Parse the embedded file content
    var words []string
    scanner := bufio.NewScanner(strings.NewReader(string(data)))
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            words = append(words, line)
        }
    }

    return words
}
```

## Benefits

1. **Self-Contained Binaries**: No need to distribute separate asset files
2. **Simplified Deployment**: Single binary deployment
3. **Version Consistency**: Assets are versioned with the code
4. **Security**: Assets are compiled into the binary, harder to tamper with
5. **Performance**: Faster access since files are memory-mapped

## Limitations

1. **Build-Time Only**: Files must exist at compile time
2. **Read-Only**: Embedded files are read-only
3. **No Dynamic Updates**: Can't modify embedded files at runtime
4. **Increased Binary Size**: All embedded files increase the binary size
5. **Path Restrictions**: Can't use `../` in embed paths

## Best Practices

1. **Organize Assets**: Keep embedded assets in dedicated directories
2. **Use Appropriate Patterns**: Use `*` for files in a directory, `**` for recursive
3. **Error Handling**: Always handle errors when reading embedded files
4. **Fallbacks**: Provide sensible defaults when embedded files can't be read
5. **Documentation**: Document what files are embedded and why

## Common Use Cases

- **Web Assets**: HTML, CSS, JS, images for web applications
- **Configuration Templates**: Default config files
- **Data Files**: JSON, CSV, text data files
- **Localization**: Language files, translation data
- **Static Content**: Help text, documentation, licenses

## Troubleshooting

### "no matching files found"
- Check that the files exist at the specified paths
- Verify the path is relative to the Go source file
- Ensure the pattern matches the actual filenames

### "invalid pattern syntax"
- Don't use `../` in embed paths
- Use forward slashes `/` even on Windows
- Glob patterns must be valid

### Files not updating
- Rebuild the binary after changing embedded files
- The embed directive is evaluated at compile time

## Migration from Filesystem Loading

When converting from filesystem-based asset loading to embedding:

1. Add the `//go:embed` directive
2. Change file reading code to use `embed.FS` methods
3. Update error handling for embedded file access
4. Test that all assets are properly embedded
5. Remove filesystem-based asset distribution

This approach ensures your Go applications are truly self-contained and portable.
