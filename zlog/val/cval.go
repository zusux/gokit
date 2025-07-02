package val

import "context"

// Valuer is returns a log value.
type Valuer func(ctx context.Context) any

// Value return the function value.
func Value(ctx context.Context, v any) any {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}
