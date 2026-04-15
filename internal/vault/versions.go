package vault

import (
	"context"
	"fmt"
)

// VersionMeta holds metadata for a single secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListVersions returns metadata for all versions of a KVv2 secret at path.
func (c *Client) ListVersions(ctx context.Context, path string) ([]VersionMeta, error) {
	metaPath := fmt.Sprintf("%s/metadata/%s", c.mount, path)

	secret, err := c.logical.ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("listing versions for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path %q", path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions key in metadata for %q", path)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format for %q", path)
	}

	var metas []VersionMeta
	for _, v := range versionsMap {
		info, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		meta := VersionMeta{}
		if ct, ok := info["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := info["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if d, ok := info["destroyed"].(bool); ok {
			meta.Destroyed = d
		}
		metas = append(metas, meta)
	}
	return metas, nil
}
