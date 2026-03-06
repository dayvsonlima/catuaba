# Models

Models are Go structs that map to database tables using [GORM](https://gorm.io/). They live in `app/models/`.

## Generating a model

```bash
catuaba g model Product name:string price:float description:text active:boolean
```

This creates two files:
- `app/models/product.go` — the Go struct
- `database/migrations/<timestamp>_create_products.up.sql` — SQL to create the table
- `database/migrations/<timestamp>_create_products.down.sql` — SQL to drop the table

## Generated code

```go
package models

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    Name        string  `json:"name"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
    Active      bool    `json:"active"`
}
```

`gorm.Model` embeds `ID`, `CreatedAt`, `UpdatedAt`, and `DeletedAt` (soft delete).

## Field types

| CLI type | Go type | SQL type | Example |
|----------|---------|----------|---------|
| `string` | `string` | `VARCHAR(255)` | `name:string` |
| `text` | `string` | `TEXT` | `bio:text` |
| `integer` | `int` | `INTEGER` | `age:integer` |
| `float` | `float64` | `REAL` | `price:float` |
| `boolean` | `bool` | `BOOLEAN` | `active:boolean` |
| `datetime` | `time.Time` | `TIMESTAMP` | `published_at:datetime` |

## Querying

Models are queried through the `database.DB` global (a `*gorm.DB` instance):

```go
// Find all
var products []models.Product
database.DB.Find(&products)

// Find by ID
var product models.Product
database.DB.First(&product, 1)

// Where clause
database.DB.Where("active = ?", true).Find(&products)

// Pagination
database.DB.Offset(0).Limit(20).Find(&products)

// Count
var count int64
database.DB.Model(&models.Product{}).Count(&count)

// Create
product := models.Product{Name: "Widget", Price: 9.99}
database.DB.Create(&product)

// Update
database.DB.Save(&product)

// Delete (soft delete — sets deleted_at)
database.DB.Delete(&product)
```

## Relationships

GORM supports all standard relationships. Add them manually to your model:

```go
// app/models/post.go
type Post struct {
    gorm.Model
    Title      string    `json:"title"`
    Body       string    `json:"body"`
    CategoryID uint      `json:"category_id"`
    Category   Category  `json:"category" gorm:"foreignKey:CategoryID"`
    Tags       []Tag     `json:"tags" gorm:"many2many:post_tags;"`
}

// app/models/category.go
type Category struct {
    gorm.Model
    Name  string `json:"name"`
    Posts []Post `json:"posts"`
}
```

Eager-loading:

```go
database.DB.Preload("Category").Preload("Tags").Find(&posts)
```

## Migrations

Catuaba uses raw SQL migrations (not AutoMigrate). Each model generator creates paired up/down files:

```sql
-- database/migrations/20240101120000_create_products.up.sql
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL DEFAULT '',
    price REAL NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    active BOOLEAN NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- database/migrations/20240101120000_create_products.down.sql
DROP TABLE IF EXISTS products;
```

Run migrations:

```bash
catuaba db migrate     # Run all pending
catuaba db rollback    # Rollback last
catuaba db status      # Show status
```

Create a manual migration (e.g., adding an index):

```bash
catuaba g migration add_index_to_products_name
```

Then edit the generated SQL files manually.

## Learn more

- [GORM documentation](https://gorm.io/docs/)
- [GORM associations](https://gorm.io/docs/belongs_to.html)
- [golang-migrate](https://github.com/golang-migrate/migrate)
