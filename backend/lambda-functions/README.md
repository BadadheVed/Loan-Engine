# ClickPe Lambda Functions - AWS SAM Deployment

Simplified Lambda deployment with shared dependencies in one package.

## Project Structure

```
backend/
└── lambda-functions/              # Single deployable package
    ├── go.mod                     # Shared dependencies for all functions
    ├── template.yaml              # SAM template
    ├── shared/                    # Shared code
    │   ├── models.go              # User & BatchResult models
    │   └── database.go            # Database connection
    ├── health/                    # Health check function
    │   └── main.go
    └── uploadcsv/                 # CSV upload function
        └── main.go
```

## Shared Dependencies

All Lambda functions share the same `go.mod` with:

- `github.com/aws/aws-lambda-go` - Lambda runtime
- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - PostgreSQL driver
- `github.com/google/uuid` - UUID support

## Single DATABASE_URL

Both functions use the **same `DATABASE_URL`** environment variable configured in `template.yaml`:

```yaml
Globals:
  Function:
    Environment:
      Variables:
        DATABASE_URL: !Ref DatabaseUrl
```

## Quick Start

### 1. Build Functions

```bash
cd backend
make build
```

### 2. Package Everything

```bash
make package
```

Creates `lambda-functions.zip` containing the entire `lambda-functions/` directory.

### 3. Deploy

```bash
# First time (guided)
make deploy-guided

# With DATABASE_URL
make deploy-with-db DATABASE_URL='postgres://user:pass@host:5432/clickpe'

# Subsequent deploys
make deploy
```

## Deployment Options

### Option 1: SAM Deploy (Recommended)

```bash
cd lambda-functions
sam build
sam deploy --guided
```

### Option 2: Upload Zip and Deploy Manually

```bash
# Create the package
make package

# Upload lambda-functions.zip to S3
aws s3 cp lambda-functions.zip s3://my-bucket/

# Update Lambda functions via AWS Console or CLI
```

### Option 3: Make Commands

```bash
make all                    # Clean, build, package
make deploy-guided          # Deploy with prompts
make deploy                 # Deploy with saved config
```

## Local Testing

```bash
# Start local API
make local

# Test health endpoint
make test-health

# Test CSV upload
make test-upload
```

## How It Works

### Shared Package Structure

- `lambda-functions/go.mod` contains ALL dependencies
- `shared/` package exports common code:
  - `shared.User` - User model
  - `shared.BatchResult` - Batch result type
  - `shared.InitDB()` - Database initialization
  - `shared.DB` - Shared database connection

### Function Imports

```go
import (
    "github.com/BadadheVed/clickpe/lambda-functions/shared"
)

func init() {
    shared.InitDB()  // Uses DATABASE_URL env var
}

func saveUsers(users []shared.User) {
    shared.DB.Create(&users)
}
```

## Makefile Commands

```bash
make build          # Build both functions
make package        # Create lambda-functions.zip
make deploy-guided  # First time deploy
make deploy         # Deploy with saved config
make local          # Start local API
make clean          # Remove build artifacts
make all            # Complete workflow
```

## Environment Variables

Set in `template.yaml` parameters:

- `DATABASE_URL` - PostgreSQL connection string (required)
- `Environment` - dev/staging/production
- `WORKER_COUNT` - Number of CSV workers (uploadcsv only)
- `BATCH_SIZE` - Batch size for inserts (uploadcsv only)

## Package Contents

`lambda-functions.zip` includes:

```
lambda-functions/
├── go.mod              # Shared dependencies
├── shared/             # Common code
├── health/
│   ├── main.go
│   └── bootstrap       # Compiled binary
├── uploadcsv/
│   ├── main.go
│   └── bootstrap       # Compiled binary
└── template.yaml       # SAM template
```

## Benefits

✅ **Single go.mod** - All dependencies in one place  
✅ **Shared code** - Models and DB logic reused  
✅ **Single DATABASE_URL** - No duplicate configuration  
✅ **Simple deployment** - Just zip and deploy  
✅ **Easy to maintain** - Everything in `lambda-functions/`

## Cost Optimization

- Both functions share code, reducing total package size
- CloudWatch Logs retention: 7 days
- Lambda timeout: 5 minutes (health), 15 minutes (uploadcsv)
- Memory: 512 MB (health), 1024 MB (uploadcsv)
