# ClickPe Backend - AWS SAM Deployment

This project contains AWS Lambda functions for the ClickPe backend API, deployable via AWS SAM (Serverless Application Model).

## Project Structure

```
backend/
├── lambda-functions/
│   ├── health/              # Health check endpoint
│   │   ├── main.go
│   │   └── go.mod
│   └── uploadcsv/           # CSV upload endpoint
│       ├── main.go
│       └── go.mod
├── template.yaml            # SAM template
├── samconfig.toml           # SAM deployment config (generated)
└── README-SAM.md            # This file
```

## Prerequisites

1. **AWS CLI** - [Install AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
2. **AWS SAM CLI** - [Install SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
3. **Go 1.24+** - [Install Go](https://golang.org/doc/install)
4. **Docker** - Required for SAM build

## Environment Variables

The following environment variables need to be configured:

- `DATABASE_URL` - PostgreSQL connection string (e.g., `postgres://user:password@host:5432/dbname`)
- `Environment` - Deployment environment (dev, staging, production)
- `CorsOrigin` - CORS allowed origin (default: `*`)

## Deployment Steps

### 1. Configure AWS Credentials

```bash
aws configure
```

### 2. Build the Application

```bash
sam build
```

This compiles the Go Lambda functions for Linux x86_64.

### 3. Deploy (First Time - Guided)

```bash
sam deploy --guided
```

You'll be prompted for:

- **Stack Name**: e.g., `clickpe-backend-dev`
- **AWS Region**: e.g., `us-east-1`
- **DatabaseUrl**: Your PostgreSQL connection string
- **Environment**: `dev`, `staging`, or `production`
- **CorsOrigin**: `*` or your specific domain
- **Confirm changes**: Y
- **Allow SAM CLI IAM role creation**: Y
- **Save arguments to configuration file**: Y

### 4. Subsequent Deployments

```bash
sam deploy
```

Uses saved configuration from `samconfig.toml`.

### 5. Deploy Specific Environment

```bash
# Deploy to dev
sam deploy --parameter-overrides Environment=dev

# Deploy to production with specific DB
sam deploy --parameter-overrides Environment=production DatabaseUrl="postgres://user:pass@prod-db:5432/clickpe"
```

## Testing Locally

### Start Local API

```bash
sam local start-api
```

API will be available at `http://127.0.0.1:3000`.

### Test Health Endpoint

```bash
curl http://127.0.0.1:3000/api/health
```

### Test Upload CSV (Local)

```bash
curl -X POST http://127.0.0.1:3000/api/uploadcsv \
  -F "file=@users.csv"
```

## Invoke Functions Directly

### Invoke Health Function

```bash
sam local invoke HealthFunction
```

### Invoke Upload CSV Function

```bash
sam local invoke UploadCSVFunction -e events/uploadcsv-event.json
```

## View Logs

### Tail logs for a function

```bash
# Health function
sam logs -n HealthFunction --stack-name clickpe-backend-dev --tail

# Upload CSV function
sam logs -n UploadCSVFunction --stack-name clickpe-backend-dev --tail
```

### CloudWatch Logs

```bash
# Health function
aws logs tail /aws/lambda/clickpe-health-dev --follow

# Upload CSV function
aws logs tail /aws/lambda/clickpe-uploadcsv-dev --follow
```

## API Endpoints (After Deployment)

After deployment, SAM outputs the API endpoints:

```
https://{api-id}.execute-api.{region}.amazonaws.com/{environment}/api/health
https://{api-id}.execute-api.{region}.amazonaws.com/{environment}/api/uploadcsv
```

Get endpoints:

```bash
aws cloudformation describe-stacks \
  --stack-name clickpe-backend-dev \
  --query 'Stacks[0].Outputs'
```

## Clean Up / Delete Stack

```bash
sam delete --stack-name clickpe-backend-dev
```

## Configuration Files

### samconfig.toml (Auto-generated)

Created after first `sam deploy --guided`. Contains saved deployment parameters.

### template.yaml

SAM template defining:

- Lambda functions
- API Gateway routes
- IAM roles
- CloudWatch log groups
- Environment variables

## Database Connection

Lambda functions connect to PostgreSQL using the `DATABASE_URL` parameter. Ensure:

1. Database is accessible from Lambda (use VPC if needed)
2. Security groups allow Lambda → Database connection
3. Connection string format: `postgres://username:password@hostname:5432/database`

### VPC Configuration (If Needed)

If your database is in a VPC, add to `template.yaml`:

```yaml
Globals:
  Function:
    VpcConfig:
      SecurityGroupIds:
        - sg-xxxxxxxxx
      SubnetIds:
        - subnet-xxxxxxxxx
        - subnet-yyyyyyyyy
```

## Performance Tuning

### Upload CSV Function

- **Timeout**: 900s (15 minutes) for large files
- **Memory**: 1024 MB
- **Workers**: 5 concurrent (configurable via `WORKER_COUNT`)
- **Batch Size**: 100 records (configurable via `BATCH_SIZE`)

Adjust in `template.yaml` as needed.

## Troubleshooting

### Build Errors

```bash
# Clean and rebuild
rm -rf .aws-sam
sam build
```

### Deployment Errors

```bash
# Validate template
sam validate

# View CloudFormation events
aws cloudformation describe-stack-events --stack-name clickpe-backend-dev
```

### Function Errors

```bash
# Check recent logs
sam logs -n UploadCSVFunction --stack-name clickpe-backend-dev
```

## Cost Optimization

- CloudWatch Logs retention: 7 days (adjustable in template)
- Lambda memory: 512 MB (health), 1024 MB (uploadcsv)
- Use reserved concurrency for production if needed

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: Deploy SAM
  run: |
    sam build
    sam deploy --no-confirm-changeset --no-fail-on-empty-changeset
```

## Security

- Database credentials passed via parameters (encrypted in CloudFormation)
- Use AWS Secrets Manager for production:

```yaml
Environment:
  Variables:
    DATABASE_URL: !Sub "{{resolve:secretsmanager:${DatabaseSecretArn}:SecretString:url}}"
```

## Support

For AWS SAM documentation: https://docs.aws.amazon.com/serverless-application-model/
