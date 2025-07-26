# GitHub Streak API

A serverless API that generates beautiful SVG images showing your GitHub contribution streak and statistics. Built with Go and deployed on Vercel.

## Features

- 🎯 **Current Streak**: Shows your current contribution streak with smart logic
- 📊 **Longest Streak**: Displays your longest contribution streak ever
- 🌍 **Timezone Support**: Respects user's timezone for accurate streak calculation
- 🌙 **Auto Theme**: Automatically switches between light/dark theme based on time
- ⚡ **Fast**: Built with Go for high performance
- 🚀 **Serverless**: Deployed on Vercel for global availability

## API Usage

### Endpoint

```
GET /api/streaks
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `username` | string | ✅ Yes | - | GitHub username to fetch data for |
| `timezone` | string | ❌ No | `Asia/Ho_Chi_Minh` | Timezone for streak calculation |

### Examples

#### Basic Usage
```
https://github-met.vercel.app/api/streaks?username=octocat
```

#### With Custom Timezone
```
https://github-met.vercel.app/api/streaks?username=octocat&timezone=America/New_York
```

#### Supported Timezones
- `Asia/Ho_Chi_Minh` (Vietnam)
- `America/New_York` (Eastern Time)
- `America/Los_Angeles` (Pacific Time)
- `Europe/London` (GMT)
- `Asia/Tokyo` (Japan)
- And many more...

### Response

The API returns an SVG image with the following information:

- **Current Streak**: Number of consecutive days with contributions
- **Longest Streak**: Your longest streak ever
- **Total Contributions**: Total contributions across all time
- **Account Creation Date**: When your GitHub account was created
- **Visual Theme**: Light theme (6 AM - 6 PM) or Dark theme (6 PM - 6 AM)

## Streak Logic

The API uses intelligent streak calculation:

- ✅ **Current Day**: If you haven't committed today but had a streak yesterday, your streak is preserved
- ❌ **Next Day**: Only resets streak if you miss commits for 2 consecutive days
- 🌍 **Timezone Aware**: Calculations respect your local timezone

## Environment Variables

To run this API locally, you need to set these environment variables:

```bash
GITHUB_TOKEN=your_github_personal_access_token
GITHUB_GRAPHQL_URL=https://api.github.com/graphql
```

### Getting GitHub Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Click "Generate new token"
3. Select scopes: `read:user`, `read:email`
4. Copy the token and set it as `GITHUB_TOKEN`

## Local Development

### Prerequisites

- Go 1.19 or higher
- GitHub Personal Access Token

### Setup

1. Clone the repository:
```bash
git clone https://github.com/your-username/github-met.git
cd github-met
```

2. Set environment variables:
```bash
export GITHUB_TOKEN=your_github_token
export GITHUB_GRAPHQL_URL=https://api.github.com/graphql
```

3. Run the development server:
```bash
go run main.go
```

### Project Structure

```
github-met/
├── api/
│   └── streaks/
│       └── index.go          # API endpoint handler
├── domain/
│   ├── contribution.service.go # GitHub API integration
│   └── streak.service.go     # Streak calculation logic
├── infrastructure/
│   └── utils/
│       ├── chunkWeeks.go     # Week processing utilities
│       ├── formatNumber.go   # Number formatting
│       ├── rangeOfYears.go   # Year range utilities
│       ├── streak.go         # Streak utilities
│       └── svgBuilder.go     # SVG generation
├── types/
│   └── index.go              # Type definitions
├── main.go                   # Entry point
├── go.mod                    # Go modules
├── go.sum                    # Dependencies checksum
└── vercel.json              # Vercel configuration
```

## Deployment

### Vercel

1. Install Vercel CLI:
```bash
npm i -g vercel
```

2. Deploy:
```bash
vercel
```

3. Set environment variables in Vercel dashboard:
   - `GITHUB_TOKEN`
   - `GITHUB_GRAPHQL_URL`

### Other Platforms

The API can be deployed to any platform that supports Go serverless functions:

- **Netlify Functions**
- **AWS Lambda**
- **Google Cloud Functions**
- **Azure Functions**

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-username/github-met/issues) page
2. Create a new issue with detailed information
3. Include your GitHub username and timezone for debugging

---

Made with ❤️ using Go and deployed on Vercel
