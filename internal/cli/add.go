package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// runAdd injects UI components directly into the user's frontend source tree.
func runAdd(args []string) error {
	if len(args) == 0 {
		return errors.New("you must specify a component to add (e.g. `glyra add button`)")
	}

	componentName := strings.ToLower(args[0])

	// Ensure we are inside a Glyra project
	if !fileExists("frontend/package.json") {
		return errors.New("glyra add only supports React/Next.js templates currently. Could not find frontend/package.json")
	}

	componentCode, ok := reactComponents[componentName]
	if !ok {
		return fmt.Errorf("component '%s' not found. Available components: button, card", componentName)
	}

	// Create frontend/src/components/ui/ directory
	uiDir := filepath.Join("frontend", "src", "components", "ui")
	if err := os.MkdirAll(uiDir, 0755); err != nil {
		return fmt.Errorf("failed to create components directory: %w", err)
	}

	// TitleCase the filename
	fileName := strings.ToUpper(componentName[:1]) + componentName[1:] + ".tsx"
	filePath := filepath.Join(uiDir, fileName)

	if fileExists(filePath) {
		fmt.Printf(theme.Red.Render("⚠ Component %s already exists. Skipping.\n"), fileName)
		return nil
	}

	if err := os.WriteFile(filePath, []byte(componentCode), 0644); err != nil {
		return fmt.Errorf("failed to write component file: %w", err)
	}

	fmt.Println(theme.Green.Render(fmt.Sprintf("✓ Successfully added %s to your project!", fileName)))
	fmt.Println(theme.Dim.Render(fmt.Sprintf("  File saved to: %s", filePath)))

	return nil
}

// In a real production app, these would be fetched from a remote registry or embedded via //go:embed
var reactComponents = map[string]string{
	"button": `import React from 'react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger';
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', ...props }, ref) => {
    // Basic inline styles to avoid needing a complex Tailwind setup for the prototype
    const baseStyle: React.CSSProperties = {
      padding: '0.5rem 1rem',
      borderRadius: '0.375rem',
      fontWeight: 500,
      cursor: 'pointer',
      border: 'none',
      transition: 'all 0.2s ease',
    };

    const variants: Record<string, React.CSSProperties> = {
      primary: { backgroundColor: '#3b82f6', color: 'white' },
      secondary: { backgroundColor: '#f3f4f6', color: '#1f2937' },
      danger: { backgroundColor: '#ef4444', color: 'white' },
    };

    return (
      <button
        ref={ref}
        style={{ ...baseStyle, ...variants[variant] }}
        className={className}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";
`,
	"card": `import React from 'react';

export const Card = ({ className, children, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  const style: React.CSSProperties = {
    backgroundColor: 'white',
    borderRadius: '0.5rem',
    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
    border: '1px solid #e5e7eb',
    overflow: 'hidden',
  };

  return (
    <div style={style} className={className} {...props}>
      {children}
    </div>
  );
};

export const CardHeader = ({ className, children, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  const style: React.CSSProperties = {
    padding: '1.5rem',
    borderBottom: '1px solid #e5e7eb',
  };
  return <div style={style} className={className} {...props}>{children}</div>;
};

export const CardContent = ({ className, children, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  const style: React.CSSProperties = {
    padding: '1.5rem',
  };
  return <div style={style} className={className} {...props}>{children}</div>;
};
`,
}
