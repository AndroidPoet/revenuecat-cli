# Pagination

RevenueCat CLI uses cursor-based pagination for all list commands. This approach is efficient for large datasets and avoids the inconsistencies of offset-based pagination.

## How It Works

Every list command returns a page of results along with a cursor pointing to the next page. You use that cursor with `--starting-after` to fetch the next page.

```
Page 1                    Page 2                    Page 3
[item_1, item_2, ...]  →  [item_21, item_22, ...]  →  [item_41, ...]
      ↑                          ↑
   --starting-after=item_20   --starting-after=item_40
```

---

## Pagination Flags

| Flag | Description | Default |
|---|---|---|
| `--limit` | Number of items per page (1-100) | `20` |
| `--starting-after` | ID of the last item from the previous page | -- |
| `--all` | Automatically fetch all pages and combine results | `false` |

---

## Manual Pagination

Fetch results page by page using `--starting-after`:

```bash
# First page
rc products list --limit 10
```

The response includes a `next_page` field with the cursor value:

```json
{
  "items": [...],
  "next_page": "prod_xxxxx"
}
```

Use that value to fetch the next page:

```bash
# Second page
rc products list --limit 10 --starting-after prod_xxxxx
```

When `next_page` is `null`, you have reached the end of the dataset.

---

## Automatic Pagination

Use `--all` to fetch every page automatically. The CLI handles cursor management internally and returns the combined result set.

```bash
rc products list --all
```

!!! warning
    Use `--all` with caution on large datasets. For projects with thousands of customers, prefer manual pagination or pipe to a file:

    ```bash
    rc customers list --all --output csv > all-customers.csv
    ```

---

## Combining with Output Formats

Pagination works with every output format:

```bash
# First 50 customers as a table
rc customers list --limit 50 --output table

# All products as CSV
rc products list --all --output csv

# Page through audit logs as YAML
rc audit-logs list --limit 25 --output yaml
```

---

## Scripting Example

Loop through all pages in a shell script:

```bash
#!/bin/bash
cursor=""

while true; do
  if [ -z "$cursor" ]; then
    response=$(rc products list --limit 100 --output json)
  else
    response=$(rc products list --limit 100 --starting-after "$cursor" --output json)
  fi

  # Process this page
  echo "$response" | jq '.items[]'

  # Get next cursor
  cursor=$(echo "$response" | jq -r '.next_page // empty')

  # Exit if no more pages
  [ -z "$cursor" ] && break
done
```

!!! tip
    For most use cases, `--all` is simpler than manual pagination. Reserve manual pagination for very large datasets where you want to process results incrementally.
