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
		return fmt.Errorf("component '%s' not found. Available components: button, card, input, checkbox, table, dialog", componentName)
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
};
`,
	"input": `import React from 'react';

export const Input = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(
  ({ className, type, ...props }, ref) => {
    const style: React.CSSProperties = {
      display: 'flex',
      height: '2.25rem',
      width: '100%',
      borderRadius: '0.375rem',
      border: '1px solid #e5e7eb',
      backgroundColor: 'white',
      padding: '0.5rem 0.75rem',
      fontSize: '0.875rem',
      outline: 'none',
      transition: 'border-color 0.2s',
    };

    return <input type={type} style={style} className={className} ref={ref} {...props} />;
  }
);
Input.displayName = "Input";
`,
	"checkbox": `import React from 'react';

export const Checkbox = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(
  ({ className, ...props }, ref) => {
    const style: React.CSSProperties = {
      height: '1rem',
      width: '1rem',
      borderRadius: '0.25rem',
      border: '1px solid #d1d5db',
      cursor: 'pointer',
    };

    return <input type="checkbox" style={style} className={className} ref={ref} {...props} />;
  }
);
Checkbox.displayName = "Checkbox";
`,
	"table": `import React from 'react';

export const Table = ({ className, ...props }: React.HTMLAttributes<HTMLTableElement>) => (
  <div style={{ width: '100%', overflowX: 'auto' }}>
    <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }} className={className} {...props} />
  </div>
);

export const TableHeader = ({ className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) => (
  <thead style={{ borderBottom: '1px solid #e5e7eb' }} className={className} {...props} />
);

export const TableBody = ({ className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) => (
  <tbody className={className} {...props} />
);

export const TableRow = ({ className, ...props }: React.HTMLAttributes<HTMLTableRowElement>) => (
  <tr style={{ borderBottom: '1px solid #f3f4f6', transition: 'background-color 0.2s' }} className={className} {...props} />
);

export const TableHead = ({ className, ...props }: React.ThHTMLAttributes<HTMLTableCellElement>) => (
  <th style={{ height: '3rem', padding: '0 1rem', textAlign: 'left', fontWeight: 500, color: '#6b7280' }} className={className} {...props} />
);

export const TableCell = ({ className, ...props }: React.TdHTMLAttributes<HTMLTableCellElement>) => (
  <td style={{ padding: '1rem' }} className={className} {...props} />
);
`,
	"dialog": `import React from 'react';

export const Dialog = ({ open, onOpenChange, children }: { open: boolean; onOpenChange: (open: boolean) => void; children: React.ReactNode }) => {
  if (!open) return null;

  const overlayStyle: React.CSSProperties = {
    position: 'fixed',
    top: 0, left: 0, right: 0, bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 50,
  };

  const contentStyle: React.CSSProperties = {
    backgroundColor: 'white',
    borderRadius: '0.5rem',
    padding: '1.5rem',
    width: '100%',
    maxWidth: '32rem',
    boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1)',
  };

  return (
    <div style={overlayStyle} onClick={() => onOpenChange(false)}>
      <div style={contentStyle} onClick={(e) => e.stopPropagation()}>
        {children}
      </div>
    </div>
  );
};

export const DialogHeader = ({ children }: { children: React.ReactNode }) => (
  <div style={{ display: 'flex', flexDirection: 'column', gap: '0.375rem', marginBottom: '1rem' }}>
    {children}
  </div>
);

export const DialogTitle = ({ children }: { children: React.ReactNode }) => (
  <h2 style={{ fontSize: '1.125rem', fontWeight: 600, margin: 0 }}>{children}</h2>
);
`,
}
